package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
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

	if err := t.templ.Execute(w, r); err != nil {
		panic(err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()

	log.Println("Webサーバー起動 port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
