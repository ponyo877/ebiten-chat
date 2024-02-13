package websocket

type Signal struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func NewSignal(event string, data []byte) *Signal {
	return &Signal{
		Event: event,
		Data:  string(data),
	}
}
