package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

var db *gorm.DB

func NewService() *gorm.DB {
	slog.Debug("db.NewService: 正在初始化数据库")
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)
	localDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("无法连接数据库", "error", err)
		panic(err)
	}
	db = localDb
	return localDb
}

func GetDb() *gorm.DB {
	if db == nil {
		err := fmt.Errorf("db未初始化")
		slog.Error("db未初始化", "error", err)
		panic(err)
	}
	return db
}
