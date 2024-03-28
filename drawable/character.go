package drawable

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/ponyo877/folks-ui/entity"
)

type Character struct {
	id    string
	x     int
	y     int
	name  string
	imgid int
	dir   entity.Dir
}

func NewCharacter(id string, x, y int, name string, imgid int, dir entity.Dir) *Character {
	return &Character{
		id,
		x,
		y,
		name,
		imgid,
		dir,
	}
}

func (c *Character) Point() (int, int) {
	return c.x, c.y
}

func (c *Character) Dir() entity.Dir {
	return c.dir
}

func (c *Character) IsMine(myID string) bool {
	return c.id == myID
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) image() *ebiten.Image {
	return CharacterImage[c.imgid]
}

func (c *Character) Bounds() (float64, float64) {
	return text.Measure(c.name, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   smallFontSize,
	}, smallFontSize)
}

func (c *Character) BoundsID() (float64, float64) {
	return text.Measure("◇"+c.id[:3], &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   logFontSize,
	}, logFontSize)
}

func (c *Character) Draw(screen *ebiten.Image) {
	opi := &ebiten.DrawImageOptions{}
	opi.Filter = ebiten.FilterLinear
	if c.dir == entity.DirLeft {
		opi.GeoM.Scale(-1, 1)
		opi.GeoM.Translate(characterWidth, 0)
	}
	opi.GeoM.Translate(float64(c.x)-characterWidth/2, float64(c.y)-characterHeight/2)
	screen.DrawImage(c.image(), opi)

	opt := &text.DrawOptions{}
	w, h := c.Bounds()
	opt.GeoM.Translate(float64(float64(c.x)-w/2), float64(c.y)+characterHeight/2)
	opt.ColorScale.ScaleWithColor(color.Black)
	opt.LineSpacing = smallFontSize
	opt.Filter = ebiten.FilterLinear
	text.Draw(screen, c.name, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   smallFontSize,
	}, opt)
	opid := &text.DrawOptions{}
	opid.ColorScale.ScaleWithColor(color.Black)
	opid.LineSpacing = logFontSize
	opid.Filter = ebiten.FilterLinear
	wi, _ := c.BoundsID()
	opid.GeoM.Translate(float64(float64(c.x)-wi/2), float64(c.y)+characterHeight/2+h)
	text.Draw(screen, "◇"+c.id[:3], &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   logFontSize,
	}, opid)
}
