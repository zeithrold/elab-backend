package ticket

import (
	"elab-backend/model/apply"
	"elab-backend/util/auth"
	"github.com/gin-gonic/gin"
)

func ApplyRoute(group *gin.RouterGroup) {
	route := group.Group("/ticket")
	route.GET("", GetTicket)
	route.PATCH("", UpdateTicket)
}

func GetTicket(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	ctx.JSON(200, apply.GetTicket(ctx, openid))
}

func UpdateTicket(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	var request apply.TicketBody
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": "请求格式错误",
		})
		return
	}
	apply.UpdateTicket(ctx, openid, &request)
	ctx.JSON(200, gin.H{
		"message": "更新成功",
	})
}
