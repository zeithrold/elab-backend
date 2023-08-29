package model

import (
	"elab-backend/model/apply"
	"elab-backend/service"
	"log/slog"
)

func Init() {
	slog.Debug("model.Init: 正在初始化数据库")
	svc := service.GetService()
	slog.Debug("model.Init: 正在迁移数据库")
	err := svc.DB.AutoMigrate(
		&apply.Config{}, &apply.Room{}, &apply.TextForm{}, &apply.Ticket{}, &apply.Selection{}, &apply.Question{})
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	slog.Debug("model.Init: 数据库初始化完成")
}
