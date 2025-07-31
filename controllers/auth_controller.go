package controllers

import (
	"giter/dto"
	"giter/initializer"
	"giter/models"
	"giter/services"
	"giter/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type IAuthControler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	CurrentUser(ctx *gin.Context)
	SignupView(ctx *gin.Context)
	LoginView(ctx *gin.Context)
}

type AuthController struct {
	service services.IAuthService
	logger  zerolog.Logger
}

func (a *AuthController) Register(ctx *gin.Context) {
	var input dto.RegisterInput

	// リクエストのフォームデータをRegisterInput構造体にバインドする
	if err := ctx.ShouldBind(&input); err != nil {
		a.logger.Error().Err(err).Msg("Failed to bind registration input")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// 入力検証
	if input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// ユーザーオブジェクトを作成し、データベースに保存する
	user := &models.User{Email: input.Email, Password: input.Password, IsGithubUser: false, CanSendEmail: false}
	user, err := a.service.Register(user)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to register user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := token.GenerateToken(user.ID)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to generate token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	// セキュアなクッキー設定
	ctx.SetCookie("jwt", token, 3600, "/", "", true, true)
	ctx.HTML(http.StatusOK, "mypage.tmpl", gin.H{
		"token": token,
	})
}

func (a *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput

	// フォームデータをバインドする
	if err := ctx.ShouldBind(&input); err != nil {
		a.logger.Error().Err(err).Msg("Failed to bind login input")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// 入力検証
	if input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	token, err := a.service.Login(&input)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to login user")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// セキュアなクッキー設定（localhostを削除してより汎用的に）
	ctx.SetCookie("jwt", token, 3600, "/", "", true, true)

	ctx.HTML(http.StatusOK, "mypage.tmpl", gin.H{
		"token": token,
	})
}

func (a *AuthController) CurrentUser(ctx *gin.Context) {
	// トークンからユーザーIDを抽出する
	userID, err := a.service.ExtractTokenId(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to extract user ID from token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var user models.User
	err = a.service.CurrentUser(&user, userID)
	if err != nil {
		a.logger.Error().Err(err).Uint("userID", userID).Msg("Failed to get current user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user.PrepareOutput(),
	})
}

func (a *AuthController) SignupView(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.tmpl", nil)
}

func (a *AuthController) LoginView(ctx *gin.Context) {
	// Indexメソッドの実装
	ctx.HTML(http.StatusOK, "login.tmpl", nil)
}

func NewAuthController(service services.IAuthService) IAuthControler {
	// デフォルトのロガーを使用してPokemonRepositoryを初期化
	return &AuthController{service: service, logger: initializer.DefaultLogger()}

}
