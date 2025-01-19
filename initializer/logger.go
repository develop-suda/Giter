package initializer

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func MiddlewareLogger(ctx *gin.Context) zerolog.Logger {
	// configured-logger
	logger := log.With().
		Str("request_id", RequestID).
		Str("ip", ctx.ClientIP()).
		Str("path", ctx.Request.URL.Path).
		Str("method", ctx.Request.Method).
		Str("params", ctx.Request.URL.RawQuery).
		Logger()

	return logger
}

func DefaultLogger() zerolog.Logger {
	return log.With().Caller().Logger()
}
