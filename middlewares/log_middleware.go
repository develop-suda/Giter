package middlewares

import (
	"giter/initializer"

	"github.com/gin-gonic/gin"
)

// LogMiddlewareはリクエストの開始と終了をログに記録するミドルウェアです。
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		initializer.SetRequestID(c)               // リクエストIDを設定
		logger := initializer.MiddlewareLogger(c) // ロガーを初期化
		logger.Info().Msg("REQUEST START")        // リクエスト開始をログに記録
		c.Next()                                  // 次のハンドラーを呼び出し
		logger.Info().Msg("REQUEST END")          // リクエスト終了をログに記録
	}
}
