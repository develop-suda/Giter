package repositories

import (
	"giter/dto"
	"giter/initializer"
	"giter/models"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type IGithubAuthRepository interface {
	Register(user *models.User) (*models.User, error)
	Login(user *dto.LoginInput) (*models.User, error)
}

type GithubAuthRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func (g *GithubAuthRepository) Register(user *models.User) (*models.User, error) {
	// First check if user already exists
	var existingUser models.User
	err := g.db.Where("email = ?", user.Email).First(&existingUser).Error
	
	if err == nil {
		// User already exists, return the existing user
		return &existingUser, nil
	}
	
	if err != gorm.ErrRecordNotFound {
		// Some other database error occurred
		return nil, err
	}

	// User doesn't exist, create new user
	user, err = user.Save(g.db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (g *GithubAuthRepository) Login(input *dto.LoginInput) (*models.User, error) {
	var user models.User

	err := g.db.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func NewGithubAuthRepository(db *gorm.DB) IGithubAuthRepository {
	return &GithubAuthRepository{db: db, logger: initializer.DefaultLogger()}
}
