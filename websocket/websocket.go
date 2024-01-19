package websocket

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/ponyo877/folks-ui/entity"
)

type WebSocket struct {
	conn *websocket.Conn
}

func NewWebSocket(host, path string) (*WebSocket, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Printf("Websocketサーバへの接続に失敗しました: %v\n", err)
		return nil, err
	}
	// defer conn.Close()
	return &WebSocket{conn: conn}, nil
}

func (w *WebSocket) Send(message *entity.Message) error {
	return w.conn.WriteJSON(MarshalMessage(message))
}

func (w *WebSocket) Receive(f func(*entity.Message)) {
	var messagePresenter MessagePresenter
	for {
		if err := w.conn.ReadJSON(&messagePresenter); err != nil {
			fmt.Printf("Messageの読み込みに失敗しました: %v\n", err)
			return
		}
		f(messagePresenter.Unmarshal())
		// websocket.JSON.Receive(ws, &rcvMsg)
		// log.Printf("Receive data=%#v\n", rcvMsg)
	}
}
