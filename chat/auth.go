package main

import (
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

// ServeHTTP はHandlerを実装
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
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
	// 末尾に「/」があってもなくても動くようにtrim
	path := strings.TrimSuffix(r.URL.Path, "/")
	segs := strings.Split(path, "/")
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
		log.Println("TODO: ログイン処理", provider)
	default:
		// 未対応のアクションの場合404
		w.WriteHeader(http.StatusNotFound)
		log.Println("未対応のアクション: ", action)
	}
}
