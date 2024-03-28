package drawable

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	text string
	x    float64
	y    float64
	clr  color.Color
}

func NewButton(text string, x, y float64, clr color.Color) *Button {
	return &Button{
		text,
		x,
		y,
		clr,
	}
}

func (b *Button) Contains(x, y int) bool {
	w, h := b.Bounds()
	return b.x-w/2 <= float64(x) && float64(x) <= b.x+w/2 && b.y-h/2 <= float64(y) && float64(y) <= b.y+h/2
}

func (b *Button) Bounds() (float64, float64) {
	return text.Measure(b.text, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, fontSize)
}

func (b *Button) Draw(screen *ebiten.Image) {
	w, h := b.Bounds()
	pX := b.x - w/2
	pY := b.y - h/2
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(pX), float64(pY))
	op.ColorScale.ScaleWithColor(color.Black)
	op.LineSpacing = fontSize
	op.Filter = ebiten.FilterLinear

	vector.DrawFilledRect(screen, float32(pX), float32(pY), float32(w), float32(h), b.clr, true)
	text.Draw(screen, b.text, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)
}
