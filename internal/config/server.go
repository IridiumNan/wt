// Package config : init , load, handle the config
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/querypresets"
)

var scanner = bufio.NewScanner(os.Stdin)

var serverConfig *model.ServerConfig = &model.ServerConfig{}

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

	slog.Debug("write the server config as " + string(configData))

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
		slog.Error("file open error in write server config", "configPath", configPath, "err", err)
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

func loadServerConfig() (config *model.ServerConfig, err error) {
	configPath, pathErr := getServerConfigPath()
	if pathErr != nil {
		err = pathErr
		return
	}

	configDir := filepath.Dir(configPath)
	if mkDirErr := os.MkdirAll(configDir, 0o755); mkDirErr != nil {
		err = mkDirErr
		return
	}

	// read the config file
	data, readErr := os.ReadFile(configPath)
	if os.IsNotExist(readErr) {
		config = queryServerConfig()
		serverConfig = config
		err = writeServerConfig(configPath)
		if err != nil {
			panic(err)
		}

		return
	}

	config = &model.ServerConfig{}

	jsonErr := json.Unmarshal(data, config)

	if jsonErr != nil {
		slog.Error("fail to unmarshal when load server config", "err", jsonErr)
		err = jsonErr
		return
	}

	fmt.Println("load config from ", configPath, " as below:")
	fmt.Println(string(data))
	return
}
