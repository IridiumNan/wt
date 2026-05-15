// Package config/client : init and read the client file
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/querypresets"
)

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

func writeClientConfig(configPath string, serverConfig *model.ClientConfig) (err error) {
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

func LoadClientConfig() (clientConfig *model.ClientConfig, err error) {
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
		err = writeClientConfig(configPath, clientConfig)
		if err != nil {
			panic(err)
		}
		return
	}

	json.Unmarshal(data, clientConfig)

	fmt.Println("load config from ", configPath, " as below:")
	fmt.Println(string(data))
	return
}
