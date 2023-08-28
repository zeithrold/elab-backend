package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func Load() {
	slog.Info("正在加载环境变量")
	err := loadDotEnvFile("prod")
	if err != nil {
		slog.Error("无法加载环境变量", "error", err)
		panic(err)
	}
	// 然后判断是否为release模式
	if os.Getenv("GIN_MODE") != "release" {
		slog.Debug("GIN_MODE是调试模式。")
		err = loadDotEnvFile("development")
		if err != nil {
			slog.Error("无法加载环境变量", "error", err)
			panic(err)
		}
	}
}

func loadDotEnvFile(fileType string) error {
	fileName := fmt.Sprintf(".env.%s", fileType)
	if fileType == "prod" {
		fileName = ".env"
	}
	_, err := os.Stat(fileName)
	if err != nil {
		slog.Info("无法找到环境变量文件", "fileName", fileName)
		return nil
	} else {
		err = godotenv.Load(".env")
	}
	return err
}
