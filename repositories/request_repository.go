package repositories

import (
	"context"
	"fmt"
	"giter/initializer"
	"giter/query"

	"github.com/google/go-github/github"
	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
)

// IRequestRepository インターフェースの定義
type IRequestRepository interface {
	// ポケモンを検索するメソッド
	GetCommits(client *graphql.Client, variables map[string]interface{}) (*query.GitHubQuery, error)
	GetRepositories(client *github.Client, username string) ([]*github.Repository, error)
}

// RequestRepository 構造体の定義
type RequestRepository struct {
	logger zerolog.Logger
}

// GetCommits implements IRequestRepository.
func (r *RequestRepository) GetCommits(client *graphql.Client, variables map[string]interface{}) (*query.GitHubQuery, error) {
	var commitsQuery query.GitHubQuery
	err := client.Query(context.Background(), &commitsQuery, variables)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &commitsQuery, nil
}

func (r *RequestRepository) GetRepositories(client *github.Client, username string) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.List(context.Background(), username, nil)
	if err != nil {
		fmt.Println(client)
		return nil, err
	}
	return repos, nil
}

// NewRequestRepository コンストラクタ関数
func NewRequestRepository() IRequestRepository {
	// デフォルトのロガーを使用してRequestRepositoryを初期化
	return &RequestRepository{logger: initializer.DefaultLogger()}
}
