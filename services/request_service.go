package services

import (
	"giter/query"
	"giter/repositories"

	"github.com/google/go-github/github"
)

type IRequestService interface {
	GetCommits(args map[string]interface{}) (*query.GitHubQuery, error)
	GetRepositories(username string) ([]*github.Repository, error)
}

type RequestService struct {
	repository repositories.IRequestRepository
}

func (r *RequestService) GetCommits(args map[string]interface{}) (*query.GitHubQuery, error) {
	commits, err := r.repository.GetCommits(args)
	if err != nil {
		return nil, err
	}

	// Execute GetCommittedDatesJST
	commits.UpdateCommittedDatesToJST()

	return commits, nil
}

func (r *RequestService) GetRepositories(username string) ([]*github.Repository, error) {
	return r.repository.GetRepositories(username)
}

// 新しいRequestServiceを作成する関数
func NewRequestService(repository repositories.IRequestRepository) IRequestService {
	return &RequestService{repository: repository}
}
