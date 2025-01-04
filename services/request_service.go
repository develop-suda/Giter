package services

import (
	"giter/query"
	"giter/repositories"

	"github.com/google/go-github/github"
	"github.com/hasura/go-graphql-client"
)

type IRequestService interface {
	GetCommits(client *graphql.Client, variables map[string]interface{}) (*query.GitHubQuery, error)
	GetRepositories(client *github.Client, username string) ([]*github.Repository, error)
}

type RequestService struct {
	repository repositories.IRequestRepository
}

func (r *RequestService) GetCommits(client *graphql.Client, variables map[string]interface{}) (*query.GitHubQuery, error) {
	// 実装をここに追加
	commits, err := r.repository.GetCommits(client, variables)
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func (r *RequestService) GetRepositories(client *github.Client, username string) ([]*github.Repository, error) {
	return r.repository.GetRepositories(client, username)
}

// 新しいRequestServiceを作成する関数
func NewRequestService(repository repositories.IRequestRepository) IRequestService {
	return &RequestService{repository: repository}
}
