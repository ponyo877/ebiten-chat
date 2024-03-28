package entity

type User struct {
	id    string
	x     int
	y     int
	name  string
	imgid int
	dir   Dir
}

func NewUser(id string, x, y int, name string, imgid int, dir Dir) *User {
	return &User{
		id:    id,
		x:     x,
		y:     y,
		name:  name,
		imgid: imgid,
		dir:   dir,
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

func (u *User) Name() string {
	return u.name
}

func (u *User) ImgID() int {
	return u.imgid
	// h := fnv.New32a()
	// h.Write([]byte(u.id))
	// return int(h.Sum32()) % 20
}

func (u *User) Dir() Dir {
	return u.dir
}
