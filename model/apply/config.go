package apply

import (
	"context"
	"elab-backend/service"
	"gorm.io/gorm"
	"log/slog"
)

type Config struct {
	gorm.Model
	// Key 是配置键。
	Key string `gorm:"type:varchar(1024)"`
	// Value 是配置值。
	Value string `gorm:"type:varchar(1024)"`
}

// GetConfig 获取配置。
//
// ctx 是上下文。
func GetConfig(ctx context.Context) map[string]string {
	slog.Debug("model.GetConfig: 正在获取配置")
	svc := service.GetService()
	var config []Config
	err := svc.DB.WithContext(ctx).Model(&Config{}).Find(&config).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	result := make(map[string]string)
	for _, c := range config {
		result[c.Key] = c.Value
	}
	return result
}
