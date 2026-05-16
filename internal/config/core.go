package config

import (
	"time"
)

type (
	WTMethod uint8
	WTType   uint8
)

const (
	WTRead    WTMethod = iota
	WTInstall WTMethod = iota
	WTWrite   WTMethod = iota
)

const (
	WTServer WTType = iota
	WTClient WTType = iota
)

func GetTokenHeadName(wtMethod WTMethod) (headName string) {
	switch wtMethod {
	case WTRead:
		headName = "X-Read-Token"
	case WTInstall:
		headName = "X-Install-Token"
	case WTWrite:
		headName = "X-Write-Token"
	}

	return
}

// GetTimeout : get timeout by wtType and wtMethod
func GetTimeout(wtType WTType, wtMethod WTMethod) (timeout time.Duration) {
	switch wtMethod {
	case WTRead:
		if wtType == WTServer {
			timeout = serverConfig.ReadTimeout
		} else {
			timeout = clientConfig.ReadTimeout
		}
	case WTInstall:
		if wtType == WTServer {
			timeout = serverConfig.InstallTimeout
		} else {
			timeout = clientConfig.InstallTimeout
		}
	case WTWrite:
		if wtType == WTServer {
			timeout = serverConfig.WriteTimeout
		} else {
			timeout = clientConfig.WriteTimeout
		}
	}
	return
}

// GetTokenList : get TokenList from server config by wtMethod
func GetTokenList(wtMethod WTMethod) (tokenList []string) {
	switch wtMethod {
	case WTRead:
		tokenList = serverConfig.ReadToken
	case WTInstall:
		tokenList = serverConfig.InstallToken
	case WTWrite:
		tokenList = serverConfig.WriteToken
	}

	return
}

func GetToken(wtMethod WTMethod) (token string) {
	switch wtMethod {
	case WTRead:
		token = clientConfig.ReadToken
	case WTInstall:
		token = clientConfig.InstallToken
	case WTWrite:
		token = clientConfig.WriteToken
	}
	return
}

func GetServerAddr(wtType WTType) (addr string) {
	switch wtType {
	case WTServer:
		addr = serverConfig.Server
	case WTClient:
		addr = clientConfig.Server
	}

	return
}
