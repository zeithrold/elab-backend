package apply

import (
	"elab-backend/model/apply"
	"github.com/gin-gonic/gin"
)

func GetConfig(ctx *gin.Context) {
	ctx.JSON(200, apply.GetConfig(ctx))
}
