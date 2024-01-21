package websocket

import (
	"github.com/ponyo877/folks-ui/entity"
)

type BodyPresenter struct {
	ID    string           `json:"id"`
	X     int              `json:"x,omitempty"`
	Y     int              `json:"y,omitempty"`
	Dir   int              `json:"dir,omitempty"`
	Text  string           `json:"text,omitempty"`
	Users []*UserPresenter `json:"users,omitempty"`
}

func NewSayBodyPresenter(body *entity.Body) *BodyPresenter {
	return &BodyPresenter{
		ID:   body.ID(),
		Text: body.Text(),
	}
}

func NewMoveBodyPresenter(body *entity.Body) *BodyPresenter {
	return &BodyPresenter{
		ID:  body.ID(),
		X:   body.X(),
		Y:   body.Y(),
		Dir: int(body.Dir()),
	}
}

func NewEnterBodyPresenter(body *entity.Body) *BodyPresenter {
	return &BodyPresenter{
		Users: NewUsersPresenter(body.Users()),
	}
}

func NewLeaveBodyPresenter(body *entity.Body) *BodyPresenter {
	return &BodyPresenter{
		ID: body.ID(),
	}
}

func NewBodyPresenter(messageType string, body *entity.Body) *BodyPresenter {
	switch messageType {
	case "say":
		return NewSayBodyPresenter(body)
	case "move":
		return NewMoveBodyPresenter(body)
	case "leave":
		return NewLeaveBodyPresenter(body)
	}
	return nil
}

func (b *BodyPresenter) Unmarshal(messageType string) *entity.Body {
	switch messageType {
	case "say":
		return entity.NewSayBody(b.ID, b.Text)
	case "move":
		return entity.NewMoveBody(b.ID, b.X, b.Y, entity.Dir(b.Dir))
	case "enter":
		var users []*entity.User
		for _, userPresenter := range b.Users {
			user := entity.NewUser(userPresenter.ID, userPresenter.X, userPresenter.Y, entity.Dir(userPresenter.Dir))
			users = append(users, user)
		}
		return entity.NewEnterBody(users)
	case "leave":
		return entity.NewLeaveBody(b.ID)
	}
	return nil
}

func MarshalBody(messageType string, body *entity.Body) *BodyPresenter {
	return NewBodyPresenter(messageType, body)
}
