package presenter

import (
	"github.com/ponyo877/folks-ui/entity"
)

type BodyPresenter struct {
	ID    string           `json:"id"`
	X     int              `json:"x,omitempty"`
	Y     int              `json:"y,omitempty"`
	Name  string           `json:"name,omitempty"`
	ImgID int              `json:"imgid,omitempty"`
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
		ID:    body.ID(),
		X:     body.X(),
		Y:     body.Y(),
		Name:  body.Name(),
		ImgID: body.ImgID(),
		Dir:   int(body.Dir()),
	}
}

func NewEnterRespBodyPresenter(body *entity.Body) *BodyPresenter {
	return &BodyPresenter{
		Users: NewUsersPresenter(body.Users()),
	}
}

func NewLeaveBodyPresenter(body *entity.Body) *BodyPresenter {
	return &BodyPresenter{
		ID: body.ID(),
	}
}

func NewEnterReqBodyPresenter(body *entity.Body) *BodyPresenter {
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
	case "enter":
		return NewEnterReqBodyPresenter(body)
	}
	return nil
}

func (b *BodyPresenter) Unmarshal(messageType string) *entity.Body {
	switch messageType {
	case "say":
		return entity.NewSayBody(b.ID, b.Name, b.Text)
	case "move":
		return entity.NewMoveBody(b.ID, b.X, b.Y, b.Name, b.ImgID, entity.Dir(b.Dir))
	case "enter":
		var users []*entity.User
		for _, up := range b.Users {
			user := entity.NewUser(up.ID, up.X, up.Y, up.Name, up.ImgID, entity.Dir(up.Dir))
			users = append(users, user)
		}
		return entity.NewEnterRespBody(users)
	case "leave":
		return entity.NewLeaveBody(b.ID)
	}
	return nil
}
