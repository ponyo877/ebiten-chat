package entity

type Body struct {
	id   string
	x    int
	y    int
	dir  Dir
	text string
}

func NewSayBody(id string, text string) *Body {
	return &Body{
		id:   id,
		text: text,
	}
}

func NewMoveBody(id string, x, y int, dir Dir) *Body {
	return &Body{
		id:  id,
		x:   x,
		y:   y,
		dir: dir,
	}
}

func NewLeaveBody(id string) *Body {
	return &Body{
		id: id,
	}
}

func (b *Body) ID() string {
	return b.id
}

func (b *Body) X() int {
	return b.x
}

func (b *Body) Y() int {
	return b.y
}

func (b *Body) Dir() Dir {
	return b.dir
}

func (b *Body) Text() string {
	return b.text
}
