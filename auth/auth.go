package auth

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
)

// MustAuth forces user to be authenticated
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

type authHandler struct {
	next http.Handler
}

// ③
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// 未承認時はログイン画面へリダイレクト
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		// 承認成功。ラップされたハンドラを呼び出す
		h.next.ServeHTTP(w, r)
	}
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// ④ 外部サービスからの認証結果を判定
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// ⑤ 外部サービスから取得した情報をアプリ用データとしてCookieにしこむ
	authCookieValue := objx.New(map[string]interface{}{
		"name":       user.Name,
		"avatar_url": user.AvatarURL,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	// メイン画面へリダイレクト
	w.Header()["Location"] = []string{"/"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header()["Location"] = []string{"/"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}
