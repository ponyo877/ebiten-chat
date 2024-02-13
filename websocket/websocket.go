package websocket

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ponyo877/folks-ui/entity"
	"nhooyr.io/websocket" // wasm対応のためgorilla/websockeから変更
	"nhooyr.io/websocket/wsjson"
)

type WebSocket struct {
	conn *websocket.Conn
}

func NewWebSocket(host, path string) (*WebSocket, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: path}
	conn, _, err := websocket.Dial(context.Background(), u.String(), nil)
	if err != nil {
		fmt.Printf("Websocketサーバへの接続に失敗しました: %v\n", err)
		return nil, err
	}
	return &WebSocket{conn: conn}, nil
}

func (w *WebSocket) Send(message *entity.Message) error {
	return wsjson.Write(context.Background(), w.conn, MarshalMessage(message))
}

func (w *WebSocket) Candidate(ICECandidate []byte) error {
	signal := NewSignal("candidate", ICECandidate)
	if err := wsjson.Write(context.Background(), w.conn, signal); err != nil {
		return err
	}
	return nil
}

func (w *WebSocket) ReadMessage() ([]byte, error) {
	_, bin, err := w.conn.Read(context.Background())
	return bin, err
}

func (w *WebSocket) Receive(f func(*entity.Message)) {
	var messagePresenter MessagePresenter
	for {
		if err := wsjson.Read(context.Background(), w.conn, &messagePresenter); err != nil {
			fmt.Printf("Messageの読み込みに失敗しました: %v\n", err)
			return
		}
		f(messagePresenter.Unmarshal())
	}
}
