package main

import (
	"fmt"
	"giter/di"
	"giter/infra"
	"giter/initializer"
	"giter/middlewares"
	"net/http"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

func main() {

	// アプリケーション設定の初期化
	initializer.Init()
	initializer.Log()

	// インフラ設定の初期化
	infra.Initialize()
	db := infra.SetupDB()

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/auth/github/callback", "user"),
	)

	clients := initializer.NewClients()
	requestController := di.InitCommitRouter(clients.RClient, clients.GClient)
	authController := di.InitAuthRouter(db)
	githubAuthController := di.InitGithubAuthRouter(db)

	r := gin.Default()

	// 各リクエストにリクエストIDを生成
	r.Use(requestid.New())
	// カスタムロギングミドルウェアを使用
	r.Use(middlewares.LogMiddleware())
	// HTMLテンプレートを読み込む
	r.LoadHTMLGlob("templates/*")

	r.GET("/", requestController.Index)
	r.GET("/commit", requestController.GetCommits)
	r.GET("/signup", authController.SignupView)
	r.POST("/signup", authController.Register)
	r.GET("/login", authController.LoginView)
	r.POST("/login", authController.Login)

	protected := r.Group("/admin")
	// JWT認証ミドルウェアを適用
	protected.Use(middlewares.JwtAuthMiddleware())
	// 認証されたユーザー情報を取得するルートを定義
	protected.GET("/user", authController.CurrentUser)

	// GitHubログイン
	r.GET("/auth/github", githubAuthController.GitHubLogin)
	r.GET("/auth/github/callback", githubAuthController.GitHubCallback)

	// 未定義のルートをホームページにリダイレクト
	r.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	})

	fmt.Println("server start")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
