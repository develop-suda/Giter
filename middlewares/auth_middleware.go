package middlewares

import (
	"giter/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware はJWT認証を行うミドルウェアを返します
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := token.TokenValid(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
