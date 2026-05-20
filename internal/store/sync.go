package store

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"os"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

// writeMetaData : atomically the byteData to MetaDataPath
func writeMetaData(byteData []byte) (err error) {
	tempPath := config.MetaDataPath + ".tmp"
	tempFile, err := os.OpenFile(
		tempPath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		slog.Error("fail to open tempFile", "tempFile", tempFile, "err", err)
	}

	_, err = tempFile.Write(byteData)
	if err != nil {
		panic(err)
	}
	closeErr := tempFile.Close()

	if closeErr != nil {
		slog.Error("tempFile close error", "tempFile", tempFile, "err", err)
		err = closeErr
		panic(err)
	}

	err = os.Rename(tempPath, config.MetaDataPath)
	if err != nil {
		panic(err)
	}

	return
}

// sync the meta data to file
func syncMetaDataToFile(metaData *model.MetaData) (err error) {
	byteData, _ := json.MarshalIndent(metaData, "", "	")

	err = writeMetaData(byteData)

	return
}

// SyncMetaDataFromDisk : sync new package from disk
func SyncMetaDataFromDisk() {
	slog.Debug("sync meta data from disk")
	dir, err := os.ReadDir(config.DataDir)
	if err != nil {
		slog.Error(
			"fial to read data dir",
			"dir", config.DataDir,
			"err", err,
		)
		return
	}

	diskPackages := make(map[string]fs.FileInfo)
	for _, entry := range dir {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.IsDir() {
			continue
		}

		diskPackages[info.Name()] = info
	}
	for pkgName := range metaData.DataMap {
		// remove package from memory which don't exsit in the disk
		_, inDisk := diskPackages[pkgName]
		if !inDisk {
			pkg := metaData.DataMap[pkgName]
			delete(metaData.DataMap, pkgName)

			if _, exist := metaData.TagMap[pkg.Tag]; exist {
				newNameList := removePkgInNameList(metaData.TagMap[pkg.Tag], pkgName)
				metaData.TagMap[pkg.Tag] = newNameList
			}

			slog.Info("remove a package from memory when sync data from disk", "pkgName", pkgName)
		}
	}

	for pkgName, info := range diskPackages {
		if _, exist := metaData.DataMap[pkgName]; !exist {
			AddPackage(info)
			slog.Info("add a package from the disk", "pkgName", info.Name(), "ModeTime", info.ModTime(), "size", info.Size())
		}
	}

	err = syncMetaDataToFile(metaData)
	if err != nil {
		slog.Error("SyncMetaDataFromDisk error when sync to file ", "error", err)
	}
}
