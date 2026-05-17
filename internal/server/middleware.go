package server

import (
	"encoding/json"
	"log/slog"
	"slices"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
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
	switch auth.WtMethod {
	case model.WTRead:
		valid = tokenInList(clientToken, config.GetTokenList(auth.WtMethod))
	case model.WTInstall:
		valid = tokenInList(clientToken, config.GetTokenList(auth.WtMethod))
	case model.WTWrite:
		valid = tokenInList(clientToken, config.GetTokenList(auth.WtMethod))
	}
	return
}
