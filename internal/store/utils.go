package store

import (
	"bufio"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
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

func IsPackageExist(fileName string, nameList []string) bool {
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
	// Can't get current serve dir so metaData remove is imposible
	// Maybe can record all metadata dir when build metaData evevy single time
	// metaData file
	// metaDataDir := filepath.Base(config.MetaDataPath)
	// if approve("remove the meta data (for tag info store) ? [y/N] ->") {
	// 	err = os.RemoveAll(metaDataDir)
	// 	slog.Info("remove the metaDataDir", "dir", metaDataDir)
	// } else {
	// 	fmt.Println("you can check and reuse or remove it manually -> ", config.MetaDataPath)
	// }
	clientConfigPath, _ := config.GetClientConfigPath()
	serverConfigPath, _ := config.GetServerConfigPath()

	if approve("remove the client_config? -> " + clientConfigPath + "[y/N] ->") {
		err = os.Remove(clientConfigPath)
	} else {
		fmt.Println("you can check and reuse or remove it manually -> ", config.ClientConfigPath)
	}

	if approve("remove the server_config.json? -> " + serverConfigPath + " [y/N] ->") {
		err = os.Remove(serverConfigPath)
	} else {
		fmt.Println("you can check and reuse or remove it manually -> ", config.ServerConfigPath)
	}

	if approve("remove the log ? -> " + config.GetLogDir() + " [y/N] ->") {
		err = os.RemoveAll(config.GetLogDir())
	} else {
		fmt.Println("you can check and reuse or remove it manually ->", config.GetLogDir())
	}

	return
}
