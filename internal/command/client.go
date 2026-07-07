package command

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gitee.com/cai-zixiang_hainan/wt/internal/client"
	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/pkg/loghelper"
)

var Done bool

var scanner = bufio.NewScanner(os.Stdin)

func ClientMain(args []string, debug bool) {
	loghelper.InitClientLogger(debug)

	err := config.InitClientConfig()
	if err != nil {
		fmt.Println("err when init the client config: ", err)
		return
	}

	command := args[CommandIndex]
	args = args[CommandIndex+1:]

	slog.Debug("check command and args for client", "command", command, "args", args)

	switch command {

	case "config", "-c":
		err = configCommand(args)
	case "search":
		err = searchCommand(args)
	case "info":
		err = infoCommand(args)
	case "install":
		err = installCommand(args)
	case "upload":
		err = uploadCommand(args)
	case "mv":
		err = mvCommand(args)
	case "rm":
		err = rmCommand(args)
	case "ls", "list":
		err = listCommand(args)
	case "sync":
		err = syncCommand()
	case "tag":
		err = tagCommand(args)
	case "reload":
		err = reloadCommand()
	case "public":
		err = publicComand(args)
	case "link", "links":
		err = linksCommand()
	case "private":
		err = privateCommand(args)
	case "list-servers":
		listServersCommand()
		err = nil
	case "change-server":
		err = changeServerCommand(args)
	case "add-server":
		err = addserverCommand(args)
	case "del-server":
		err = delServerCommand(args)

	}

	if err != nil {
		fmt.Println("fail to exec ", command, " err :", err.Error())
	}
	if !Done {
		fmt.Println("Bad Usage")
		Usage(command)
	}
}

func configCommand(args []string) (err error) {
	switch len(args) {
	case 0:
		return errors.New("args required")
	case 1:
		if args[FirstTargetIndex] == "show" {
			config.ClientConfigShow()
			Done = true
		}
	case 2:
		err = config.AlterClientConfig(args[FirstTargetIndex], args[SecondTargetIndex])
		Done = true
	}

	return
}

func searchCommand(args []string) (err error) {
	if len(args) == 0 {
		return errors.New("search pattern is required")
	}

	err = client.SearchRequest(args[FirstTargetIndex])
	Done = true

	return
}

func listCommand(args []string) (err error) {
	switch len(args) {
	case 0:
		err = client.ListRequest(config.DefaultTagTemp)
		Done = true
	case 1:
		err = client.ListRequest(args[FirstTargetIndex])
		Done = true
	}
	return
}

func infoCommand(args []string) (err error) {
	if len(args) == 0 {
		Done = true
		return errors.New("package name is required")
	}
	err = client.InfoRequest(args[FirstTargetIndex])
	Done = true
	return
}

func installCommand(args []string) (err error) {
	if len(args) == 0 {
		Done = true
		return errors.New("package name is required")
	}
	err = client.InstallRequest(args[FirstTargetIndex])
	Done = true
	return
}

func uploadCommand(args []string) (err error) {
	switch len(args) {
	case 0:
		return errors.New("local package path is required")
	case 1:
		err = client.UploadRequest(args[FirstTargetIndex], "")
		Done = true
	case 2:
		err = client.UploadRequest(args[FirstTargetIndex], args[SecondTargetIndex])
		Done = true

	}

	return
}

func mvCommand(args []string) (err error) {
	if len(args) != 2 {
		return errors.New("args num error")
	}

	err = client.MvRequest(args[FirstTargetIndex], args[SecondTargetIndex])
	Done = true

	return
}

func rmCommand(args []string) (err error) {
	if len(args) != 1 {
		return errors.New("args num error")
	}

	err = client.RmRequest(args[FirstTargetIndex])
	Done = true

	return
}

func syncCommand() (err error) {
	err = client.SyncRequest()
	Done = true
	return
}

func tagCommand(args []string) (err error) {
	switch len(args) {
	case 1:
		switch args[FirstTargetIndex] {
		case "ls", "list":
			err = client.TagListRequest()
			Done = true
		}
	case 2:
		switch args[FirstTargetIndex] {
		case "add":
			err = client.AddTagRequest(args[SecondTargetIndex])
			Done = true
		case "rm":
			err = client.TagRmRequest(args[SecondTargetIndex])
			Done = true
		default:
			err = client.UpdateTagRequest(args[FirstTargetIndex], args[SecondTargetIndex])
			Done = true
		}
	}

	return
}

func reloadCommand() (err error) {
	err = client.ReloadRequest()
	Done = true
	return
}

func publicComand(args []string) (err error) {
	err = client.PublicRequest(args[FirstTargetIndex])
	Done = true
	return
}

func linksCommand() (err error) {
	err = client.LinksRequest()
	Done = true
	return
}

func privateCommand(args []string) (err error) {
	err = client.PrivateRequest(args[FirstTargetIndex])
	Done = true
	return
}

func listServersCommand() {
	config.ListAvailableServers()
	Done = true
}

func changeServerCommand(args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("invalid args : alias or host required")
	}
	if strings.Contains(args[0], "http") {
		host := args[0]

		config.ChangeCurrServer("", host)

		Done = true
		return
	}

	alias := args[0]

	config.ChangeCurrServer(alias, "")

	Done = true

	return
}

func addserverCommand(args []string) (err error) {
	if len(args) < 2 {
		Done = true
		return fmt.Errorf("invalid args : you must provide <alias> and <host>")
	}

	config.AddAvialableServer(args[0], args[1])
	Done = true

	return
}

func delServerCommand(args []string) (err error) {
	if len(args) == 0 {
		Done = true
		return fmt.Errorf("invalid args : alias or host required")
	}

	if strings.Contains(args[0], "http") {
		host := args[0]

		config.DelAvialableServer("", host)

		Done = true

		return
	}

	alias := args[0]

	config.DelAvialableServer(alias, "")

	Done = true

	return
}
