package textform

import (
	"elab-backend/model/apply"
	"elab-backend/util/auth"
	"github.com/gin-gonic/gin"
)

func ApplyRoute(group *gin.RouterGroup) {
	route := group.Group("/textform")
	route.GET("", GetTextForm)
	route.GET("/question", GetQuestionList)
	route.PATCH("/:id", UpdateTextForm)
}

func GetQuestionList(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	ctx.JSON(200, apply.GetQuestionList(ctx, openid))
}

func GetTextForm(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	ctx.JSON(200, apply.GetTextForm(ctx, openid))
}

func UpdateTextForm(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	var request apply.UpdateTextFormRequest
	var requestUri apply.UpdateTextFormRequestUri
	if err := ctx.ShouldBindUri(&requestUri); err != nil {
		ctx.JSON(400, gin.H{
			"message": "请求格式错误",
		})
		return
	}
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": "请求格式错误",
		})
		return
	}
	request.Id = requestUri.Id
	apply.UpdateTextForm(ctx, openid, &request)
	ctx.JSON(200, gin.H{
		"message": "更新成功",
	})
}
