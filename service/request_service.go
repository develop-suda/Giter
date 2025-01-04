package services

import (
	"giter/query"
	"giter/repositories"

	"github.com/hasura/go-graphql-client"
)

type IRequestService interface {
	GetCommits(client *graphql.Client, variables map[string]interface{}) (*query.GitHubQuery, error)
	GetRepositories(client *graphql.Client) (*query.WelcomeElement, error)
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

func (r *RequestService) GetRepositories(client *graphql.Client) (*query.WelcomeElement, error) {
	return r.repository.GetRepositories(client)
}

// 新しいRequestServiceを作成する関数
func NewRequestService(repository repositories.IRequestRepository) IRequestService {
	return &RequestService{repository: repository}
}
