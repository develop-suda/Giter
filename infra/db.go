package infra

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDB() *gorm.DB {
	// DSNをフォーマット
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("HOST_NAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PORT"),
	)

	var (
		db  *gorm.DB
		err error
	)

	// データベースに接続
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// エラーチェック
	if err != nil {
		panic("Failed to connect database")
	}

	return db
}
