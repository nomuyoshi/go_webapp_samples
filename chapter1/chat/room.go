package main

import (
	"log"
	"net/http"

	"chapter1/trace"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

// room はクライアントとの接続管理、メッセージ受付を行う
type room struct {
	forward chan *message    // ルーム内にいるクライアントに転送するメッセージを保持するチャネル
	join    chan *client     // ルームに入室しようとしているクライアント管理のチャネル
	leave   chan *client     // ルームから退室しようとしているクライアント管理のチャネル
	clients map[*client]bool // ルームに入室中のクライアント管理
	tracer  trace.Tracer     // tracer はロガー
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 入室
			r.clients[client] = true
			r.tracer.Trace("新規クライアントが入室")
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			// 退室したらsendチャネルは不要（メッセージを受け取らない）なので閉じる
			close(client.send)
			r.tracer.Trace("クライアントが退室")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージ受信 mgs: ", string(msg.Message))
			// forward チャネルにメッセージが送信されてきたら
			// 入室中のクライアントのsendチャネルにメッセージを送信
			// sendチャネルに送信したら、クライアントのwriteメソッドがwebsocketに書きこむ
			for c := range r.clients {
				select {
				case c.send <- msg:
					// メッセージ送信
					r.tracer.Trace("メッセージ送信成功")
				default:
					delete(r.clients, c)
					close(c.send)
					r.tracer.Trace("メッセージ送信失敗")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキー取得失敗: ", err)
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() { r.leave <- client }()

	go client.write()
	client.read()
}
