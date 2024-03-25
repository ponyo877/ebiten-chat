package entity

import (
	"time"
)

type Message struct {
	messageType string
	body        *Body
	createdAt   time.Time
}

// NewMessage
func NewMessage(messageType string, body *Body, createdAt time.Time) *Message {
	return &Message{
		messageType: messageType,
		body:        body,
		createdAt:   createdAt,
	}
}

// MessageType
func (m *Message) MessageType() string {
	return m.messageType
}

// Body
func (m *Message) Body() *Body {
	return m.body
}

// CreatedAt
func (m *Message) CreatedAt() time.Time {
	return m.createdAt
}
