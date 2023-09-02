package auth

import (
	"elab-backend/util/auth"
	"github.com/gin-gonic/gin"
)

func DeleteAccount(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	err := auth.DeleteAccount(ctx, openid)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "删除成功",
	})
}
