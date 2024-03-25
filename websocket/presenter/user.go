package presenter

import "github.com/ponyo877/folks-ui/entity"

type UserPresenter struct {
	ID  string `json:"id"`
	X   int    `json:"x"`
	Y   int    `json:"y"`
	Dir int    `json:"dir"`
}

func NewUserPresenter(user *entity.User) *UserPresenter {
	return &UserPresenter{
		ID:  user.ID(),
		X:   user.X(),
		Y:   user.Y(),
		Dir: int(user.Dir()),
	}
}

func NewUsersPresenter(users []*entity.User) []*UserPresenter {
	var usersPresenter []*UserPresenter
	for _, user := range users {
		userPresenter := NewUserPresenter(user)
		usersPresenter = append(usersPresenter, userPresenter)
	}
	return usersPresenter
}
