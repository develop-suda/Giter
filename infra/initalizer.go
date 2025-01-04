package infra

import (
	"log"

	"github.com/joho/godotenv"
)

// Initialize 関数は .env ファイルを読み込みます。
// エラーが発生した場合、ログにエラーメッセージを出力してプログラムを終了します。
func Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
