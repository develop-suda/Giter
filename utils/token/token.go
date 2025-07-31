package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

// 指定されたユーザーIDに基づいてJWTトークンを生成する
func GenerateToken(id uint) (string, error) {
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func extractTokenString(ctx *gin.Context) string {
	// First try to get token from Authorization header
	bearToken := ctx.Request.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 && strings.ToLower(strArr[0]) == "bearer" {
		return strArr[1]
	}

	// If not found in header, try to get from cookie
	token, err := ctx.Cookie("jwt")
	if err == nil && token != "" {
		return token
	}

	return ""
}

func parseToken(tokenString string) (*jwt.Token, error) {
	// トークンの検証、トークンのパースを実行
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// トークンが有効かどうかを検証
func TokenValid(ctx *gin.Context) error {
	tokenString := extractTokenString(ctx)
	if tokenString == "" {
		return fmt.Errorf("no token provided")
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// トークンからユーザーIDを取得
func ExtractTokenId(ctx *gin.Context) (uint, error) {
	tokenString := extractTokenString(ctx)
	if tokenString == "" {
		return 0, fmt.Errorf("no token provided")
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token claims")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}

	return uint(userId), nil
}
