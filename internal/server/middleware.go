package server

import (
	"slices"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
)

// tokenInList : check if token in tokenList
func tokenInList(token string, tokenList []string) (inList bool) {
	inList = slices.Contains(tokenList, token)
	return
}

func IsAccess(wtMethod config.WTMethod, clientToken string) (valid bool) {
	if clientToken == "" {
		return
	}
	switch wtMethod {
	case config.WTRead:
		valid = tokenInList(clientToken, config.GetTokenList(wtMethod))
	case config.WTInstall:
		valid = tokenInList(clientToken, config.GetTokenList(wtMethod))
	case config.WTWrite:
		valid = tokenInList(clientToken, config.GetTokenList(wtMethod))
	}
	return
}
