package command

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
	FlagIdx      = 1
	FirstTarget  = 2
	SecondTarget = 3
	ThirdTarget  = 4
)

func handleServerZeroTarget(args []string) (err error) {
	cmdFlag := args[FlagIdx]
	switch cmdFlag {
	case "log":
		// readLog
		err = loghelper.ReadServerLog()
		Done = true
	default:
		err = fmt.Errorf("unknow command : %s", cmdFlag)
		Done = true

	}

	return
}

func handleServerOneTarget(cmdFlag string, args []string) (err error) {
	switch cmdFlag {
	case "-d", "dir":
		err = store.InitData(args[FirstTarget])
		Done = false
	case "-c", "config":
		if args[FirstTarget] == "show" {
			config.ServerConfigShow()
			Done = true
		}
	default:
		err = fmt.Errorf("unknow command : %s", cmdFlag)
		Done = true
	}
	return
}

func execCommand(args []string) (err error) {
	cmdFlag := args[FlagIdx]

	switch len(args) {
	case 2:
		err = handleServerZeroTarget(args)
	case 3:
		err = handleServerOneTarget(cmdFlag, args)
	case 5:
		if cmdFlag == "config" || cmdFlag == "-c" {
			err = config.AlterServerConfig(args[FirstTarget], args[SecondTarget], args[ThirdTarget])
			if err != nil {
				fmt.Println("restart the server to make config reload")
			}
			Done = true
		}
	}

	return
}

func ServerMain(args []string, debug bool) {
	err := loghelper.InitServerLogger(debug)
	if err != nil {
		panic(err)
	}

	err = config.InitServerConfig()
	if err != nil {
		fmt.Println("server config init err:", err)
		panic(err)
	}

	// this args contains server as args[0]
	if len(args) > 1 {
		err = execCommand(args)
	}
	if err != nil {
		slog.Error("fail to exec server command: " + err.Error())
		return
	}
	if Done {
		return
	}

	if len(args) == 1 {
		err = store.InitData("")
		if err != nil {
			slog.Error(err.Error())
			return
		}
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
