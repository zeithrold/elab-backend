package auth

import "github.com/gin-gonic/gin"

func NewHandler(r *gin.RouterGroup) {
	route := r.Group("/auth")
	route.DELETE("", DeleteAccount)
}
