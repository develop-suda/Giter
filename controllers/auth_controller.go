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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーオブジェクトを作成し、データベースに保存する
	user := &models.User{Email: input.Email, Password: input.Password, IsGithubUser: false}
	user, err := a.service.Register(user)
	if err != nil {
	}

	token, err := token.GenerateToken(user.ID)
	if err != nil {
	}

	ctx.SetCookie("jwt", token, 3600, "/", "localhost", true, true)
	ctx.HTML(http.StatusOK, "mypage.tmpl", gin.H{
		"token": token,
	})
}

func (a *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput

	// フォームデータをバインドする
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.service.Login(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ctx.SetCookie("jwt", cookie, 3600, "/", "localhost", false, true)
	ctx.SetCookie("jwt", token, 3600, "/", "localhost", true, true)

	ctx.HTML(http.StatusOK, "mypage.tmpl", gin.H{
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
