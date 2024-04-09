package drawable

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

func (c *Character) ShortID() string {
	return fmt.Sprintf("◇%s", c.id[:3])
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) image() *ebiten.Image {
	return CharacterImages[c.imgid]
}

func (c *Character) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	if c.dir == entity.DirLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(characterWidth, 0)
	}
	op.GeoM.Translate(float64(c.x)-characterWidth/2, float64(c.y)-characterHeight/2)
	screen.DrawImage(c.image(), op)
	nameTxt := NewText(float64(c.x), float64(c.y)+characterHeight/2, SmallFontSize, c.name, color.Black, arcadeFaceSource)
	nameTxt.Draw(screen, true)

	_, h := nameTxt.Bounds()
	idTxt := NewText(float64(c.x), float64(c.y)+characterHeight/2+h, LogFontSize, c.ShortID(), color.Black, arcadeFaceSource)
	idTxt.Draw(screen, true)
}
