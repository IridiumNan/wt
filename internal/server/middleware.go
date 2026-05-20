package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/pkg/httphelper"
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
