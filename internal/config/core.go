package config

import (
	"errors"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

func GetTokenHeadName(wtMethod model.WTMethod) (headName string) {
	switch wtMethod {
	case model.WTRead:
		headName = "X-Read-Token"
	case model.WTInstall:
		headName = "X-Install-Token"
	case model.WTWrite:
		headName = "X-Write-Token"
	}

	return
}

// GetTimeout : get timeout by wtType and wtMethod
func GetTimeout(wtType model.WTType, wtMethod model.WTMethod) (timeout time.Duration) {
	switch wtMethod {
	case model.WTRead:
		if wtType == model.WTServer {
			timeout = serverConfig.ReadTimeout
		} else {
			timeout = clientConfig.ReadTimeout
		}
	case model.WTInstall:
		if wtType == model.WTServer {
			timeout = serverConfig.InstallTimeout
		} else {
			timeout = clientConfig.InstallTimeout
		}
	case model.WTWrite:
		if wtType == model.WTServer {
			timeout = serverConfig.WriteTimeout
		} else {
			timeout = clientConfig.WriteTimeout
		}
	}
	return
}

// GetTokenList : get TokenList from server config by wtMethod
func GetTokenList(wtMethod model.WTMethod) (tokenList []string) {
	switch wtMethod {
	case model.WTRead:
		tokenList = serverConfig.ReadToken
	case model.WTInstall:
		tokenList = serverConfig.InstallToken
	case model.WTWrite:
		tokenList = serverConfig.WriteToken
	}

	return
}

func GetToken(wtMethod model.WTMethod) (token string) {
	switch wtMethod {
	case model.WTRead:
		token = clientConfig.ReadToken
	case model.WTInstall:
		token = clientConfig.InstallToken
	case model.WTWrite:
		token = clientConfig.WriteToken
	}
	return
}

func GetServerAddr(wtType model.WTType) (addr string) {
	switch wtType {
	case model.WTServer:
		addr = serverConfig.Server
	case model.WTClient:
		addr = clientConfig.Server
	}

	return
}

func GetTagList() (tagList []string) {
	for tag := range serverConfig.TagTokenMap {
		tagList = append(tagList, tag)
	}

	return
}

func GetTagTokenList(tag string, wtMethod model.WTMethod) (tokenList []string, err error) {
	tagAuthTokens, ok := serverConfig.TagTokenMap[tag]

	if !ok {
		return nil, errors.New("tag not found")
	}

	switch wtMethod {
	case model.WTRead:
		tokenList = tagAuthTokens.ReadTokens
	case model.WTInstall:
		tokenList = tagAuthTokens.Installtokens
	case model.WTWrite:
		tokenList = tagAuthTokens.WriteTokens
	}

	return
}

func AddTagTokenList(tag string) {
	if _, exist := serverConfig.TagTokenMap[tag]; exist {
		return
	}

	serverConfig.TagTokenMap[tag] = model.TagAuthTokens{
		ReadTokens:    []string{},
		Installtokens: []string{},
		WriteTokens:   []string{},
	}
}

func DeleteTagTokenList(tag string) {
	if _, exist := serverConfig.TagTokenMap[tag]; !exist {
		return
	}

	delete(serverConfig.TagTokenMap, tag)
}

func GetAllTags() (tagList []string) {
	for tag := range serverConfig.TagTokenMap {
		if tag == DefaultTagStatic || tag == DefaultTagTemp {
			continue
		}
		tagList = append(tagList, tag)
	}

	return
}
