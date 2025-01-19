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
	RegisterUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	CurrentUser(ctx *gin.Context)
	LoginView(ctx *gin.Context)
}

type AuthController struct {
	service services.IAuthService
	logger  zerolog.Logger
}

func (a *AuthController) RegisterUser(ctx *gin.Context) {
	var input dto.RegisterInput

	// リクエストのJSONデータをRegisterInput構造体にバインドする
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーオブジェクトを作成し、データベースに保存する
	user := &models.User{Username: input.Username, Password: input.Password}
	user, err := a.service.RegisterUser(user)
	if err != nil {
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}

func (a *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.service.Login(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (a *AuthController) CurrentUser(ctx *gin.Context) {
	// トークンからユーザーIDを抽出する
	userID, err := a.service.ExtractTokenId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err = a.service.CurrentUser(&user, userID)

	ctx.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}

func (a *AuthController) LoginView(ctx *gin.Context) {
	// Indexメソッドの実装
	ctx.HTML(http.StatusOK, "login.tmpl", nil)
}

func NewAuthController(service services.IAuthService) IAuthControler {
	// デフォルトのロガーを使用してPokemonRepositoryを初期化
	return &AuthController{service: service, logger: initializer.DefaultLogger()}

}
