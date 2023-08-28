package redis

import (
	"context"
	"log/slog"
	"time"
)

const Timeout = time.Second * 10
const RetryTimes = 50
const RetryInterval = time.Millisecond * 200

type GetLockTimeoutError struct{}

func (e *GetLockTimeoutError) Error() string {
	return "获取锁超时"
}

// GetLock 获取锁。
//
// ctx 是上下文。
// key 是锁的键。
func GetLock(ctx context.Context, key string) (func(), error) {
	slog.Debug("redis.GetLock: 正在获取锁", "key", key)
	for i := 0; i < RetryTimes; i++ {
		ok, err := client.SetNX(ctx, key, 1, Timeout).Result()
		if err != nil {
			slog.Error("无法获取锁", "error", err)
			return nil, err
		}
		if ok {
			slog.Debug("redis.GetLock: 成功获取锁", "key", key)
			return func() {
				client.Del(ctx, key)
			}, nil
		}
		slog.Debug("redis.GetLock: 未能获取锁，正在重试", "key", key, "retry", i)
		time.Sleep(RetryInterval)
	}
	slog.Debug("redis.GetLock: 未能获取锁，超时", "key", key)
	return nil, &GetLockTimeoutError{}
}
