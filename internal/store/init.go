package store

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

// global metaData variant
var metaData = &model.MetaData{
	DataMap: map[string]*model.Package{},
	TagMap:  map[string][]string{},
}

// SyncMetaDataFromDisk which tag all packages as temp group
func buildMetaDataFromDisk() {
	dir, err := os.ReadDir(config.DataDir)
	if err != nil {
		slog.Error(
			"fail to read data dir",
			slog.String("dir", config.DataDir),
			slog.Any("err", err),
		)

		return
	}

	metaData = &model.MetaData{
		DataMap: map[string]*model.Package{},
		TagMap:  map[string][]string{},
	}
	var defaultPack *model.Package
	for _, entry := range dir {
		info, _ := entry.Info()

		if info.IsDir() {
			continue
		}

		defaultPack = getDefaultPackage(info)
		metaData.DataMap[info.Name()] = defaultPack
		metaData.TagMap[config.DefaultTagTemp] = append(metaData.TagMap[config.DefaultTagTemp], info.Name())
	}

	byteData, _ := json.MarshalIndent(metaData, "", "	")

	err = writeMetaData(byteData)
	if err != nil {
		slog.Error(
			"fail to write the metaData",
			slog.Any("err", err),
		)
	}
}

func InitData(manualPath string) (err error) {
	slog.Info("init data with manaully config the data path", "dataPath", manualPath)
	err = config.InitDataDir(manualPath)
	if err != nil {
		return fmt.Errorf("initData fail : %w", err)
	}

	// there is no need to check the dir which has done by InitDataDir
	// metaDataPath := config.MetaDataPath
	// metaDataDir := filepath.Dir(metaDataPath)
	//
	// err = os.MkdirAll(metaDataDir, 0o755)
	// if err != nil {
	// 	return err
	// }

	byteData, err := os.ReadFile(config.MetaDataPath)
	slog.Debug("check from MetaDataPath byte data and err", "path", config.MetaDataPath, "byteData", string(byteData), "err", err)

	// if the metadata.json is not exsit which mean this dir is not used before
	// than init the metadata.json by tag all packages with temp
	if os.IsNotExist(err) {
		buildMetaDataFromDisk()
		slog.Info("meta_data file not found, create new one", "filePath", config.MetaDataPath)
		slog.Warn("all packages will be taged as temp")

		err = nil
		return
	}

	// if the meta data exist which means this dir has be used
	// sync lastest file info to meta data
	SyncMetaDataFromDisk()
	err = syncMetaDataToFile(metaData)
	if err != nil {
		slog.Error("fail to sync memory data to file", "filePath", config.MetaDataPath, "err", err)
		return
	}
	err = json.Unmarshal(byteData, metaData)
	if err != nil {
		slog.Error(
			"fail to unmarshal metaData when load json file",
			"err", err.Error(),
		)
		panic(err)
	}

	slog.Info("load meta data successfully ", "packageDir", config.DataDir, "metaDataPath", config.MetaDataPath)

	return
}
