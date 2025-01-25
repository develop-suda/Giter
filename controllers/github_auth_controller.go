package controllers

import (
	"context"
	"giter/initializer"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog"
	"github.com/stretchr/objx"
)

// TODO:JWTログインと整合性をとる→DB設計、などなど
// TODO:signup後のDB登録処理
// TODO:signupとlogin処理の書き分け
// TODO:認証後の処理

type IGitHubAuthController interface {
	GitHubLogin(ctx *gin.Context)
	GitHubCallback(ctx *gin.Context)
}

type GitHubAuthController struct {
	logger zerolog.Logger
}

// GitHubLogin はGitHub認証の開始を処理します。
// このメソッドは、リクエストにプロバイダ名として "github" を設定し、
// gothicパッケージのBeginAuthHandler関数を呼び出して認証フローを開始します。
func (a *GitHubAuthController) GitHubLogin(ctx *gin.Context) {
	ctx.Request = contextWithProviderName(ctx, "github")
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// GitHubCallback はGitHub認証のコールバックを処理するメソッドです。
// 認証が成功した場合、ユーザー情報を取得し、クッキーに保存してメイン画面にリダイレクトします。
// 認証に失敗した場合、エラーメッセージを返します。
func (a *GitHubAuthController) GitHubCallback(ctx *gin.Context) {
	// ユーザー認証を完了し、ユーザー情報を取得
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		// 認証に失敗した場合、エラーメッセージを返す
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザー情報をクッキーに保存
	authCookieValue := objx.New(map[string]interface{}{
		"name":       user.Name,
		"avatar_url": user.AvatarURL,
	}).MustBase64()
	ctx.SetCookie("auth", authCookieValue, 0, "/", "", false, true)

	// メイン画面へリダイレクト
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func contextWithProviderName(c *gin.Context, provider string) *http.Request {
	return c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
}

func NewGitHubAuthController() IGitHubAuthController {
	return &GitHubAuthController{logger: initializer.DefaultLogger()}
}
