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
	VERSION = "0.0.1"
	CMDINX  = 1
)

func main() {
	args := os.Args

	switch args[CMDINX] {
	case "uninstall":
		err := store.Uninstall()
		if err != nil {
			fmt.Println("error when uninstall: ", err)
		}
	case "server", "serve":
		command.ServerMain(args[CMDINX:], DEBUG)
	default:
		command.ClientMain(args[CMDINX:], DEBUG)
	}
}
