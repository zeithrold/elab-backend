package room

import (
	"elab-backend/model/apply"
	"elab-backend/util/auth"
	"github.com/gin-gonic/gin"
)

func ApplyRoute(group *gin.RouterGroup) {
	route := group.Group("/room")
	route.GET("", GetRoomList)
	route.GET("/date", GetRoomDateList)
	route.POST("/selection", SetSelection)
	route.GET("/selection", GetSelection)
}

func GetRoomList(ctx *gin.Context) {
	date := ctx.Query("date")
	ctx.JSON(200, apply.GetRoomList(ctx, date))
}

func GetRoomDateList(ctx *gin.Context) {
	ctx.JSON(200, apply.GetRoomDateList(ctx))
}

func SetSelection(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	var request apply.SetRoomSelectionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{
			"message": "请求格式错误",
		})
		return
	}
	err := apply.SetSelection(ctx, openid, request.RoomId)
	if err != nil {
		switch v := err.(type) {
		case *apply.RoomFullError:
			ctx.JSON(400, gin.H{
				"message": v.Error(),
			})
			return
		case *apply.DuplicateSelectionError:
			ctx.JSON(400, gin.H{
				"message": v.Error(),
			})
			return
		}
	}
	ctx.JSON(200, gin.H{
		"message": "更新成功",
	})
}

func GetSelection(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	selection, err := apply.GetSelection(ctx, openid)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
	}
	ctx.JSON(200, gin.H{
		"room_id": selection.RoomId,
	})
}