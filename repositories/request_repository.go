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
	GetCommits(args map[string]interface{}) (*query.GitHubQuery, error)
	GetRepositories(username string) ([]*github.Repository, error)
}

// RequestRepository 構造体の定義
type RequestRepository struct {
	RESTClient    *github.Client
	GraphQLClient *graphql.Client
	logger        zerolog.Logger
}

// GetCommits implements IRequestRepository.
func (r *RequestRepository) GetCommits(args map[string]interface{}) (*query.GitHubQuery, error) {
	var commitsQuery query.GitHubQuery
	err := r.GraphQLClient.Query(context.Background(), &commitsQuery, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &commitsQuery, nil
}

func (r *RequestRepository) GetRepositories(username string) ([]*github.Repository, error) {
	// オプションをつけないとアルファベット順に取得する
	opts := &github.RepositoryListOptions{Sort: "pushed", ListOptions: github.ListOptions{PerPage: 100}}
	repos, _, err := r.RESTClient.Repositories.List(context.Background(), username, opts)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// NewRequestRepository コンストラクタ関数
func NewRequestRepository(RClient *github.Client, GClient *graphql.Client) IRequestRepository {
	// デフォルトのロガーを使用してRequestRepositoryを初期化
	return &RequestRepository{logger: initializer.DefaultLogger(), RESTClient: RClient, GraphQLClient: GClient}
}
