package websocket

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ponyo877/folks-ui/entity"
	"github.com/ponyo877/folks-ui/websocket/presenter"
	"nhooyr.io/websocket" // wasm対応のためgorilla/websockeから変更
	"nhooyr.io/websocket/wsjson"
)

type WebSocket struct {
	conn *websocket.Conn
}

func NewWebSocket(scheme, host, path string) (*WebSocket, error) {
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	conn, _, err := websocket.Dial(context.Background(), u.String(), nil)
	if err != nil {
		fmt.Printf("Websocketサーバへの接続に失敗しました: %v\n", err)
		return nil, err
	}
	return &WebSocket{conn: conn}, nil
}

func (w *WebSocket) Send(message *entity.SocketMessage) error {
	return wsjson.Write(context.Background(), w.conn, presenter.MarshalMessage(message))
}

func (w *WebSocket) Close() error {
	return w.conn.Close(websocket.StatusNormalClosure, "exit")
}

func (w *WebSocket) Receive(f func(*entity.SocketMessage)) {
	var messagePresenter presenter.MessagePresenter
	for {
		if err := wsjson.Read(context.Background(), w.conn, &messagePresenter); err != nil {
			fmt.Printf("Messageの読み込みに失敗しました: %v\n", err)
			return
		}
		f(messagePresenter.Unmarshal())
	}
}
