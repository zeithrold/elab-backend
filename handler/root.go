package handler

import (
	"elab-backend/handler/apply"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Init() *gin.Engine {
	slog.Info("正在初始化路由")
	r := gin.Default()
	apply.ApplyRoute(r)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	return r
}
