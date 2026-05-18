package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultTagTemp   = "temp"
	DefaultTagStatic = "static"
	ClientConfigPath = "/.config/water-repo/client_config.json"
	ServerConfigPath = "/.config/water-repo/server_config.json"
	metaDataFileName = "meta_data.json"
	LogFile          = "log.txt"
)

var (
	DataDir      = "."
	LogDir       = GetLogDir()
	MetaDataPath string
)

// AppName is the name of the application used for directory paths
const AppName = "water-repo"

// internal/presets/config/path.go

func GetLogDir() string {
	var baseDir string

	switch runtime.GOOS {
	case "darwin": // macOS
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, "Library", "Logs")
	case "windows":
		appData := os.Getenv("LOCALAPPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, "AppData", "Local")
		}
		baseDir = filepath.Join(appData, AppName, "Logs")
	default: // Linux & Others
		// 优先使用 XDG_STATE_HOME，否则回退到 ~/.local/state
		xdgStateHome := os.Getenv("XDG_STATE_HOME")
		if xdgStateHome != "" {
			baseDir = xdgStateHome
		} else {
			home, _ := os.UserHomeDir()
			baseDir = filepath.Join(home, ".local", "state")
		}
	}

	return filepath.Join(baseDir, AppName)
}

// InitDataDir : make the data dir and the meta data dir in advance without write any data
func InitDataDir(manualPath string) (err error) {
	if manualPath == "" {
		manualPath = DataDir
	}
	err = os.MkdirAll(manualPath, 0o755)
	if err != nil {
		return fmt.Errorf("error when make data dir: %w", err)
	}

	DataDir = manualPath

	metaDataDir := filepath.Join(DataDir, ".metaData")
	err = os.MkdirAll(metaDataDir, 0o755)

	MetaDataPath = filepath.Join(metaDataDir, metaDataFileName)
	return
}

func getServerConfigPath() (serverConfigPath string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	serverConfigPath = filepath.Join(home, ServerConfigPath)

	return
}

func getClientConfigPath() (clientConfigPath string, err error) {
	home, err := os.UserHomeDir()

	clientConfigPath = filepath.Join(home, ClientConfigPath)

	return
}

func GetServerLogPath() (path string) {
	return filepath.Join(LogDir, LogFile)
}
