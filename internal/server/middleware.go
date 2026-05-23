package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/pkg/httphelper"
)

var (
	FunnelOn    = false
	baseLink, _ = getLinkPrefix()
	linkDir, _  = filepath.Abs(filepath.Join(config.DataDir, "public"))
	subPath     = "/wt/public"
	pkgLinkPool []string
)

// tokenInList : check if token in tokenList
func tokenInList(token string, tokenList []string) (inList bool) {
	inList = slices.Contains(tokenList, token)
	return
}

func isAccess(auth model.Auth) (valid bool) {
	clientToken := auth.Token

	authByte, _ := json.Marshal(auth)
	slog.Debug("check auth", "auth", string(authByte))
	if clientToken == "" {
		return
	}
	if tokenInList(clientToken, config.GetTokenList(auth.WtMethod)) {
		return true
	}

	if auth.Tag != "" {
		tagToken, err := config.GetTagTokenList(auth.Tag, auth.WtMethod)
		if err == nil && tokenInList(auth.Token, tagToken) {
			return true
		}
	}

	return
}

func GetAvailableTags(clientToken string, wtMethod model.WTMethod) []string {
	allTags := []string{}
	allTags = append(allTags, config.DefaultTagTemp)
	allTags = append(allTags, config.DefaultTagStatic)
	allTags = append(allTags, config.GetAllTags()...)

	if tokenInList(clientToken, config.GetTokenList(wtMethod)) {
		return allTags
	}

	var accessible []string
	for i := range allTags {
		tagTokens, err := config.GetTagTokenList(allTags[i], wtMethod)

		if err == nil && tokenInList(clientToken, tagTokens) {
			accessible = append(accessible, allTags[i])
		}
	}

	return accessible
}

// TokenOk : check the header token for single tag or global token
func TokenOk(w http.ResponseWriter, r *http.Request, wtMethod model.WTMethod, tag string) (pass bool) {
	pass = true
	tokenHeadName := config.GetTokenHeadName(wtMethod)

	var errMsg string
	switch wtMethod {
	case model.WTRead:
		errMsg = "has no access to read"
	case model.WTInstall:
		errMsg = "has no access to install"
	case model.WTWrite:
		errMsg = "has no access to write"
	}

	auth := model.Auth{
		WtMethod: wtMethod,
		Token:    r.Header.Get(tokenHeadName),
		ErrMsg:   errMsg,
		Tag:      tag,
	}

	if !isAccess(auth) {
		httphelper.SendJSONResponse(
			w,
			http.StatusForbidden,
			model.ForbiddenResponse(errMsg),
		)
		pass = false
	}

	return
}

func getLinkPrefix() (baseLink *url.URL, err error) {
	slog.Debug("Getting Tailscale link prefix", "command", "tailscale status")
	cmd := exec.Command("tailscale", "dns", "status")

	lines, err := cmd.Output()
	if err != nil {
		return
	}

	strLines := strings.Split(string(lines), "\n")

	var tailnetName string
	for i := range strLines {
		if strings.Contains(strLines[i], "Other devices in your tailnet") {
			targetLine := strings.Fields(strLines[i])
			tailnetName = targetLine[10][:len(targetLine[10])-1]
		}
	}

	link := "https://" + tailnetName

	slog.Debug("Tailscale link prefix obtained", "tailnet_name", tailnetName, "link", link)

	base, _ := url.Parse(link)

	linkPrefix := base.JoinPath(subPath)
	baseLink = linkPrefix

	return
}

func turnOnFunnel() (err error) {
	slog.Info("Starting Tailscale Funnel", "data_dir", config.DataDir, "sub_path", subPath)
	linkDir = filepath.Join(config.DataDir, "public")
	linkDir, _ = filepath.Abs(linkDir)

	err = os.MkdirAll(linkDir, 0o755)
	if err != nil {
		slog.Error("Failed to create public directory", "error", err.Error(), "path", linkDir)
		return
	}

	slog.Debug("Executing tailscale funnel command", "command", fmt.Sprintf("tailscale funnel --bg --set-path %s %s", subPath, linkDir))
	cmd := exec.Command("tailscale", "funnel", "--bg", "--set-path", subPath, linkDir)
	_, err = cmd.Output()
	if err != nil {
		slog.Error("Failed to start Tailscale Funnel", "error", err.Error())
		fmt.Println(err)
		return
	}

	slog.Info("Tailscale Funnel started successfully", "link_dir", linkDir)
	FunnelOn = true

	return
}

func getPackageLink(pkgName string) string {
	pkgLink := baseLink.JoinPath(pkgName)

	return pkgLink.String()
}

func addNewLinkIfNotInPool(newLink string) {
	if slices.Contains(pkgLinkPool, newLink) {
		return
	}

	pkgLinkPool = append(pkgLinkPool, newLink)
}

func exposeSinglePackage(pkgName string) (link string, err error) {
	slog.Info("Exposing package via Funnel", "package_name", pkgName, "funnel_active", FunnelOn)

	if !FunnelOn {
		slog.Info("Funnel not active, starting now")
		err = turnOnFunnel()
		if err != nil {
			slog.Error("Failed to start Funnel", "error", err.Error())
			return
		}
	}

	pkgPath := filepath.Join(config.DataDir, pkgName)
	pkgPath, _ = filepath.Abs(pkgPath)

	softLinkPath := filepath.Join(linkDir, pkgName)

	slog.Debug("Creating symlink", "source", pkgPath, "target", softLinkPath)
	err = os.Symlink(pkgPath, softLinkPath)
	link = getPackageLink(pkgName)

	if os.IsExist(err) {
		slog.Info("Symlink already exists", "package_name", pkgName, "link", link)
		addNewLinkIfNotInPool(link)
		err = nil
		return
	}
	if err != nil {
		slog.Error("Failed to create symlink", "error", err.Error(), "package_name", pkgName)
		return
	}

	addNewLinkIfNotInPool(link)
	slog.Info("Package exposed successfully", "package_name", pkgName, "public_link", link)

	return
}

func privateSinglePackage(pkgName string) (err error) {
	slog.Info("Making package private", "package_name", pkgName)

	softLinkPath := filepath.Join(linkDir, pkgName)
	slog.Debug("Removing symlink", "symlink_path", softLinkPath)

	err = os.Remove(softLinkPath)
	if err != nil && os.IsNotExist(err) {
		slog.Info("Symlink does not exist, nothing to remove", "package_name", pkgName)
		err = nil
		return
	}

	if err != nil {
		slog.Error("Failed to remove symlink", "error", err.Error(), "package_name", pkgName)
		return err
	}

	pkgLinkPool = slices.DeleteFunc(pkgLinkPool, func(link string) bool {
		return strings.Contains(link, pkgName)
	})

	slog.Info("Package made private successfully", "package_name", pkgName)
	return
}
