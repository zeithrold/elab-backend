package main

import (
	"elab-backend/handler"
	"elab-backend/model"
	"elab-backend/service"
	"elab-backend/util/config"
	"fmt"
	"log/slog"
)

func main() {
	slog.Info("正在启动Web服务器")
	config.Load()
	service.Init()
	model.Init()
	r := handler.Init()
	err := r.Run(":2333")
	if err != nil {
		fmt.Println(err)
	}
}
