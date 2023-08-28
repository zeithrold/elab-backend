package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Init() *gin.Engine {
	slog.Info("正在初始化路由")
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	return gin.Default()
}
