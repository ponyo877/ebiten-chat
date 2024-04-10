package drawable

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Text struct {
	x     float64
	y     float64
	size  float64
	text  string
	bgclr color.Color
	font  *text.GoTextFaceSource
}

func NewText(x, y, size float64, text string, bgclr color.Color, font *text.GoTextFaceSource) *Text {
	return &Text{
		x, y,
		size,
		text,
		bgclr,
		font,
	}
}

func (t *Text) Bounds() (float64, float64) {
	return text.Measure(t.text, &text.GoTextFace{
		Source: t.font,
		Size:   t.size,
	}, t.size)
}

func (t *Text) Contains(x, y int, center bool) bool {
	w, h := t.Bounds()
	if center {
		return t.x-w/2 <= float64(x) && float64(x) <= t.x+w/2 && t.y <= float64(y) && float64(y) <= t.y+h
	}
	return t.x <= float64(x) && float64(x) <= w && t.y <= float64(y) && float64(y) <= t.y+h
}

func (t *Text) Draw(screen *ebiten.Image, center bool) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(t.bgclr)
	op.LineSpacing = t.size
	op.Filter = ebiten.FilterLinear
	w, _ := t.Bounds()
	x := float64(t.x)
	if center {
		x = float64(t.x) - w/2
	}
	op.GeoM.Translate(x, float64(t.y))
	text.Draw(screen, t.text, &text.GoTextFace{
		Source: t.font,
		Size:   t.size,
	}, op)
}
