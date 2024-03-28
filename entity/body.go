package entity

import "fmt"

type Body struct {
	id    string
	x     int
	y     int
	name  string
	imgid int
	dir   Dir
	text  string
	users []*User
}

func NewSayBody(id, name, text string) *Body {
	return &Body{
		id:   id,
		name: name,
		text: text,
	}
}

func NewMoveBody(id string, x, y int, name string, imgid int, dir Dir) *Body {
	if imgid > 20 {
		panic(fmt.Sprintf("imgid must be less than 20, but got %d", imgid))
	}
	return &Body{
		id:    id,
		x:     x,
		y:     y,
		name:  name,
		imgid: imgid,
		dir:   dir,
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

func (b *Body) X() int {
	return b.x
}

func (b *Body) Y() int {
	return b.y
}

func (b *Body) Name() string {
	return b.name
}

func (b *Body) ImgID() int {
	return b.imgid
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
