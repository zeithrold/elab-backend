package status

import (
	"elab-backend/model/apply"
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
	ctx.JSON(200, apply.GetStatus(ctx, openid))
}
