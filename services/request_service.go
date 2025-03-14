package services

import (
	"giter/query"
	"giter/repositories"

	"github.com/google/go-github/github"
)

type IRequestService interface {
	GetCommits(repoName string, email string) (*query.SimpleCommits, error)
	GetRepositories(email string) ([]*github.Repository, error)
}

type RequestService struct {
	repository repositories.IRequestRepository
}

func (r *RequestService) GetCommits(repoName string, email string) (*query.SimpleCommits, error) {
	args := map[string]any{
		"USER_NAME":       email,
		"REPOSITORY_NAME": repoName,
	}
	commits, err := r.repository.GetCommits(args)
	if err != nil {
		return nil, err
	}

	// Convert commits to SimpleCommits
	sRepo := commits.ToSimpleCommits()

	// Execute GetCommittedDatesJST
	sRepo.UpdateCommittedDatesToJST()

	sRepo.RemoveDuplicateCommits()

	return sRepo, nil
}

func (r *RequestService) GetRepositories(email string) ([]*github.Repository, error) {
	return r.repository.GetRepositories(email)
}

// 新しいRequestServiceを作成する関数
func NewRequestService(repository repositories.IRequestRepository) IRequestService {
	return &RequestService{repository: repository}
}
