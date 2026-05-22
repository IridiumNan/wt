package main

import (
	_ "embed"
	"fmt"
	"os"

	"gitee.com/cai-zixiang_hainan/wt/internal/command"
	"gitee.com/cai-zixiang_hainan/wt/internal/store"
)

const (
	DEBUG   = false
	VERSION = "0.2.1"
	CMDINX  = 1
)

func helpCommand(args []string) {
	if len(args) == 0 {
		fmt.Println(command.DefaultManual)
		return
	}
	switch args[0] {
	case "simple":
		fmt.Println(command.SimpleManual)
	case "advance":
		fmt.Println(command.AdvanceManual)
	default:
		command.Usage(args[0])
	}
}

func main() {
	args := os.Args

	switch args[CMDINX] {
	case "uninstall":
		err := store.Uninstall()
		if err != nil {
			fmt.Println("error when uninstall: ", err)
		}
	case "help", "-h", "--help":
		helpCommand(args[CMDINX+1:])
	case "server", "serve":
		command.ServerMain(args[CMDINX:], DEBUG)
	default:
		// command.NewClientMain(args[CMDINX:], DEBUG)
		command.ClientMain(args[CMDINX:], DEBUG)
	}
}
