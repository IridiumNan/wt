package config

import (
	"os"
	"path/filepath"

	"gitee.com/cai-zixiang_hainan/wt/internal/presets/commonpresets"
)

func getServerConfigPath() (ServerConfigPath string, err error) {
	home, err := os.UserHomeDir()

	ServerConfigPath = filepath.Join(home, commonpresets.ServerConfigPath)

	return
}

func getClientConfigPath() (ClientConfigPath string, err error) {
	home, err := os.UserHomeDir()

	ClientConfigPath = filepath.Join(home, commonpresets.ClientConfigPath)

	return
}
