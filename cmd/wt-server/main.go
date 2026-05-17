package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/server"
	"gitee.com/cai-zixiang_hainan/wt/internal/store"
	"gitee.com/cai-zixiang_hainan/wt/pkg/loghelper"
)

const (
	DEBUG   = true
	VERSION = "0.0.1"
)

func init() {
	loghelper.JSONLog(DEBUG)
}

func main() {
	err := config.InitServerConfig()
	if err != nil {
		fmt.Println("server config init err:", err)
		panic(err)
	}

	err = store.InitMetaData()
	if err != nil {
		panic(err)
	}

	serveAddr := config.GetServerAddr(model.WTServer)
	readTimeout := config.GetTimeout(model.WTServer, model.WTRead)
	writeTimeout := config.GetTimeout(model.WTServer, model.WTWrite)

	router := server.NewRouter()

	httpServer := &http.Server{
		Addr:         serveAddr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  time.Second * 120,
	}

	slog.Debug("server start", "addr", serveAddr)
	slog.Debug("load timeout config", "readtimeout", readTimeout, "writetimeout", writeTimeout)

	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Println("serve fail :", err)

		return
	}
}
