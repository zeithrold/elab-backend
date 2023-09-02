package service

import (
	"elab-backend/service/auth0"
	"elab-backend/service/db"
	"elab-backend/service/redis"
	"github.com/auth0/go-auth0/management"
	"github.com/pkg/errors"
	libRedis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
)

type Service struct {
	DB      *gorm.DB
	Redis   *libRedis.Client
	AuthAPI *management.Management
}

var service *Service

func Init() {
	slog.Info("正在初始化服务")
	service = &Service{}
	service.Redis = redis.NewService()
	service.DB = db.NewService()
	service.AuthAPI = auth0.NewService()
}

func GetService() *Service {
	if service == nil {
		err := errors.New("service未初始化")
		slog.Error("service未初始化", "error", err)
		panic(err)
	}
	return service
}
