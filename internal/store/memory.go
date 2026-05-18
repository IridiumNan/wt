// Package store : manage the meta file data in the memory and sync it to file when updated
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

var MapLock sync.RWMutex

// SearchPackage : return the pkg if find pkg which matchs the pattern, and return nil if not found
func SearchPackage(pattern string) (results []*model.Package) {
	results = []*model.Package{}
	for pkgName, pkg := range metaData.DataMap {
		if strings.Contains(strings.ToLower(pkgName), strings.ToLower(pattern)) {
			results = append(results, pkg)
		}
	}
	if len(results) == 0 {
		return nil
	}
	return
}

// GetPackageInfo : get byte info of specific pkg
func GetPackageInfo(pkgName string) (byteData []byte, err error) {
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
func UpdateTag(pkgName string, newTag string) (err error) {
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

	err = syncMetaDataToFile(metaData)
	return
}

// AddPackage : add single package
func AddPackage(info fs.FileInfo) {
	MapLock.Lock()
	defer MapLock.Unlock()

	if oldPkg, exist := metaData.DataMap[info.Name()]; exist {
		oldTag := oldPkg.Tag
		if nameList, ok := metaData.TagMap[oldTag]; ok {
			var newNameList []string
			for _, n := range nameList {
				if n != oldPkg.Name {
					newNameList = append(newNameList, n)
				}
			}
			metaData.TagMap[oldTag] = newNameList
		}
	}

	metaData.DataMap[info.Name()] = getDefaultPackage(info)
	metaData.TagMap[config.DefaultTagTemp] = append(metaData.TagMap[config.DefaultTagTemp], info.Name())

	err := syncMetaDataToFile(metaData)
	if err != nil {
		panic(err)
	}
}

// RenamePackage : rename single package without disk operation
func RenamePackage(oldName string, newName string) (err error) {
	// modify the meateData
	MapLock.Lock()
	defer MapLock.Unlock()

	pkg := metaData.DataMap[oldName]
	pkg.Name = newName
	metaData.DataMap[newName] = pkg

	oldNameList := metaData.TagMap[pkg.Tag]

	for i := range oldNameList {
		if oldNameList[i] == oldName {
			newNameList := append(oldNameList[:i], oldNameList[i+1:]...)
			newNameList = append(newNameList, newName)
			metaData.TagMap[pkg.Tag] = newNameList
			break
		}
	}

	delete(metaData.DataMap, oldName)

	err = syncMetaDataToFile(metaData)
	return
}

// DeletePackageByName : delete single package
func DeletePackageByName(pkgName string) (err error) {
	// delet real pkg
	pkgPath := filepath.Join(config.DataDir, pkgName)
	err = os.Remove(pkgPath)
	if err != nil {
		slog.Debug("error when remove file", "func", "memory.DeletePackageByName", "error", err)
		return
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

	err = syncMetaDataToFile(metaData)
	if err != nil {
		panic(err)
	}

	return
}

// ListPackagesByTag : return the NameList by tag
func ListPackagesByTag(targetTag string) (nameList []string) {
	MapLock.RLock()
	defer MapLock.RUnlock()

	nameList = metaData.TagMap[targetTag]
	slog.Debug("check tagMap", "map", fmt.Sprint(metaData.TagMap))

	slog.Debug("List package by tag", "tag", targetTag, "found_res", nameList)
	if len(nameList) == 0 {
		return nil
	}
	return
}

// DeletePackageByTag : delete tag group
func DeletePackageByTag(metaData *model.MetaData, tagName string) {
	MapLock.RLock()

	fileList := metaData.TagMap[tagName]
	MapLock.RUnlock()

	for i := range fileList {
		DeletePackageByName(fileList[i])
	}
}
