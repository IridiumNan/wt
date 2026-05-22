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
	cmd := exec.Command("tailscale", "status")

	var out []byte
	out, err = cmd.Output()
	if err != nil {
		return
	}

	strOut := string(out)
	tailnetName := strings.Fields(strOut)[2]
	link := "https://" + tailnetName

	base, _ := url.Parse(link)

	linkPrefix := base.JoinPath(subPath)
	baseLink = linkPrefix

	return
}

func turnOnFunnel() (err error) {
	linkDir = filepath.Join(config.DataDir, "public")
	linkDir, _ = filepath.Abs(linkDir)

	err = os.MkdirAll(linkDir, 0o755)
	if err != nil {
		return
	}

	cmd := exec.Command("tailscale", "funnel", "--bg", "--set-path", subPath, linkDir)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

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
	if !FunnelOn {
		err = turnOnFunnel()
		if err != nil {
			return
		}
	}
	pkgPath := filepath.Join(config.DataDir, pkgName)
	pkgPath, _ = filepath.Abs(pkgPath)

	softLinkPath := filepath.Join(linkDir, pkgName)

	err = os.Symlink(pkgPath, softLinkPath)
	link = getPackageLink(pkgName)

	if os.IsExist(err) {
		addNewLinkIfNotInPool(link)
		err = nil
		return
	}
	if err != nil {
		return
	}

	addNewLinkIfNotInPool(link)

	return
}

func privateSinglePackage(pkgName string) (err error) {
	softLinkPath := filepath.Join(linkDir, pkgName)
	err = os.Remove(softLinkPath)
	if err != nil && os.IsNotExist(err) {
		err = nil
		return
	}

	if err != nil {
		return err
	}

	pkgLinkPool = slices.DeleteFunc(pkgLinkPool, func(link string) bool {
		return strings.Contains(link, pkgName)
	})

	return
}
