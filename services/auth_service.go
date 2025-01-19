package services

import (
	"giter/dto"
	"giter/initializer"
	"giter/models"
	"giter/repositories"
	"giter/utils/token"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	RegisterUser(user *models.User) (*models.User, error)
	Login(user *dto.LoginInput) (string, error)
	CurrentUser(user *models.User, userID uint) error
	ExtractTokenId(ctx *gin.Context) (uint, error)
}

type AuthService struct {
	repository repositories.IAuthRepository
	logger     zerolog.Logger
}

func (a *AuthService) RegisterUser(user *models.User) (*models.User, error) {
	return a.repository.RegisterUser(user)
}

func (a *AuthService) Login(input *dto.LoginInput) (string, error) {
	// ユーザ名で検索をかける
	// DBから1行取得
	user, err := a.repository.Login(input)
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

func (a *AuthService) CurrentUser(user *models.User, userID uint) error {
	return a.repository.CurrentUser(user, userID)
}

func (a *AuthService) ExtractTokenId(ctx *gin.Context) (uint, error) {
	return token.ExtractTokenId(ctx)
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
	// デフォルトのロガーを使用してPokemonRepositoryを初期化
	return &AuthService{repository: repository, logger: initializer.DefaultLogger()}
}
