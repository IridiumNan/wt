package store

import (
	"encoding/json"
	"log/slog"
	"os"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/commonpresets"
)

// writeMetaData : atomically the byteData to MetaDataPath
func writeMetaData(byteData []byte) (err error) {
	tempPath := commonpresets.MetaDataFile + ".tmp"
	tempFile, err := os.OpenFile(
		tempPath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		slog.Error("fail to open tempFIle", "tempFile", tempFile, "err", err)
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

	err = os.Rename(tempPath, commonpresets.MetaDataFile)
	if err != nil {
		panic(err)
	}

	return
}

// sync the meta data to file
func syncMetaData(metaData *model.MetaData) (err error) {
	byteData, _ := json.MarshalIndent(metaData, "", "	")

	err = writeMetaData(byteData)

	return
}
