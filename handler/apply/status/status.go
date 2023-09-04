package status

import (
	"elab-backend/model/apply"
	"elab-backend/service/redis"
	"elab-backend/util/auth"
	"github.com/gin-gonic/gin"
)

func ApplyRoute(group *gin.RouterGroup) {
	route := group.Group("/status")
	route.GET("", GetStatus)
}

func GetStatus(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	unlock, err := redis.GetLock(ctx, "textform:"+openid)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"message": "服务器错误",
		})
	}
	defer unlock()
	ctx.JSON(200, apply.GetStatus(ctx, openid))
}
