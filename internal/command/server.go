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

func execTagTokenCommand(args []string) (err error) {
	tag := args[SecondTargetIndex]
	TokenJSONKey := args[ThirdTargetIndex]
	operation := args[FourthTargetIndex]
	token := args[FifthTargetIndex]

	switch operation {
	case "add":
		switch TokenJSONKey {
		case config.ReadTokenJSONKey:
			config.AddTagToken(tag, token, model.WTRead)
		case config.InstallTokenJSONKey:
			config.AddTagToken(tag, token, model.WTInstall)
		case config.WriteTokenJSONKey:
			config.AddTagToken(tag, token, model.WTWrite)
		}
		fmt.Printf("add the %s -> %s for tag %s", TokenJSONKey, token, tag)
		fmt.Println("you need to restart the server to make this config reload")
	case "rm":
		return

	}

	config.SyncServerConfig()

	return
}

func newExecCommand(args []string) (err error) {
	command := args[CommandIndex]

	// stript the command
	args = args[1:]

	slog.Debug("exec command in server", "command", command, "args", args)
	switch command {
	case "log":
		err = loghelper.ReadServerLog()
		Done = true
	case "config", "-c":
		if len(args) < 3 {
			switch args[FirstTargetIndex] {
			case "show":
				config.ServerConfigShow()
				Done = true
			case config.ReadTokenJSONKey, config.InstallTokenJSONKey, config.WriteTokenJSONKey:
				err = config.AlterServerConfig(args[FirstTargetIndex], args[SecondTargetIndex], args[ThirdTargetIndex])
				Done = true
			}
		} else if len(args) == 5 && args[FirstTargetIndex] == "tag" {
			err = execTagTokenCommand(args)
			Done = true

		}
	case "dir", "-d":
		err = store.InitData(args[FirstTargetIndex])
		Done = false
	default:
		err = fmt.Errorf("unknown command: %s", command)
		Done = true
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
		fmt.Printf("server config init err: %s\nyou can try to wt uninstall to remove the previous tags data and rebuild it", err.Error())
		return
	}

	// args has been stript the server command, and begin with command like "config" "log"
	if len(args) > 1 {
		err = newExecCommand(args[1:])
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
