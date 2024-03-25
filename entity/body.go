package entity

import (
	"hash/fnv"
)

type Body struct {
	id    string
	x     int
	y     int
	dir   Dir
	text  string
	users []*User
}

func NewSayBody(id, text string) *Body {
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

func NewEnterReqBody(id string) *Body {
	return &Body{
		id: id,
	}
}

func NewEnterRespBody(users []*User) *Body {
	return &Body{
		users: users,
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

func (b *Body) ImgIdx() int {
	h := fnv.New32a()
	h.Write([]byte(b.id))
	return int(h.Sum32()) % 21
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

func (b *Body) Users() []*User {
	return b.users
}
