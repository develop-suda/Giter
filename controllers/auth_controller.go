package controllers

import (
	"giter/dto"
	"giter/initializer"
	"giter/models"
	"giter/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type IAuthControler interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	CurrentUser(c *gin.Context)
}

type AuthController struct {
	service services.IAuthService
	logger  zerolog.Logger
}

func (a *AuthController) RegisterUser(c *gin.Context) {
	var input dto.RegisterInput

	// リクエストのJSONデータをRegisterInput構造体にバインドする
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーオブジェクトを作成し、データベースに保存する
	user := &models.User{Username: input.Username, Password: input.Password}
	user, err := a.service.RegisterUser(user)
	if err != nil {
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}

func (a *AuthController) Login(c *gin.Context) {
	var input dto.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.service.Login(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (a *AuthController) CurrentUser(c *gin.Context) {
	// トークンからユーザーIDを抽出する
	userID, err := a.service.ExtractTokenId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err = a.service.CurrentUser(&user, userID)

	c.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}

func NewAuthController(service services.IAuthService) IAuthControler {
	// デフォルトのロガーを使用してPokemonRepositoryを初期化
	return &AuthController{service: service, logger: initializer.DefaultLogger()}
}
