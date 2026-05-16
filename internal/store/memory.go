// Package store : manage the meta file data in the memory and sync it to file when updated
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/commonpresets"
)

var MapLock sync.RWMutex

// getDefaultPackage : handle the info the tag the package as temp group
func getDefaultPackage(info fs.FileInfo) (pack *model.Package) {
	pack = &model.Package{
		Name:    info.Name(),
		Tag:     commonpresets.DefaultTagTemp,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	return
}

// SearchPackage : return the pkgName if find pkg which matchs the pattern, and return "" if not found
func SearchPackage(metaData *model.MetaData, pattern string) (pkgName string) {
	for pkg := range metaData.DataMap {
		if strings.Contains(strings.ToLower(pkg), strings.ToLower(pattern)) {
			pkgName = pkg
			return
		}
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

// GetPackageInfo : get byte info of specific pkg
func GetPackageInfo(metaData *model.MetaData, pkgName string) (byteData []byte, err error) {
	MapLock.RLock()
	pkgInfo, ok := metaData.DataMap[pkgName]
	MapLock.RUnlock()

	if !ok {
		err = errors.New("pkg not found")
		return
	}

	byteData, err = json.Marshal(pkgInfo)
	if err != nil {
		panic(err)
	}
	return
}

// UpdateTag : update the tag
func UpdateTag(metaData *model.MetaData, pkgName string, newTag string) (err error) {
	MapLock.Lock()
	defer MapLock.Unlock()

	pkg, ok := metaData.DataMap[pkgName]
	if !ok {
		err = errors.New("pkg not found")
		return
	}

	oldTag := pkg.Tag
	oldNameList := metaData.TagMap[oldTag]
	newNameList := removePkgInNameList(oldNameList, pkgName)

	metaData.TagMap[oldTag] = newNameList
	metaData.TagMap[newTag] = append(metaData.TagMap[newTag], pkgName)

	pkg.Tag = newTag

	err = syncMetaData(metaData)
	return
}

// AddPackage : add single package
func AddPackage(metaData *model.MetaData, info fs.FileInfo) {
	MapLock.Lock()
	defer MapLock.Unlock()
	metaData.DataMap[info.Name()] = getDefaultPackage(info)
	metaData.TagMap[commonpresets.DefaultTagTemp] = append(metaData.TagMap[commonpresets.DefaultTagTemp], info.Name())

	err := syncMetaData(metaData)
	if err != nil {
		panic(err)
	}
}

// RenamePackage : rename single package
func RenamePackage(metaData *model.MetaData, oldName string, newName string) (err error) {
	// rename real pkg
	oldPath := filepath.Join(commonpresets.DataDir, oldName)
	newPath := filepath.Join(commonpresets.DataDir, newName)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("err when rename the file:", err)
	}

	// modify the meateData
	MapLock.Lock()
	defer MapLock.Unlock()

	value := metaData.DataMap[oldName]
	value.Name = newName
	metaData.DataMap[newName] = value

	oldNameList := metaData.TagMap[value.Tag]

	for i := range oldNameList {
		if oldNameList[i] == oldName {
			newNameList := append(oldNameList[:i], oldNameList[i+1:]...)
			metaData.TagMap[value.Tag] = newNameList
			break
		}
	}

	delete(metaData.DataMap, oldName)

	err = syncMetaData(metaData)
	return
}

// DeletePackageByName : delete single package
func DeletePackageByName(metaData *model.MetaData, pkgName string) {
	// delet real pkg
	err := os.Remove(pkgName)
	if err != nil {
		fmt.Println("err when remove file :", err)
	}

	MapLock.Lock()
	defer MapLock.Unlock()

	oldTag := metaData.DataMap[pkgName].Tag
	oldNameList := metaData.TagMap[oldTag]

	for i := range oldNameList {
		if oldNameList[i] == pkgName {
			newNameList := append(oldNameList[:i], oldNameList[i+1:]...)
			metaData.TagMap[oldTag] = newNameList
			break
		}
	}

	delete(metaData.DataMap, pkgName)

	err = syncMetaData(metaData)
	if err != nil {
		panic(err)
	}
}

// ListPackagesByTag : return the NameList by tag
func ListPackagesByTag(metaData *model.MetaData, tagName string) (nameList []string) {
	MapLock.RLock()
	defer MapLock.RUnlock()

	nameList = metaData.TagMap[tagName]
	return
}

// DeletePackageByTag : delete tag group
func DeletePackageByTag(metaData *model.MetaData, tagName string) {
	MapLock.RLock()

	fileList := metaData.TagMap[tagName]
	MapLock.RUnlock()

	for i := range fileList {
		DeletePackageByName(metaData, fileList[i])
	}
}

// InitMetaData which tag all packages as temp group
func initMetaData() (metaData *model.MetaData) {
	dir, err := os.ReadDir(commonpresets.DataDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	metaData = &model.MetaData{}
	var defaultPack *model.Package
	for _, entry := range dir {
		info, _ := entry.Info()

		defaultPack = getDefaultPackage(info)
		metaData.DataMap[info.Name()] = defaultPack
		metaData.TagMap[commonpresets.DefaultTagTemp] = append(metaData.TagMap[commonpresets.DefaultTagTemp], info.Name())
	}

	byteData, _ := json.MarshalIndent(metaData, "", "	")

	err = writeMetaData(byteData)
	if err != nil {
		fmt.Println("file to write the metaData :", err)
	}

	return
}

// LoadMetaData : load the meta data for file
func LoadMetaData() (metaData *model.MetaData) {
	dataPath := commonpresets.MetaDataPath

	byteData, err := os.ReadFile(dataPath)

	if os.IsNotExist(err) {
		metaData = initMetaData()
		return
	}

	err = json.Unmarshal(byteData, metaData)
	if err != nil {
		fmt.Println("fail to Unmarshal")
		panic(err)
	}

	fmt.Println("load meta data from ", dataPath, " sucessfully")
	fmt.Println()
	return
}
