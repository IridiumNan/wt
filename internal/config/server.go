// Package config : init , load, handle the config
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/querypresets"
)

var scanner = bufio.NewScanner(os.Stdin)

// query and return server
func queryHostPort(query model.Query) (server string) {
	fmt.Println(query.Head)
	fmt.Println("example:", query.Example)
	fmt.Println("default value:", query.Default)
	fmt.Printf("enter your value[default if blank] ->")
	scanner.Scan()
	server = strings.TrimSpace(scanner.Text())

	if server == "" {
		server = query.Default
	}
	fmt.Println()
	return
}

// set timeout with query input like "enter the Install timeout for client:"
// return the timeOut which is time.Duration
func queryTimeout(query model.Query) (timeOut time.Duration) {
	timeString := ""

	fmt.Println(query.Head)
	fmt.Println("example:", query.Example)
	fmt.Println("default value:", query.Default)
	fmt.Printf("enter your value[default if blank] ->")
	scanner.Scan()
	timeString = strings.TrimSpace(scanner.Text())

	if timeString == "" {
		timeString = query.Default
	}

	timeOut, err := time.ParseDuration(timeString)
	if err != nil {
		fmt.Printf("error when parse %s", timeString)
		fmt.Println(err)
		timeOut = time.Second * 0
	}

	fmt.Println()
	return
}

// set Servertokens with query like ">>set the read tokens for server<<"
func queryServerToken(query model.Query) (tokens []string) {
	var tokenNum int
	var tempToken string

	fmt.Println(query.Head)
	fmt.Println("example:", query.Example)
	fmt.Println("default value:", query.Default)

	fmt.Printf("enter the number of tocken you want to add:")
	scanner.Scan()
	numStr := strings.TrimSpace(scanner.Text())

	if numStr == "" {
		tokens = []string{}
		return
	}

	_, err := fmt.Sscanf(numStr, "%d", &tokenNum)
	if err != nil || tokenNum < 0 {
		fmt.Println("invalid number, skipping token configuration")
		fmt.Println()
		return []string{}
	}
	fmt.Println("setting your value[default if blank]")
	for i := 0; i < tokenNum; i++ {
		fmt.Printf("enter the %d th token ->", i+1)
		scanner.Scan()
		tempToken = strings.TrimSpace(scanner.Text())

		tokens = append(tokens, tempToken)
	}
	fmt.Println()
	return
}

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
	n, writeErr := configFile.Write(configData)
	if writeErr != nil {
		panic(writeErr)
	}
	fmt.Println(n)

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
	if os.IsNotExist(err) {
		serverConfig = queryServerConfig()
		err = writeServerConfig(configPath, serverConfig)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(string(data))
	fmt.Println("load config as below:")
	fmt.Println("load the config from ", configPath)
	return
}
