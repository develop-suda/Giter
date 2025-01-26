package controllers

import (
	"context"
	"giter/initializer"
	"giter/models"
	"giter/services"
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
// TODO:GitHubログイン後のJWT???

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
// 認証が成功した場合、ユーザー情報を取得し、クッキーに保存してメイン画面にリダイレクトします。
// 認証に失敗した場合、エラーメッセージを返します。
func (g *GithubAuthController) GitHubCallback(ctx *gin.Context) {
	// ユーザー認証を完了し、ユーザー情報を取得
	userData, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		// 認証に失敗した場合、エラーメッセージを返す
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーオブジェクトを作成し、データベースに保存する
	user := &models.User{Email: userData.Email, Password: "", IsGithubUser: true, CanSendEmail: false}
	user, err = g.service.Register(user)
	if err != nil {
	}

	// ユーザー情報をクッキーに保存
	authCookieValue := objx.New(map[string]interface{}{
		"name":       userData.Name,
		"avatar_url": userData.AvatarURL,
	}).MustBase64()
	ctx.SetCookie("auth", authCookieValue, 0, "/", "", false, true)

	// メイン画面へリダイレクト
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func contextWithProviderName(c *gin.Context, provider string) *http.Request {
	return c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
}

func NewGithubAuthController(service services.IGithubAuthService) IGithubAuthController {
	return &GithubAuthController{service: service, logger: initializer.DefaultLogger()}
}
