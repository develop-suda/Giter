package controllers

import (
	"fmt"
	"giter/initializer"
	"giter/query"
	"giter/services"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/rs/zerolog"

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

func (c *RequestController) GetCommits(ctx *gin.Context) {
	repositories, err := c.GetRepositories(ctx)
	if err != nil {
		return
	}

	var commits []query.SimpleCommits

	for i, v := range repositories {
		repo := *v.Name
		fmt.Printf("Repository %d: %s\n", i, repo)
		commit, err := c.service.GetCommits(repo, "develop-suda")
		commits = append(commits, *commit)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	result := query.ToCommits(&commits)

	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		"commits": result,
	})
}

func (c *RequestController) GetRepositories(ctx *gin.Context) ([]*github.Repository, error) {

	email := "develop-suda"

	repos, err := c.service.GetRepositories(email)
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
