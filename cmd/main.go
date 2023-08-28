package main

import (
	"elab-backend/handler"
	"elab-backend/service"
	"elab-backend/util/config"
	"fmt"
	"log/slog"
)

func main() {
	slog.Info("正在启动Web服务器")
	config.Load()
	service.Init()
	r := handler.Init()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
