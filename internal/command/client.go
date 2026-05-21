package command

import (
	"errors"
	"fmt"
	"log/slog"

	"gitee.com/cai-zixiang_hainan/wt/internal/client"
	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/pkg/loghelper"
)

const (
	CommandIndex      = 0
	FirstTargetIndex  = 0
	SecondTargetIndex = 1
)

var Done bool

func Usage(command string) {
	fmt.Println("\n=== Usage Guide ===")

	switch command {
	case "search":
		fmt.Println("  wt search <package-name>")
		fmt.Println("  Example: wt search my-app")

	case "info":
		fmt.Println("  wt info <package-name>")
		fmt.Println("  Example: wt info my-app-v1")

	case "install":
		fmt.Println("  wt install <package-name>")
		fmt.Println("  Example: wt install my-app-v1")

	case "upload":
		fmt.Println("  wt upload <file-path> [package-name]")
		fmt.Println("  Example: wt upload ./build/app.tar.gz")
		fmt.Println("  Example: wt upload ./build/app.tar.gz my-app-v1")

	case "mv":
		fmt.Println("  wt mv <old-name> <new-name>")
		fmt.Println("  Example: wt mv my-app-v1 my-app-v2")

	case "rm":
		fmt.Println("  wt rm <package-name>")
		fmt.Println("  Example: wt rm my-app-v1")

	case "list", "ls":
		fmt.Println("  wt list [tag]")
		fmt.Println("  Example: wt list")
		fmt.Println("  Example: wt list latest")

	case "sync":
		fmt.Println("  wt sync")
		fmt.Println("  Description: Sync local metadata with server")

	case "help", "--help", "-h":
		fmt.Println("  wt help")
		fmt.Println("  Description: Show this help manual")

	default:
		fmt.Println("  Unknown command:", command)
		fmt.Printf("\nAvailable commands:\n")
		fmt.Println("  search   - Search for packages")
		fmt.Println("  info     - Show package information")
		fmt.Println("  install  - Download and install a package")
		fmt.Println("  upload   - Upload a package to server")
		fmt.Println("  mv       - Rename a package")
		fmt.Println("  rm       - Remove a package")
		fmt.Println("  list     - List packages by tag")
		fmt.Println("  sync     - Sync metadata with server")
		fmt.Println("  help     - Show help information")
	}

	fmt.Println("\nFor more information, use: wt help")
}

func handleZeroTarget(command string) (err error) {
	switch command {
	case "sync":
		err = client.SyncRequest()
		Done = true
	case "ls", "list":
		err = client.ListRequest(config.DefaultTagTemp)
		Done = true
	case "help", "--help", "-h":
		fmt.Println(DefaultManual)
		Done = true
	case "tags":
		err = client.TagListRequest()
		Done = true
	}

	return
}

func handleOneTarget(command string, args []string) (err error) {
	firstTarget := args[FirstTargetIndex]

	switch command {
	case "search":
		err = client.SearchRequest(firstTarget)
		Done = true
	case "info":
		err = client.InfoRequest(firstTarget)
		Done = true
	case "install":
		err = client.InstallRequest(firstTarget)
		Done = true
	case "rm":
		err = client.RmRequest(firstTarget)
		Done = true
	case "list", "ls":
		err = client.ListRequest(firstTarget)
		Done = true
	case "upload":
		err = client.UploadRequest(firstTarget, "")
		Done = true
	case "config", "-c":
		if firstTarget == "show" {
			config.ClientConfigShow()
			Done = true
		}
	case "help", "--help", "-h":
		if firstTarget == "simple" {
			fmt.Println(SimpleManual)
		} else if firstTarget == "advance" {
			fmt.Println(AdvanceManual)
		} else {
			fmt.Print(DefaultManual)
		}
		Done = true
	}

	return
}

func handleTwoTarget(command string, args []string) (err error) {
	firstTarget := args[FirstTargetIndex]
	secondTarget := args[SecondTargetIndex]

	switch command {
	case "upload":
		err = client.UploadRequest(firstTarget, secondTarget)
		Done = true
	case "mv":
		err = client.MvRequest(firstTarget, secondTarget)
		Done = true
	case "config":
		err = config.AlterClientConfig(firstTarget, secondTarget)
		Done = true
	}
	return
}

func ClientMain(args []string, debug bool) {
	loghelper.InitClientLogger(debug)

	err := config.InitClientConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	command := args[CommandIndex]
	slog.Debug("client exec command", "command", command)
	args = args[1:]

	switch len(args) {
	case 0:
		err = handleZeroTarget(command)
	case 1:
		err = handleOneTarget(command, args)
	case 2:
		err = handleTwoTarget(command, args)
	default:
		fmt.Printf("too many arguments for command %s\n", command)
		Usage(command)
		return
	}

	if err != nil {
		fmt.Println("fail to exec ", command, " err :", err.Error())
		Usage(command)
	}
	if !Done {
		fmt.Println("Bad Usage")
		Usage(command)
	}
}

func NewClientMain(args []string, debug bool) {
	loghelper.InitClientLogger(debug)

	err := config.InitClientConfig()
	if err != nil {
		fmt.Println("err when init the client config: ", err)
		return
	}

	command := args[CommandIndex]
	args = args[CommandIndex+1:]

	fmt.Println("command:", command)
	fmt.Println("enter args:", args)

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
	}

	if err != nil {
		Usage(command)
	}

	if err != nil {
		fmt.Println("fail to exec ", command, " err :", err.Error())
		Usage(command)
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

	err = client.SearchRequest(args[0])
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
		err = client.UploadRequest(args[FirstTargetIndex], args[SecondTarget])
		Done = true

	}

	return
}

func mvCommand(args []string) (err error) {
	if len(args) != 2 {
		return errors.New("args num error")
	}

	err = client.MvRequest(args[FirstTargetIndex], args[SecondTarget])
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
	fmt.Println("enter the tag command branch")
	switch len(args) {
	case 1:
		err = client.TagListRequest()
		Done = true
	case 2:
		switch args[0] {
		case "add":
			fmt.Println("enter the tag add branch")
			err = client.AddTagRequest(args[SecondTargetIndex])
			Done = true
		}
	}

	return
}
