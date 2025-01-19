package middlewares

import (
	"giter/initializer"

	"github.com/gin-gonic/gin"
)

// LogMiddlewareはリクエストの開始と終了をログに記録するミドルウェアです。
func LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		initializer.SetRequestID(ctx)               // リクエストIDを設定
		logger := initializer.MiddlewareLogger(ctx) // ロガーを初期化
		logger.Info().Msg("REQUEST START")          // リクエスト開始をログに記録
		ctx.Next()                                  // 次のハンドラーを呼び出し
		logger.Info().Msg("REQUEST END")            // リクエスト終了をログに記録
	}
}
