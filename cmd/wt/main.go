package main

import (
	_ "embed"
	"fmt"
	"os"

	"gitee.com/cai-zixiang_hainan/wt/internal/client"
	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/pkg/loghelper"
)

const (
	DEBUG   = false
	VERSION = "0.0.1"
)

func init() {
	loghelper.JSONLog(DEBUG)
}

func main() {
	args := os.Args

	err := config.InitClientConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	client.ClientMain(args)
}
