package folks

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/folks-ui/entity"
)

type SpeechBubble struct {
	message   *entity.ChatMessage
	x         float64
	y         float64
	createdAt int64
}

var (
	speechBubbleAltitudeUnit float64 = 0.03
)

// NewSpeechBubble creates a new SpeechBubble
func NewSpeechBubble(message *entity.ChatMessage, x, y int) (*SpeechBubble, error) {
	createdAt := time.Now().Unix()
	return &SpeechBubble{
		message:   message,
		x:         float64(x),
		y:         float64(y),
		createdAt: createdAt,
	}, nil
}

// Draw draws the SpeechBubble
func (s *SpeechBubble) Draw(screen *ebiten.Image, now time.Time) {
	w, h := text.Measure(s.message.Content(), &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, fontSize)
	pX := s.x - w/2
	pY := s.y - h - s.altitude(now)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(pX), float64(pY))
	op.ColorScale.ScaleWithColor(color.Black)
	op.LineSpacing = fontSize
	op.Filter = ebiten.FilterLinear

	vector.DrawFilledRect(screen, float32(pX), float32(pY), float32(w), float32(h), color.White, true)
	text.Draw(screen, s.message.Content(), &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)
}

// Altitude returns the altitude of the SpeechBubble
func (s *SpeechBubble) altitude(now time.Time) float64 {
	return float64(s.message.ElapsedMilli(now)) * speechBubbleAltitudeUnit
}
