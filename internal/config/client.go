// Package config/client : init and read the client file
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/querypresets"
)

var clientConfig *model.ClientConfig

func queryClientConfig() (clientConfig *model.ClientConfig) {
	clientConfig = &model.ClientConfig{
		Server: queryHostPort(querypresets.ClientHostPortQuery),

		ReadTimeout:    queryTimeout(querypresets.ClientReadTimeoutQuery),
		InstallTimeout: queryTimeout(querypresets.ClientInstallTimeoutQuery),
		WriteTimeout:   queryTimeout(querypresets.ClientWriteTimeoutQuery),

		ReadToken:    queryClientToken(querypresets.ClientReadTokenQuery),
		InstallToken: queryClientToken(querypresets.ClientInstallTokenQuery),
		WriteToken:   queryClientToken(querypresets.ClientWriteTokenQuery),
	}
	return
}

func writeClientConfig(configPath string) (err error) {
	configData, jsonErr := json.MarshalIndent(*clientConfig, "", "	")

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

// InitClientConfig : init the client config
func InitClientConfig() (err error) {
	clientConfig, err = loadClientConfig()
	return
}

// loadClientConfig : the client config from file
func loadClientConfig() (clientConfig *model.ClientConfig, err error) {
	configPath, err := getClientConfigPath()
	if err != nil {
		return
	}

	configDir := filepath.Dir(configPath)
	if err = os.MkdirAll(configDir, 0o755); err != nil {
		return
	}

	// read the config file
	data, err := os.ReadFile(configPath)
	clientConfig = &model.ClientConfig{}
	if os.IsNotExist(err) {
		fmt.Printf("client config file not exsit, begin to init...\n\n")
		clientConfig = queryClientConfig()
		err = writeClientConfig(configPath)
		if err != nil {
			panic(err)
		}
		return
	}

	err = json.Unmarshal(data, clientConfig)
	if err != nil {
		fmt.Println("err when Unmarshal", err)
		panic(err)
	}

	fmt.Println("load config from ", configPath, " as below:")
	fmt.Println(string(data))
	return
}

// GetServerHostPortFromClient : return the host:port from the client config
func GetServerHostPortFromClient() (server string) {
	return clientConfig.Server
}

// GetClientReadTimeout : return the readTimeout from the client config
func GetClientReadTimeout() (readTimeout time.Duration) {
	return clientConfig.ReadTimeout
}

// GetClientInstallTimeout : return the install Timeout from the client config
func GetClientInstallTimeout() (installTimeout time.Duration) {
	return clientConfig.InstallTimeout
}

// GetClientWriteTimeout : return the write Timeout from the client config
func GetClientWriteTimeout() (writeTimeout time.Duration) {
	return clientConfig.WriteTimeout
}

// GetClientReadToken : return the read token  from the client config
func GetClientReadToken() (token string) {
	return clientConfig.ReadToken
}

// GetClientInstallToken : return the install token from the client config
func GetClientInstallToken() (token string) {
	return clientConfig.InstallToken
}

// GetClientWriteToken : return the write token  from the client config
func GetClientWriteToken() (token string) {
	return clientConfig.WriteToken
}
