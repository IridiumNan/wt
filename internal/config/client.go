// Package config/client : init and read the client file
package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/querypresets"
)

var clientConfig *model.ClientConfig

func queryClientConfig() {
	clientConfig = &model.ClientConfig{
		Server: queryHostPort(querypresets.ClientHostPortQuery),

		ReadTimeout:    queryTimeout(querypresets.ClientReadTimeoutQuery),
		InstallTimeout: queryTimeout(querypresets.ClientInstallTimeoutQuery),
		WriteTimeout:   queryTimeout(querypresets.ClientWriteTimeoutQuery),

		ReadToken:    queryClientToken(querypresets.ClientReadTokenQuery),
		InstallToken: queryClientToken(querypresets.ClientInstallTokenQuery),
		WriteToken:   queryClientToken(querypresets.ClientWriteTokenQuery),
	}
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
	err = loadClientConfig()
	return
}

// loadClientConfig : the client config from file
func loadClientConfig() (err error) {
	configPath, err := GetClientConfigPath()
	if err != nil {
		return
	}

	configDir := filepath.Dir(configPath)
	if err = os.MkdirAll(configDir, 0o755); err != nil {
		return
	}

	// read the config file
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		fmt.Println("expected config file: ", configPath)
		fmt.Printf("client config file not exsit, begin to init...\n\n")
		queryClientConfig()
		fmt.Println("config content", clientConfig)
		err = writeClientConfig(configPath)
		if err != nil {
			panic(err)
		}
		return
	}

	clientConfig = &model.ClientConfig{}
	err = json.Unmarshal(data, clientConfig)
	if err != nil {
		fmt.Println("err when Unmarshal", err)
		panic(err)
	}

	slog.Debug("load config from" + configPath)
	slog.Debug(string(data))
	return
}

func AlterClientConfig(key string, value string) (err error) {
	switch key {
	case "server":
		clientConfig.Server = value
	case "read_token":
		clientConfig.ReadToken = value
	case "install_token":
		clientConfig.InstallToken = value
	case "write_token":
		clientConfig.WriteToken = value
	default:
		err = fmt.Errorf("unknow element or element not allow to alter :%s", key)
		return
	}

	configPath, _ := GetClientConfigPath()
	err = writeClientConfig(configPath)
	if err != nil {
		return
	}

	fmt.Println("alter the ", key, " -> ", value, " successfully")

	return
}

func ClientConfigShow() {
	fmt.Println("-------------- client config ------------------")
	clientCOnfigPath, _ := GetClientConfigPath()
	fmt.Println("config file :", clientCOnfigPath)
	fmt.Println("-------------- config content -----------------")
	fmt.Println("default server:", GetServerAddr(model.WTClient))
	fmt.Println("read_token:", GetToken(model.WTRead))
	fmt.Println("install_token:", GetToken(model.WTInstall))
	fmt.Println("write_token:", GetToken(model.WTWrite))
	fmt.Println("------------------------------------------------")

	fmt.Println("you can config any element by Usage: wt config element value")
	fmt.Println("example: wt config read_token fowehgiqwojfaoinhgvaoij")
}
