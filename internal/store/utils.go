package store

import (
	"bufio"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

var scanner = bufio.NewScanner(os.Stdin)

// getDefaultPackage : handle the info the tag the package as temp group
func getDefaultPackage(info fs.FileInfo) (pack *model.Package) {
	pack = &model.Package{
		Name:    info.Name(),
		Tag:     config.DefaultTagTemp,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	return
}

func removePkgInNameList(oldNameList []string, pkgName string) (newNameList []string) {
	for i := range oldNameList {
		if oldNameList[i] == pkgName {
			newNameList = append(oldNameList[:i], oldNameList[i+1:]...)
			return
		}
	}
	return
}

func isPackageExist(fileName string, nameList []string) bool {
	return slices.Contains(nameList, fileName)
}

func approve(query string) bool {
	fmt.Print(query)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	if choice == "" {
		choice = "n"
	}

	slog.Info("remove the file and check", "query", query, "choice", choice)
	return strings.ToLower(choice) == "y"
}

func Uninstall() (err error) {
	// metaData file
	metaDataDir := filepath.Base(config.MetaDataPath)

	if approve("remove the meta data (for tag info store) ? [y/N] ->") {
		err = os.RemoveAll(metaDataDir)
		slog.Info("remove the metaDataDir", "dir", metaDataDir)
	} else {
		fmt.Println("you can check and reuse or remove it manually -> ", config.MetaDataPath)
	}

	if approve("remove the client_config? [y/N] ->") {
		err = os.Remove(config.ClientConfigPath)
		slog.Info("remove the client config file", "filePath", config.ClientConfigPath)

	} else {
		fmt.Println("you can check and reuse or remove it manually -> ", config.ClientConfigPath)
	}

	if approve("remove the server_config.json? [y/N] ->") {
		err = os.Remove(config.ServerConfigPath)
		slog.Info("remove the server config file", "filePath", config.ServerConfigPath)
	} else {
		fmt.Println("you can check and reuse or remove it manually -> ", config.ServerConfigPath)
	}

	if approve("remove the log ? [y/N] ->") {
		err = os.RemoveAll(config.GetLogDir())
		slog.Info("remove the log dir", "filePath", config.GetLogDir())
	} else {
		fmt.Println("you can check and reuse or remove it manually ->", config.GetLogDir())
	}

	return
}
