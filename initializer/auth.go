package initializer

import (
	"giter/auth"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/gorilla/pat"
	// "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/stretchr/objx"
)

func init() {
	/*
		// gothで利用するCookieの設定を上書きする場合
		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
		store.MaxAge(86400 * 60) // セッション期限の設定(60日) デフォルトでは30日
		store.Options.Secure = true // Cookieのセキュア設定 デフォルトではfalse
		gothic.Store = store // 上書きする
	*/

	// ①
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/auth/github/callback"),
	)
}

type templateHandler struct {
	filename string
	once     sync.Once
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(
			template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := make(map[string]interface{})

	// ⑥ アプリ用Cookieからユーザー情報を取得する
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	// ② patを使ってルーティング設定
	p := pat.New()
	p.Get("/auth/{provider}/callback", auth.CallbackHandler)
	p.Get("/auth/{provider}", gothic.BeginAuthHandler)
	p.Get("/logout", auth.LogoutHandler)
	p.Add("GET", "/login", &templateHandler{filename: "login.html"})
	p.Add("GET", "/", auth.MustAuth(&templateHandler{filename: "index.html"}))

	// WEBサーバを起動
	log.Fatal(http.ListenAndServe(":3000", p))
}
