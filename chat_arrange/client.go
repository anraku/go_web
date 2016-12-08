package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// clientはチャットを行っている1人のユーザを表しています。
type client struct {
	// socketはこのクライアントのためのWebSocketです。
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです。
	send chan *message
	// room はこのクライアントが参加しているチャットルームです。
	room *room
	// userDataはユーザーに関する情報を保持します
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			//時間のフォーマットを定義
			// const layout = "2006-01-02 15:04:05"
			// var now = time.Now()
			// var dateStr = now.Format(layout)
			msg.When = time.Now()
			if err != nil {
				panic(err)
			}
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
