package folks

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SpeechBubble struct {
	message   *Message
	pX1000    float64
	pY1000    float64
	createdAt int64
}

var (
	speechBubbleWidthUnit    float32 = 6.0
	speechBubbleHeight       float32 = 20.0
	speechBubbleAltitudeUnit float64 = 0.00005
	pointScale               int     = 1000
)

// NewSpeechBubble creates a new SpeechBubble
func NewSpeechBubble(message *Message, pointX, pointY int) (*SpeechBubble, error) {
	createdAt := time.Now().Unix()
	return &SpeechBubble{
		message:   message,
		pX1000:    float64(pointX * pointScale),
		pY1000:    float64(pointY * pointScale),
		createdAt: createdAt,
	}, nil
}

// Draw draws the SpeechBubble
func (s *SpeechBubble) Draw(screen *ebiten.Image, now time.Time) {
	speechBubbleWidth := speechBubbleWidthUnit * s.message.Size()
	pX := s.pX1000 / float64(pointScale)
	pY := s.pY1000/float64(pointScale) - s.altitude(now)
	vector.DrawFilledRect(screen, float32(pX), float32(pY), speechBubbleWidth, speechBubbleHeight, color.White, true)

	// message over pangle
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(pX), float64(pY))
	op.ColorScale.ScaleWithColor(color.Black)
	op.LineSpacing = smallFontSize
	op.Filter = ebiten.FilterLinear
	text.Draw(screen, s.message.Content(), &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   smallFontSize,
	}, op)
}

// Altitude returns the altitude of the SpeechBubble
func (s *SpeechBubble) altitude(now time.Time) float64 {
	return float64(s.message.ElapsedMicro(now)) * speechBubbleAltitudeUnit
}
