package main

// room はクライアントとの接続管理、メッセージ受付を行う
type room struct {
	forward chan []byte // ルーム内にいるクライアントに転送するメッセージを保持するチャネル
}
