package websocket

import (
	"encoding/json"

	"github.com/ponyo877/folks-ui/entity"
)

type SignalPresenter struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func NewSignalPresenter(event string, data []byte) (*SignalPresenter, error) {
	v, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &SignalPresenter{
		Event: event,
		Data:  string(v),
	}, nil
}

func UnmarshalSignal(data []byte) (*entity.Signal, error) {
	var s SignalPresenter
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return entity.NewSignal(s.Event, s.Data), nil
}
