package folks

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	lifespan         int64 = 5000
	maxContentLength int   = 100
)

type Message struct {
	characterID string
	content     string
	createdAt   time.Time
}

// NewMessage creates a new Message
func NewMessage(characterID, content string, createdAt time.Time) (*Message, error) {
	if utf8.RuneCountInString(content) > maxContentLength {
		return nil, fmt.Errorf("content is too long")
	}
	return &Message{
		characterID,
		content,
		createdAt,
	}, nil
}

// CharacterID returns the characterID of the Message
func (m *Message) CharacterID() string {
	return m.characterID
}

// Content returns the content of the Message
func (m *Message) Content() string {
	return m.content
}

func (m *Message) Format() string {
	escaped := strings.Replace(m.content, "\n", " ", -1)
	return fmt.Sprintf("(id_%3s) %v [%02d/%02d %02d:%02d:%02d]", m.characterID[:3], escaped, m.createdAt.Month(), m.createdAt.Day(), m.createdAt.Hour(), m.createdAt.Minute(), m.createdAt.Second())
}

// Size returns the size of the Message
func (m *Message) Size() float32 {
	return float32(len(m.content))
}

func (m *Message) ElapsedMilli(now time.Time) int64 {
	return now.UnixMilli() - m.createdAt.UnixMilli()
}

// IsExpired returns true if the Message is expired
func (m *Message) IsExpired(now time.Time) bool {
	return m.ElapsedMilli(now) > lifespan
}
