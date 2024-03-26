package entity

import (
	"time"
)

type SocketMessage struct {
	messageType string
	body        *Body
	createdAt   time.Time
}

func NewSocketMessage(messageType string, body *Body, createdAt time.Time) *SocketMessage {
	return &SocketMessage{
		messageType: messageType,
		body:        body,
		createdAt:   createdAt,
	}
}

func (m *SocketMessage) MessageType() string {
	return m.messageType
}

func (m *SocketMessage) Body() *Body {
	return m.body
}

func (m *SocketMessage) CreatedAt() time.Time {
	return m.createdAt
}
