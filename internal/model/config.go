package model

import (
	"time"
)

// ClientConfig struct: store the config of client which load from the ClientConfigPath
type ClientConfig struct {
	// host:port
	Server string `json:"server"`

	// timeout for read
	ReadTimeout time.Duration `json:"read_timeout"`
	// timeout for install
	InstallTimeout time.Duration `json:"install_timeout"`
	// timeout for write
	WriteTimeout time.Duration `json:"write_timeout"`

	// auth token for read
	ReadToken string `json:"read_token"`
	// auth token for install
	InstallToken string `json:"install_token"`
	// auth token for write
	WriteToken string `json:"write_token"`
}

type TagAuthTokens struct {
	ReadTokens    []string `json:"read_token"`
	Installtokens []string `json:"install_token"`
	WriteTokens   []string `json:"write_token"`
}

// ServerConfig struct: store the config of server which load from the ServerConfigPath
type ServerConfig struct {
	// host:port
	Server string `json:"server"`

	// timeout for read
	ReadTimeout time.Duration `json:"read_timeout"`
	// timeout for install
	InstallTimeout time.Duration `json:"install_timeout"`
	// timeout for write
	WriteTimeout time.Duration `json:"write_timeout"`

	// auth tokens for read
	ReadToken []string `json:"read_token"`
	// auth tokens for install
	InstallToken []string `json:"install_token"`
	// auth tokens for write
	WriteToken []string `json:"write_token"`

	// tokens for tag
	TagTokenMap map[string]TagAuthTokens `json:"tag_token"`
}

func addIfNotEmpty(tokenList []string, newToken string) []string {
	if newToken != "" {
		tokenList = append(tokenList, newToken)
	}

	return tokenList
}

func (sc *ServerConfig) AddTagTokens(tagName string, readToken string, installToken string, writeToken string) {
	authTokens, ok := sc.TagTokenMap[tagName]
	if !ok {
		sc.TagTokenMap[tagName] = TagAuthTokens{}
	}

	authTokens.ReadTokens = addIfNotEmpty(authTokens.ReadTokens, readToken)
	authTokens.Installtokens = addIfNotEmpty(authTokens.Installtokens, installToken)
	authTokens.WriteTokens = addIfNotEmpty(authTokens.WriteTokens, writeToken)

	sc.TagTokenMap[tagName] = authTokens
}
