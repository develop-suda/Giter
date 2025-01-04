package initializer

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func MiddlewareLogger(c *gin.Context) zerolog.Logger {
	// configured-logger
	logger := log.With().
		Str("request_id", RequestID).
		Str("ip", c.ClientIP()).
		Str("path", c.Request.URL.Path).
		Str("method", c.Request.Method).
		Str("params", c.Request.URL.RawQuery).
		Logger()

	return logger
}

func DefaultLogger() zerolog.Logger {
	return log.With().Caller().Logger()
}
