package auth

import (
	"elab-backend/util/auth"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// EnsureValidToken 用于检查用户是否已经登录。
// 需要注意，在客户端中，需要指定audience。
//
//	 const { authorize } = useAuth0();
//		await authorize({
//		  audience: process.env.AUTH0_AUDIENCE,
//		})
func EnsureValidToken() gin.HandlerFunc {
	slog.Debug("middleware.auth.EnsureValidToken: 进入Token验证中间件")
	jwtValidator := auth.GetValidator()
	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("jwt验证失败", "error", err)
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusUnauthorized)
		//_, err = w.Write([]byte(`{"message":"用户验证失败。", "description":"您的Token无效。"}`))
		//if err != nil {
		//	slog.Error("无法写入错误信息", "error", err)
		//	panic(err)
		//}
	}
	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)
	return func(c *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			c.Request = r
			c.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)
		if encounteredError {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message":     "用户验证失败。",
				"description": "您的Token无效。",
			})
		}
	}
}
