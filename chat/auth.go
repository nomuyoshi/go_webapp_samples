package main

import "net/http"

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
