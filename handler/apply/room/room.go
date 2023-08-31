package room

import (
	"elab-backend/model/apply"
	"elab-backend/service/redis"
	"elab-backend/util/auth"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func ApplyRoute(group *gin.RouterGroup) {
	route := group.Group("/room")
	route.GET("", GetRoomList)
	route.GET("/date", GetRoomDateList)
	route.POST("/selection", SetSelection)
	route.DELETE("/selection", ClearSelection)
	route.GET("/selection", GetSelection)
}

func GetRoomList(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(400, gin.H{
			"message": "请求格式错误，缺少参数 date",
		})
		return
	}
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
	unlock, err := redis.GetLock(ctx, "room_selection")
	if err != nil {
		slog.Error("handler.apply.room.SetSelection: 获取锁失败", "err", err)
		ctx.JSON(400, gin.H{
			"message": "请求失败",
		})
		return
	}
	defer unlock()
	err = apply.SetSelection(ctx, openid, request.Id)
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

func ClearSelection(ctx *gin.Context) {
	token := auth.GetToken(ctx)
	openid := token.RegisteredClaims.Subject
	isSelectionExists := apply.CheckIsSelectionExists(ctx, openid)
	if !isSelectionExists {
		ctx.JSON(400, gin.H{
			"message": "请求失败，用户未选择",
		})
		return
	}
	unlock, err := redis.GetLock(ctx, "room_selection")
	if err != nil {
		slog.Error("handler.apply.room.SetSelection: 获取锁失败", "err", err)
		ctx.JSON(400, gin.H{
			"message": "请求失败",
		})
		return
	}
	defer unlock()
	err = apply.ClearSelection(ctx, openid)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "清除成功",
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
		return
	}
	ctx.JSON(200, gin.H{
		"id": selection.RoomId,
	})
}
