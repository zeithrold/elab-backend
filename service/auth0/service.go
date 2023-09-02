package auth0

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"log/slog"
	"os"
)

func NewService() *management.Management {
	slog.Debug("service.auth0.NewService: 正在创建authAPI")
	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_MANAGEMENT_API_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_MANAGEMENT_API_CLIENT_SECRET")
	api, err := management.New(
		domain,
		management.WithClientCredentials(context.Background(), clientID, clientSecret),
	)
	if err != nil {
		slog.Error("无法创建authAPI", "error", err)
	}
	slog.Debug("service.auth0.NewService: authAPI创建成功", "access_token", api)
	return api
}
