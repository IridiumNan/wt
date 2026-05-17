// Package commonpresets store some static or const variant
package commonpresets

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultTagTemp   = "temp"
	DefaultTagStatic = "static"
	DataDir          = "./"
	metaDataFileName = "meta_data.json"
)

// AppName is the name of the application used for directory paths
const AppName = "water-repo"

// GetMetaDataDir returns the appropriate data directory for the current OS
func GetMetaDataDir() string {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\water-repo
		appData := os.Getenv("APPDATA")
		if appData == "" {
			// Fallback if APPDATA is not set
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		baseDir = appData
	case "darwin":
		// macOS: ~/Library/Application Support/water-repo
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, "Library", "Application Support")
	default:
		// Linux and others: ~/.local/share/water-repo
		// Check XDG_DATA_HOME first, fallback to ~/.local/share
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		if xdgDataHome != "" {
			baseDir = xdgDataHome
		} else {
			home, _ := os.UserHomeDir()
			baseDir = filepath.Join(home, ".local", "share")
		}
	}

	return filepath.Join(baseDir, AppName)
}

// MetaDataDir is the global data directory path (initialized at startup)
var MetaDataDir = GetMetaDataDir()

// MetaDataFile is the path to the metadata JSON file
var MetaDataFile = filepath.Join(MetaDataDir, metaDataFileName)
