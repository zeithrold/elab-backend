package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func Load() {
	slog.Info("正在加载环境变量")
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		slog.Info("GIN_MODE为空，设置为debug")
		mode = "debug"
		err := os.Setenv("GIN_MODE", "debug")
		if err != nil {
			slog.Error("无法设置GIN_MODE", "error", err)
			panic(err)
		}
	}
	err := loadDotEnvFile("release")
	if err != nil {
		slog.Error("无法加载环境变量", "error", err)
		panic(err)
	}
	// 然后判断是否为release模式
	if mode != "release" {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
		slog.Info("GIN_MODE是调试模式。")
		err = loadDotEnvFile("development")
		if err != nil {
			slog.Error("无法加载环境变量", "error", err)
			panic(err)
		}
	}
}

func loadDotEnvFile(fileType string) error {
	fileName := fmt.Sprintf(".env.%s", fileType)
	if fileType == "release" {
		fileName = ".env"
	}
	_, err := os.Stat(fileName)
	if err != nil {
		slog.Info("无法找到环境变量文件", "fileName", fileName)
		return nil
	} else {
		err = godotenv.Load(fileName)
	}
	return err
}
