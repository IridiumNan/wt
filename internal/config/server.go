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
	configPath, pathErr := GetServerConfigPath()
	if pathErr != nil {
		err = pathErr
		return
	}

	configDir := filepath.Dir(configPath)
	slog.Debug("mk the config dir", "dir", configDir, "configPath", configPath)
	if err = os.MkdirAll(configDir, 0o755); err != nil {
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

	slog.Debug("load config...", "config Path", configPath, "config content", string(data))
	return
}

func ServerConfigShow() {
	fmt.Println("----------------- server config ---------------")
	configPath, _ := GetServerConfigPath()
	fmt.Println("config file: ", configPath)
	fmt.Println("----------------- config content ---------------")
	fmt.Println("server: ", serverConfig.Server)
	fmt.Println("read_token: ", serverConfig.ReadToken)
	fmt.Println("install_token: ", serverConfig.InstallToken)
	fmt.Println("write_token: ", serverConfig.WriteToken)
	fmt.Println("-------------------------------------------------")
	fmt.Println("you can alter the config manually by edit the configPath and restart the server")
}

func AddToken(token string, wtMethod model.WTMethod) (err error) {
	switch wtMethod {
	case model.WTRead:
		serverConfig.ReadToken = append(serverConfig.ReadToken, token)
	case model.WTInstall:
		serverConfig.InstallToken = append(serverConfig.InstallToken, token)
	case model.WTWrite:
		serverConfig.WriteToken = append(serverConfig.WriteToken, token)
	}

	configPath, _ := GetServerConfigPath()
	err = writeServerConfig(configPath)

	return
}

func AlterServerConfig(object string, operation string, value string) (err error) {
	if operation == "add" {
		switch object {
		case "read_token":
			err = AddToken(value, model.WTRead)
			if err == nil {
				fmt.Println("add read_token :", value)
			}
		case "install_token":
			err = AddToken(value, model.WTInstall)
			if err == nil {
				fmt.Println("add install_token :", value)
			}
		case "write_token":
			err = AddToken(value, model.WTWrite)
			if err == nil {
				fmt.Println("add write_token :", value)
			}
		}
	}

	return
}
