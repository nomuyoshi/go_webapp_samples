package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"webapp_samples/trace"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP はHTTPリクエストを処理する
// インターフェース Handler を実装
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// テンプレートのコンパイル
	// 1度コンパイルすれば良いのでsync.Onceを使っている
	// sync.Onceなら複数のゴルーチンからServeHTTPが呼ばれても1度しか実行されない
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	if err := t.templ.Execute(w, data); err != nil {
		panic(err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET_KEY"), "http://localhost:8080/auth/callback/google"),
	)
	// r := newRoom(UseAuthAvatar)
	// r := newRoom(UseGravatarAvatar)
	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)

	// ログイン後
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/upload", MustAuth(&templateHandler{filename: "upload.html"}))
	http.Handle("/uploader", MustAuth(uploadHandler{}))
	http.Handle("/avatars/", MustAuth(
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars")))))
	http.Handle("/room", r)
	http.HandleFunc("/logout", logoutHandler)
	go r.run()

	log.Println("Webサーバー起動 port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
