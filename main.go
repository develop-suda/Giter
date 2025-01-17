package main

import (
	"fmt"
	"giter/controllers"
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
	requestController := di.InitializeRouter(clients.RClient, clients.GClient)
	authController := controllers.NewAuthController(db)

	r := gin.Default()

	// 各リクエストにリクエストIDを生成
	r.Use(requestid.New())
	// カスタムロギングミドルウェアを使用
	r.Use(middlewares.LogMiddleware())
	// HTMLテンプレートを読み込む
	r.LoadHTMLGlob("templates/*")

	r.GET("/", requestController.Index)
	r.GET("/commit", requestController.GetCommits)
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	// 未定義のルートをホームページにリダイレクト
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	fmt.Println("server start")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
