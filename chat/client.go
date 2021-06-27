package main

import (
	"github.com/gorilla/websocket"
)

// client はチャットしている1人のユーザーを表す
type client struct {
	socket *websocket.Conn // WebSocket
	send   chan []byte     // メッセージが送られてくるチャネル
	room   *room           // クライアントが入っているルーム
}

// read はwebsocketからメッセージを読み込む
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err != nil {
			// websocketから読み込んだメッセージをroomのforwardチャネルに送信
			c.room.forward <- msg
		} else {
			break
		}
	}
}

// write はsendチャネルに溜まったメッセージをwebsocketに書き込む
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}

	c.socket.Close()
}
