package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"strconv"
)

var client *redis.Client

// NewService 用于初始化Redis服务。
func NewService() *redis.Client {
	slog.Info("service.redis.NewService: 正在初始化Redis服务")
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		slog.Error("无法解析REDIS_DB", "error", err)
		panic(err)
	}
	localClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	TestConnection(ctx, localClient)
	client = localClient
	return localClient
}

// TestConnection 获取Redis客户端。
//
// ctx 是上下文。
// client 是Redis客户端。
func TestConnection(ctx context.Context, client *redis.Client) {
	slog.Info("service.redis.TestConnection: 正在检查Redis连接")
	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		slog.Error("无法连接Redis", "error", err)
		panic(err)
	}
	value, err := client.Get(ctx, "foo").Result()
	if err != nil {
		slog.Error("无法连接Redis", "error", err)
		panic(err)
	}
	if value != "bar" {
		slog.Error("Redis无法设置值，value应该是'bar'", "error", err, "value", value)
		panic(err)
	}
	client.Del(ctx, "foo")
	slog.Info("service.redis.TestConnection: 完成Redis连接检查")
}

// GetClient 获取Redis客户端。
func GetClient() *redis.Client {
	if client == nil {
		err := "client未初始化"
		slog.Error("service.redis.GetClient: client未初始化", "error", err)
		panic(err)
	}
	return client
}
