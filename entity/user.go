package entity

type User struct {
	id  string
	x   int
	y   int
	dir Dir
}

func NewUser(id string, x, y int, dir Dir) *User {
	return &User{
		id:  id,
		x:   x,
		y:   y,
		dir: dir,
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) X() int {
	return u.x
}

func (u *User) Y() int {
	return u.y
}

func (u *User) Dir() Dir {
	return u.dir
}
