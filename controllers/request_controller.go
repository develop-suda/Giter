package controllers

import (
	"context"
	"fmt"
	"giter/initializer"
	services "giter/service"
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
)

// IRequestControllerインターフェースの定義
type IRequestController interface {
	GetCommits(ctx *gin.Context)
	Index(ctx *gin.Context)
	Err(ctx *gin.Context)
}

// RequestController構造体の定義
type RequestController struct {
	service services.IRequestService
	logger  zerolog.Logger
}

// GetCommitsメソッドの実装
func (c *RequestController) GetCommits(ctx *gin.Context) {
	// Implement the logic for GetCommits here
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GRAPHQL_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient("https://api.github.com/users/develop-suda/repos", httpClient)
	repositories, err := c.service.GetRepositories(client)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(repositories)

	// client := graphql.NewClient("https://api.github.com/graphql", httpClient)
}

// Indexメソッドの実装
func (c *RequestController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tmpl", nil)
}

// Errメソッドの実装
func (c *RequestController) Err(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "error.tmpl", gin.H{
		"error": "エラー用画面テスト表示",
	})
}

// NewRequestControllerファクトリーメソッドの実装
func NewRequestController(service services.IRequestService) IRequestController {
	// ロガーの初期化
	logger := initializer.DefaultLogger()
	return &RequestController{service: service, logger: logger}
}
