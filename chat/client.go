package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client はチャットしている1人のユーザーを表す
type client struct {
	socket   *websocket.Conn        // WebSocket
	send     chan *message          // メッセージが送られてくるチャネル
	room     *room                  // クライアントが入っているルーム
	userData map[string]interface{} // ユーザー情報
}

// read はwebsocketからメッセージを読み込む
func (c *client) read() {
	for {
		var msg message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now().Format("2006/01/02 15:04:05")
			msg.Name = c.userData["name"].(string)
			// websocketから読み込んだメッセージをroomのforwardチャネルに送信
			c.room.forward <- &msg
		} else {
			break
		}
	}
}

// write はsendチャネルに溜まったメッセージをwebsocketに書き込む
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}

	c.socket.Close()
}
