package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

// ServeHTTP はHandlerを実装
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何かしらの別のエラー
		panic(err)
	}

	// ラップされた次のハンドラを呼び出す
	h.next.ServeHTTP(w, r)
}

// MustAuth は任意のハンドラーをラップしたauthHandlerを返す
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler はソーシャルログイン処理を行う
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	// パスは「/」から始まるので次のようなスライスになっている
	// パスの形式: /auth/{action}/{provider}
	// ["", "auth", "{action}", "{provider}"]
	if len(segs) != 4 {
		log.Println("未対応のpath: ", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		// TODO: facebook, githubログイン
		authProvider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました: ", provider, "-", err)
		}
		loginURL, err := authProvider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("loginURL取得に失敗しました: ", provider, "-", err)
		}
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		authProvider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました: ", provider, "-", err)
		}
		creds, err := authProvider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("アクセストークン取得失敗: ", provider, "-", err)
		}
		user, err := authProvider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザー取得に失敗: ", provider, "-", err)
		}
		m := md5.New()
		io.WriteString(m, user.Email())
		userID := fmt.Sprintf("%x", m.Sum(nil))
		authCookieValue := objx.New(map[string]interface{}{
			"userID":    userID,
			"name":      user.Name(),
			"avatarURL": user.AvatarURL(),
			"email":     user.Email(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		// 未対応のアクションの場合404
		w.WriteHeader(http.StatusNotFound)
		log.Println("未対応のアクション: ", action)
	}
}

// logoutHandler はログアウト処理
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
