package websocket

import (
	"encoding/json"

	"github.com/ponyo877/folks-ui/entity"
)

type SignalPresenter struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func NewSignalPresenter(event string, data string) *SignalPresenter {
	return &SignalPresenter{
		Event: event,
		Data:  data,
	}
}

func UnmarshalSignal(data []byte) (*entity.Signal, error) {
	var s SignalPresenter
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return entity.NewSignal(s.Event, s.Data), nil
}
