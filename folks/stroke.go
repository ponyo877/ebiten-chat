package folks

type Stroke struct {
	source StrokeSource

	// initX and initY represents the position when dragging starts.
	initX int
	initY int

	// x and y represents the current position
	x int
	y int

	released bool

	// draggingObject represents a object (sprite in this case)
	// that is being dragged.
	draggingObject any
}

func NewStroke(source StrokeSource) *Stroke {
	cx, cy := source.Position()
	return &Stroke{
		source: source,
		initX:  cx,
		initY:  cy,
		x:      cx,
		y:      cy,
	}
}

func (s *Stroke) Update() {
	if s.released {
		return
	}
	if s.source.IsJustReleased() {
		s.released = true
		return
	}
	s.x, s.y = s.source.Position()
}

func (s *Stroke) IsReleased() bool {
	return s.released
}

func (s *Stroke) Position() (int, int) {
	return s.x, s.y
}

func (s *Stroke) PositionDiff() (int, int) {
	dx := s.x - s.initX
	dy := s.y - s.initY
	return dx, dy
}

func (s *Stroke) DraggingObject() any {
	return s.draggingObject
}

func (s *Stroke) SetDraggingObject(object any) {
	s.draggingObject = object
}
