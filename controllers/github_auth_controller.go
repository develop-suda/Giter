package controllers

import (
	"context"
	"giter/initializer"
	"giter/models"
	"giter/services"
	"giter/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog"
)

// Fixed: JWT login consistency achieved through standardized JWT token usage
// Fixed: DB registration process implemented for GitHub signup
// Fixed: GitHub login now uses JWT tokens consistently
// Note: Additional features like user profile management could be added in the future

type IGithubAuthController interface {
	GitHubLogin(ctx *gin.Context)
	GitHubCallback(ctx *gin.Context)
}

type GithubAuthController struct {
	service services.IGithubAuthService
	logger  zerolog.Logger
}

// GitHubLogin はGitHub認証の開始を処理します。
// このメソッドは、リクエストにプロバイダ名として "github" を設定し、
// gothicパッケージのBeginAuthHandler関数を呼び出して認証フローを開始します。
func (g *GithubAuthController) GitHubLogin(ctx *gin.Context) {
	ctx.Request = contextWithProviderName(ctx, "github")
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// GitHubCallback はGitHub認証のコールバックを処理するメソッドです。
// 認証が成功した場合、ユーザー情報を取得し、JWTトークンを生成してクッキーに保存し、メイン画面にリダイレクトします。
// 認証に失敗した場合、エラーメッセージを返します。
func (g *GithubAuthController) GitHubCallback(ctx *gin.Context) {
	// ユーザー認証を完了し、ユーザー情報を取得
	userData, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		// 認証に失敗した場合、エラーメッセージを返す
		g.logger.Error().Err(err).Msg("GitHub authentication failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "GitHub authentication failed"})
		return
	}

	// 入力検証
	if userData.Email == "" {
		g.logger.Error().Msg("GitHub user has no email address")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "GitHub account must have a public email address"})
		return
	}

	// ユーザーオブジェクトを作成し、データベースに保存する
	user := &models.User{Email: userData.Email, Password: "", IsGithubUser: true, CanSendEmail: false}
	user, err = g.service.Register(user)
	if err != nil {
		g.logger.Error().Err(err).Str("email", userData.Email).Msg("Failed to register GitHub user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
		return
	}

	// JWTトークンを生成（通常認証と統一）
	token, err := token.GenerateToken(user.ID)
	if err != nil {
		g.logger.Error().Err(err).Uint("userID", user.ID).Msg("Failed to generate JWT token for GitHub user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	// セキュアなクッキー設定でJWTトークンを保存
	ctx.SetCookie("jwt", token, 3600, "/", "", true, true)

	// メイン画面へリダイレクト
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func contextWithProviderName(c *gin.Context, provider string) *http.Request {
	return c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
}

func NewGithubAuthController(service services.IGithubAuthService) IGithubAuthController {
	return &GithubAuthController{service: service, logger: initializer.DefaultLogger()}
}
