package controllers

import (
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
		c.logger.Error().Err(err).Msg("Failed to get repositories")
		ctx.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Failed to retrieve repositories",
		})
		return
	}

	var commits []query.SimpleCommits

	for i, v := range repositories {
		repo := *v.Name
		c.logger.Debug().Int("index", i).Str("repo", repo).Msg("Processing repository")
		commit, err := c.service.GetCommits(repo, "develop-suda") // TODO: make username configurable
		if err != nil {
			c.logger.Error().Err(err).Str("repo", repo).Msg("Failed to get commits for repository")
			continue // Skip this repository but continue with others
		}
		commits = append(commits, *commit)
	}

	result := query.ToCommits(&commits)

	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		"commits": result,
	})
}

func (c *RequestController) GetRepositories(ctx *gin.Context) ([]*github.Repository, error) {
	// TODO: Get username from authenticated user context instead of hardcoding
	email := "develop-suda"

	repos, err := c.service.GetRepositories(email)
	if err != nil {
		c.logger.Error().Err(err).Str("email", email).Msg("Failed to get repositories")
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
