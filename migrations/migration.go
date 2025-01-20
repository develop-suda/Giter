package main

import (
	"giter/infra"
	"giter/models"
)

// pokemonsテーブルを削除し再作成後、データ投入
func main() {
	infra.Initialize()    // インフラの初期化
	db := infra.SetupDB() // データベースのセットアップ

	db.Migrator().DropTable(&models.User{})
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("Failed to migrate database") // エラーが発生した場合、パニックを起こす
	}
}
