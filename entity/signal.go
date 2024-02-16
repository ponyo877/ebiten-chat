package entity

type Signal struct {
	event string
	data  string
}

func NewSignal(event string, data string) *Signal {
	return &Signal{
		event: event,
		data:  data,
	}
}

func (s *Signal) Event() string {
	return s.event
}

func (s *Signal) Data() []byte {
	return []byte(s.data)
}
