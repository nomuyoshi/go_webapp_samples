package main

// room はクライアントとの接続管理、メッセージ受付を行う
type room struct {
	forward chan []byte      // ルーム内にいるクライアントに転送するメッセージを保持するチャネル
	join    chan *client     // ルームに入室しようとしているクライアント管理のチャネル
	leave   chan *client     // ルームから退室しようとしているクライアント管理のチャネル
	clients map[*client]bool // ルームに入室中のクライアント管理
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 入室
			r.clients[client] = true
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			// 退室したらsendチャネルは不要（メッセージを受け取らない）なので閉じる
			close(client.send)
		case msg := <-r.forward:
			// forward チャネルにメッセージが送信されてきたら
			// 入室中のクライアントのsendチャネルにメッセージを送信
			// sendチャネルに送信したら、クライアントのwriteメソッドがwebsocketに書きこむ
			for c := range r.clients {
				select {
				case c.send <- msg:
					// メッセージ送信
				default:
					delete(r.clients, c)
					close(c.send)
				}
			}
		}
	}
}
