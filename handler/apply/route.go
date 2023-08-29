package apply

import (
	"elab-backend/handler/apply/room"
	"elab-backend/handler/apply/status"
	"elab-backend/handler/apply/textform"
	"elab-backend/handler/apply/ticket"
	"elab-backend/middleware/auth"
	"github.com/gin-gonic/gin"
)

func NewHandler(r *gin.RouterGroup) {
	group := r.Group("/apply")
	group.Use(auth.EnsureValidToken())
	group.GET("/config", GetConfig)
	room.ApplyRoute(group)
	status.ApplyRoute(group)
	textform.ApplyRoute(group)
	ticket.ApplyRoute(group)
}
