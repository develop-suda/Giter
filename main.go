package main

import (
	"fmt"
	"giter/di"
	"giter/infra"
	"giter/initializer"
	"giter/middlewares"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func main() {

	// アプリケーション設定の初期化
	initializer.Init()
	initializer.Log()

	// インフラ設定の初期化
	infra.Initialize()
	db := infra.SetupDB()

	clients := initializer.NewClients()
	requestController := di.InitCommitRouter(clients.RClient, clients.GClient)
	authController := di.InitAuthRouter(db)

	r := gin.Default()

	// 各リクエストにリクエストIDを生成
	r.Use(requestid.New())
	// カスタムロギングミドルウェアを使用
	r.Use(middlewares.LogMiddleware())
	// HTMLテンプレートを読み込む
	r.LoadHTMLGlob("templates/*")

	r.GET("/", requestController.Index)
	r.GET("/commit", requestController.GetCommits)
	r.POST("/register", authController.RegisterUser)
	r.POST("/login", authController.Login)

	protected := r.Group("/admin")
	// JWT認証ミドルウェアを適用
	protected.Use(middlewares.JwtAuthMiddleware())
	// 認証されたユーザー情報を取得するルートを定義
	protected.GET("/user", authController.CurrentUser)

	// 未定義のルートをホームページにリダイレクト
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	fmt.Println("server start")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
