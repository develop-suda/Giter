package controllers

import (
	"context"
	"fmt"
	"giter/initializer"
	"giter/services"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
)

// IRequestControllerインターフェースの定義
type IRequestController interface {
	GetCommits(ctx *gin.Context)
	GetRepositories(ctx *gin.Context) ([]*github.Repository, error)
	Index(ctx *gin.Context)
	Err(ctx *gin.Context)
}

// RequestController構造体の定義
type RequestController struct {
	service services.IRequestService
	logger  zerolog.Logger
}

// GetCommitsメソッドの実装
// TODO: リポジトリ一覧取得してそれでループ回す
// TODO: client作るところ別関数にする
// TODO: log出力させる
func (c *RequestController) GetCommits(ctx *gin.Context) {
	repositories, err := c.GetRepositories(ctx)
	if err != nil {
		return
	}

	fmt.Println(repositories)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GRAPHQL_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	args := map[string]any{
		"USER_NAME":       "develop-suda",
		"REPOSITORY_NAME": "Giter",
	}
	commits, err := c.service.GetCommits(client, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(commits)

}

func (c *RequestController) GetRepositories(ctx *gin.Context) ([]*github.Repository, error) {
	client := github.NewClient(nil)
	username := "develop-suda"

	repos, err := c.service.GetRepositories(client, username)
	if err != nil {
		return nil, err
	}
	return repos, nil
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
