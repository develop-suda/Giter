package repositories

import (
	"giter/dto"
	"giter/initializer"
	"giter/models"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	Register(user *models.User) (*models.User, error)
	Login(user *dto.LoginInput) (*models.User, error)
	CurrentUser(user *models.User, userID uint) error
}

type AuthRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func (a *AuthRepository) Register(user *models.User) (*models.User, error) {
	// パスワードをハッシュ化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	user, err = user.Save(a.db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthRepository) Login(input *dto.LoginInput) (*models.User, error) {
	var user models.User

	err := a.db.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *AuthRepository) CurrentUser(user *models.User, userID uint) error {
	// ユーザーIDに基づいてユーザー情報をデータベースから取得する
	err := a.db.First(&user, userID).Error
	if err != nil {
		return err
	}

	return nil
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{db: db, logger: initializer.DefaultLogger()}
}
