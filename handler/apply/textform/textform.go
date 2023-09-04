package textform

import (
	"elab-backend/model/apply"
	"elab-backend/service/redis"
	"elab-backend/util/auth"
	"github.com/gin-gonic/gin"
)

func ApplyRoute(group *gin.RouterGroup) {
	textFormRoute := group.Group("/textform")
	textFormRoute.Use(LockMiddleware())
	textFormRoute.GET("", GetTextForm)
	textFormRoute.PATCH("/:id", UpdateTextForm)
	questionRoute := group.Group("/question")
	questionRoute.Use(LockMiddleware())
	questionRoute.GET("", GetQuestionList)
	questionRoute.GET("/:id", GetQuestion)
}

func LockMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := auth.GetToken(ctx)
		openid := token.RegisteredClaims.Subject
		unlock, err := redis.GetLock(ctx, "textform:"+openid)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"message": "服务器错误",
			})
		}
		defer unlock()
		ctx.Next()
	}
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

func GetQuestion(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	var request apply.GetQuestionRequestUri
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": "请求格式错误",
		})
		return
	}
	ctx.JSON(200, apply.GetQuestion(ctx, openid, request.Id))
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
