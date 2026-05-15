// Package config : init , load, handle the config
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/querypresets"
)

var scanner = bufio.NewScanner(os.Stdin)

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

func writeServerConfig(configPath string, serverConfig *model.ServerConfig) (err error) {
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

func LoadServerConfig() (serverConfig *model.ServerConfig, err error) {
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
		err = writeServerConfig(configPath, serverConfig)
		if err != nil {
			panic(err)
		}
	}

	json.Unmarshal(data, serverConfig)

	fmt.Println("load config from ", configPath, " as below:")
	fmt.Println(string(data))
	return
}
