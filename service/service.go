package service

import (
	"elab-backend/service/db"
	"elab-backend/service/redis"
	"github.com/pkg/errors"
	libRedis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
)

type Service struct {
	DB    *gorm.DB
	Redis *libRedis.Client
}

var service *Service

func Init() {
	slog.Info("正在初始化服务")
	service = &Service{}
	service.Redis = redis.NewService()
	service.DB = db.NewService()
}

func GetService() *Service {
	if service == nil {
		err := errors.New("service未初始化")
		slog.Error("service未初始化", "error", err)
		panic(err)
	}
	return service
}
