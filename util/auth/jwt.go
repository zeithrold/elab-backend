package auth

import (
	"context"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/url"
	"os"
	"time"
)

var v *validator.Validator

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func GetValidator() *validator.Validator {
	slog.Debug("util.auth.GetValidator: 正在获取验证器")
	if v != nil {
		slog.Debug("util.auth.GetValidator: 验证器已存在")
		return v
	}
	slog.Debug("util.auth.GetValidator: 验证器不存在，正在创建")
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		slog.Error("无法解析issuerURL", "error", err)
		panic(err)
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			}),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		slog.Error("无法创建jwtValidator", "error", err)
		panic(err)
	}
	v = jwtValidator
	return jwtValidator
}

func GetToken(ctx *gin.Context) *validator.ValidatedClaims {
	slog.Debug("util.auth.GetToken: 正在获取Token")
	authHeader := ctx.GetHeader("Authorization")
	vail := GetValidator()
	token, err := vail.ValidateToken(ctx, authHeader[7:])
	if err != nil {
		slog.Error("无法获取Token", "error", err)
		panic(err)
	}
	return token.(*validator.ValidatedClaims)
}
