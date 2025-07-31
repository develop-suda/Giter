package infra

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Initialize 関数は .env ファイルを読み込み、必要な環境変数を検証します。
// エラーが発生した場合、ログにエラーメッセージを出力してプログラムを終了します。
func Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 必要な環境変数を検証
	validateRequiredEnvVars()
}

// validateRequiredEnvVars は必要な環境変数が設定されているかを検証します
func validateRequiredEnvVars() {
	requiredVars := []string{
		"API_SECRET",
		"TOKEN_HOUR_LIFESPAN",
		"GITHUB_KEY",
		"GITHUB_SECRET",
		"GITHUB_CALLBACK_URL",
	}

	for _, envVar := range requiredVars {
		if value := os.Getenv(envVar); value == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}

	// TOKEN_HOUR_LIFESPAN が有効な数値かチェック
	if tokenLifespan := os.Getenv("TOKEN_HOUR_LIFESPAN"); tokenLifespan != "" {
		if _, err := strconv.Atoi(tokenLifespan); err != nil {
			log.Fatalf("TOKEN_HOUR_LIFESPAN must be a valid integer, got: %s", tokenLifespan)
		}
	}

	// API_SECRET の長さをチェック（セキュリティ）
	if apiSecret := os.Getenv("API_SECRET"); len(apiSecret) < 32 {
		log.Fatal("API_SECRET must be at least 32 characters long for security")
	}
}
