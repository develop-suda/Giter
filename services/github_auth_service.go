package services

import (
	"giter/dto"
	"giter/initializer"
	"giter/models"
	"giter/repositories"
	"giter/utils/token"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type IGithubAuthService interface {
	Register(user *models.User) (*models.User, error)
	Login(user *dto.LoginInput) (string, error)
}

type GithubAuthService struct {
	repository repositories.IGithubAuthRepository
	logger     zerolog.Logger
}

func (g *GithubAuthService) Register(user *models.User) (*models.User, error) {
	return g.repository.Register(user)
}

func (g *GithubAuthService) Login(input *dto.LoginInput) (string, error) {
	// ユーザ名で検索をかける
	// DBから1行取得
	user, err := g.repository.Login(input)
	if err != nil {
		return "", err
	}

	// パスワード検証
	// パスワードが正しい場合nilを返却する
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", err
	}

	token, err := token.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewGithubAuthService(repository repositories.IGithubAuthRepository) IGithubAuthService {
	// デフォルトのロガーを使用してPokemonRepositoryを初期化
	return &GithubAuthService{repository: repository, logger: initializer.DefaultLogger()}
}
