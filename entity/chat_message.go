package entity

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	lifespan         int64 = 10000
	MaxContentLength int   = 20
)

type ChatMessage struct {
	characterID string
	content     string
	createdAt   time.Time
}

// NewChatMessage creates a new Message
func NewChatMessage(characterID, content string, createdAt time.Time) (*ChatMessage, error) {
	if utf8.RuneCountInString(content) > MaxContentLength {
		return nil, fmt.Errorf("content is too long")
	}
	return &ChatMessage{
		characterID,
		content,
		createdAt,
	}, nil
}

// CharacterID returns the characterID of the ChatMessage
func (m *ChatMessage) CharacterID() string {
	return m.characterID
}

// Content returns the content of the ChatMessage
func (m *ChatMessage) Content() string {
	return m.content
}

func (m *ChatMessage) Format() string {
	escaped := strings.Replace(m.content, "\n", " ", -1)
	return fmt.Sprintf("(id_%3s) %v [%02d/%02d %02d:%02d:%02d]", m.characterID[:3], escaped, m.createdAt.Month(), m.createdAt.Day(), m.createdAt.Hour(), m.createdAt.Minute(), m.createdAt.Second())
}

// Size returns the size of the ChatMessage
func (m *ChatMessage) Size() float32 {
	return float32(len(m.content))
}

func (m *ChatMessage) ElapsedMilli(now time.Time) int64 {
	return now.UnixMilli() - m.createdAt.UnixMilli()
}

// IsExpired returns true if the ChatMessage is expired
func (m *ChatMessage) IsExpired(now time.Time) bool {
	return m.ElapsedMilli(now) > lifespan
}
