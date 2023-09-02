package handler

import (
	"elab-backend/handler/apply"
	"elab-backend/handler/auth"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Init() *gin.Engine {
	slog.Info("handler.Init: 正在初始化路由")
	r := gin.Default()
	endpoint := r.Group("/v1")
	apply.NewHandler(endpoint)
	auth.NewHandler(endpoint)
	endpoint.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	return r
}
