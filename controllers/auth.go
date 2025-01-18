package controllers

import (
	"giter/models"
	"giter/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type IAuthControler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	CurrentUser(c *gin.Context)
}

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) IAuthControler {
	// デフォルトのロガーを使用してPokemonRepositoryを初期化
	return &AuthController{db: db}
}

func (a *AuthController) Register(c *gin.Context) {
	var input RegisterInput

	// リクエストのJSONデータをRegisterInput構造体にバインドする
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーオブジェクトを作成し、データベースに保存する
	user := &models.User{Username: input.Username, Password: input.Password}
	user, err := user.Save(a.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthController) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := models.GenerateToken(input.Username, input.Password, a.db)
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
	userId, err := token.ExtractTokenId(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// ユーザーIDに基づいてユーザー情報をデータベースから取得する
	err = a.db.First(&user, userId).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}
