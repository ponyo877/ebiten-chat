package folks

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var maxMessagesCount = 50

type MessageArea struct {
	messages []*Message
	x        float64
	y        float64
}

var (
	// messageAreaWidth           float32     = 450.0
	messageAreaHeight          float32     = smallFontSize*float32(maxMessagesCount) + 5
	messageAreaBackgroundColor color.Color = color.RGBA{0x22, 0x22, 0x22, 0x88}
)

// NewMessageArea creates a new MessageArea
func NewMessageArea(x, y int) *MessageArea {
	return &MessageArea{
		messages: []*Message{},
		x:        float64(x),
		y:        float64(y),
	}
}

func (s *MessageArea) AddMessage(msg *Message) {
	s.messages = append(s.messages, msg)
}

// Draw draws the MessageArea
func (s *MessageArea) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(s.x), float32(s.y), MessageAreaWidth, messageAreaHeight, messageAreaBackgroundColor, true)
	limit := min(len(s.messages), maxMessagesCount)
	for i, msg := range s.messages[len(s.messages)-limit:] {
		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(color.Gray16{0xffff})
		op.LineSpacing = smallFontSize
		op.Filter = ebiten.FilterLinear
		op.GeoM.Translate(s.x, s.y+float64(i)*smallFontSize)
		text.Draw(screen, msg.Format(), &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   smallFontSize,
		}, op)
	}
}
