package main

import (
	"fmt"
	"net/http"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
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

	serveAddr := config.GetServerAddr(config.WTServer)
	readTimeout := config.GetTimeout(config.WTServer, config.WTRead)
	writeTimeout := config.GetTimeout(config.WTServer, config.WTWrite)

	router := server.NewRouter()

	httpServer := &http.Server{
		Addr:         serveAddr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  time.Second * 120,
	}

	fmt.Println("server start in addr: ", serveAddr)
	fmt.Println("Read Timeout config: ", readTimeout)
	fmt.Println("write Timeout config: ", writeTimeout)

	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Println("serve fail :", err)

		return
	}
}
