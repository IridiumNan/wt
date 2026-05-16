// Package config : init , load, handle the config
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/querypresets"
)

var scanner = bufio.NewScanner(os.Stdin)

var serverConfig *model.ServerConfig

// queryServerConfig : generate a serverconfig by query and called by LoadServerConfig
func queryServerConfig() (serverConfig *model.ServerConfig) {
	serverConfig = &model.ServerConfig{
		Server: queryHostPort(querypresets.ServerHostPortQuery),

		ReadTimeout:    queryTimeout(querypresets.ServerReadTimeoutQuery),
		InstallTimeout: queryTimeout(querypresets.ServerInstallTimeoutQuety),
		WriteTimeout:   queryTimeout(querypresets.ServerWriteTimeoutQuery),

		ReadToken:    queryServerToken(querypresets.ServerReadTokenQuery),
		InstallToken: queryServerToken(querypresets.ServerInstallTokenQuery),
		WriteToken:   queryServerToken(querypresets.ServerWriteTokenQuery),
	}

	return
}

// writeServerConfig : write server config file to configPath
func writeServerConfig(configPath string) (err error) {
	configData, jsonErr := json.MarshalIndent(*serverConfig, "", "	")

	if jsonErr != nil {
		panic(jsonErr)
	}

	configFile, fileErr := os.OpenFile(
		configPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o644,
	)

	if fileErr != nil {
		err = fileErr
		fmt.Println("file open error:", err)
		return
	}
	defer configFile.Close()
	_, writeErr := configFile.Write(configData)
	if writeErr != nil {
		panic(writeErr)
	}

	fmt.Println(">>>> you can check or alter your config -> ", configPath, " <<<<")

	return
}

func InitServerConfig() (err error) {
	serverConfig, err = loadServerConfig()
	return
}

func loadServerConfig() (serveserverConfig *model.ServerConfig, err error) {
	configPath, err := getServerConfigPath()
	if err != nil {
		return
	}

	configDir := filepath.Dir(configPath)
	if err = os.MkdirAll(configDir, 0o755); err != nil {
		return
	}

	// read the config file
	data, err := os.ReadFile(configPath)
	serverConfig = &model.ServerConfig{}
	if os.IsNotExist(err) {
		serverConfig = queryServerConfig()
		err = writeServerConfig(configPath)
		if err != nil {
			panic(err)
		}
	}

	json.Unmarshal(data, serverConfig)

	fmt.Println("load config from ", configPath, " as below:")
	fmt.Println(string(data))
	return
}

// GetServerHostPortFromServer : return the host:port from the server config
func GetServerHostPortFromServer() (server string) {
	return serverConfig.Server
}

// GetServerReadTimeout : return the readTimeout from the server config
func GetServerReadTimeout() (readTimeout time.Duration) {
	return serverConfig.ReadTimeout
}

// GetServerInstallTimeout : return the install Timeout from the server config
func GetServerInstallTimeout() (installTimeout time.Duration) {
	return serverConfig.InstallTimeout
}

// GetServerWriteTimeout : return the write Timeout from the server config
func GetServerWriteTimeout() (writeTimeout time.Duration) {
	return serverConfig.WriteTimeout
}

// GetServerReadTokenList : return the read token list from the server config
func GetServerReadTokenList() (tokenList []string) {
	return serverConfig.ReadToken
}

// GetServerInstallTokenList : return the server token list from the server config
func GetServerInstallTokenList() (tokenList []string) {
	return serverConfig.InstallToken
}

// GetServerWriteTokenList : return the write token list from the server config
func GetServerWriteTokenList() (tokenList []string) {
	return serverConfig.WriteToken
}
