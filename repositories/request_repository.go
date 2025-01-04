package repositories

import (
	"context"
	"fmt"
	"giter/initializer"
	"giter/query"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
)

// IPokemonRepository インターフェースの定義
type IRequestRepository interface {
	// ポケモンを検索するメソッド
	GetCommits(client *graphql.Client, variables map[string]interface{}) (*query.GitHubQuery, error)
	GetRepositories(client *graphql.Client) (*query.WelcomeElement, error)
}

// PokemonRepository 構造体の定義
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

func (r *RequestRepository) GetRepositories(client *graphql.Client) (*query.WelcomeElement, error) {
	var repositoriesQuery query.WelcomeElement
	err := client.Query(context.Background(), &repositoriesQuery, nil)
	if err != nil {
		return nil, err
	}

	return &repositoriesQuery, nil
}

// NewPokemonRepository コンストラクタ関数
func NewRequestRepository() IRequestRepository {
	// デフォルトのロガーを使用してPokemonRepositoryを初期化
	return &RequestRepository{logger: initializer.DefaultLogger()}
}
