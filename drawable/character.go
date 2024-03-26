package drawable

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ponyo877/folks-ui/entity"
)

type Character struct {
	id    string
	image *ebiten.Image
	x     int
	y     int
	dir   entity.Dir
}

func NewCharacter(id string, image *ebiten.Image, x, y int, dir entity.Dir) *Character {
	return &Character{
		id,
		image,
		x,
		y,
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

func (c *Character) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	if c.dir == entity.DirLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(characterWidth, 0)
	}
	op.GeoM.Translate(float64(c.x)-characterWidth/2, float64(c.y)-characterHeight/2)
	screen.DrawImage(c.image, op)
}
