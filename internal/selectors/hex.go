package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HexSelect struct {
	input    *input.Input
	clicked  []world.Coords
	count    int
	maxRange int
	origin   world.Coords
	isDone   bool
	isAtk    bool
}

func NewHexSelect() *HexSelect {
	return &HexSelect{}
}

func (s *HexSelect) Init(input *input.Input) {
	s.isDone = false
	s.input = input
	s.clicked = []world.Coords{}
}

func (s *HexSelect) SetValues(values ActionValues) {
	s.origin = values.Source.Coords
	s.maxRange = util.Max(values.Range, values.Move)
	s.count = values.Targets
	s.isAtk = values.Damage > 0
}

func (s *HexSelect) Update() {
	if !s.isDone {
		x := s.input.Coords.X
		y := s.input.Coords.Y
		checks := floor.PathChecks{
			NotFilled:  true,
			Unoccupied: false,
			NonEmpty:   false,
			Orig:       s.origin,
		}
		hex := floor.CurrentFloor.IsLegal(s.input.Coords, checks)
		legal := hex != nil && world.DistanceSimple(s.origin, s.input.Coords) <= s.maxRange
		if legal {
			if s.input.Select.JustPressed() {
				s.input.Select.Consume()
				// add to or remove from the clicked array
				found := -1
				for i, hex := range s.clicked {
					if hex.X == x && hex.Y == y {
						found = i
						break
					}
				}
				if found == -1 {
					// add to clicked array
					s.clicked = append(s.clicked, s.input.Coords)
				} else {
					// remove from clicked array
					s.clicked[len(s.clicked)-1], s.clicked[found] = s.clicked[found], s.clicked[len(s.clicked)-1]
					s.clicked = s.clicked[:len(s.clicked)-1]
				}
			}
		}
		for _, sel := range s.clicked {
			if s.isAtk {
				AddSelectUI(Attack, sel.X, sel.Y)
			} else {
				AddSelectUI(Move, sel.X, sel.Y)
			}
		}
		if legal {
			if s.isAtk {
				AddSelectUI(Attack, x, y)
			} else {
				AddSelectUI(MoveSolid, x, y)
			}
		} else {
			AddSelectUI(Blank, x, y)
		}
	}
	if s.input.Cancel.JustPressed() {
		s.input.Cancel.Consume()
		s.clicked = []world.Coords{}
	}
}

func (s *HexSelect) IsDone() bool {
	return len(s.clicked) == s.count || s.isDone
}

func (s *HexSelect) Finish() []world.Coords {
	return s.clicked
}

type EmptyHexSelect struct {
	input    *input.Input
	clicked  []world.Coords
	count    int
	maxRange int
	origin   world.Coords
	isDone   bool
}

func NewEmptyHexSelect() *EmptyHexSelect {
	return &EmptyHexSelect{}
}

func (s *EmptyHexSelect) Init(input *input.Input) {
	s.isDone = false
	s.input = input
	s.clicked = []world.Coords{}
}

func (s *EmptyHexSelect) SetValues(values ActionValues) {
	s.origin = values.Source.Coords
	s.maxRange = util.Max(values.Range, values.Move)
	s.count = values.Targets
}

func (s *EmptyHexSelect) Update() {
	if !s.isDone {
		x := s.input.Coords.X
		y := s.input.Coords.Y
		hex := floor.CurrentFloor.Get(s.input.Coords)
		legal := hex != nil && hex.Occupant == nil && !hex.Empty && world.DistanceSimple(s.origin, s.input.Coords) <= s.maxRange
		if legal {
			if s.input.Select.JustPressed() {
				s.input.Select.Consume()
				// add to or remove from the clicked array
				found := -1
				for i, hex := range s.clicked {
					if hex.X == x && hex.Y == y {
						found = i
						break
					}
				}
				if found == -1 {
					// add to clicked array
					s.clicked = append(s.clicked, s.input.Coords)
				} else {
					// remove from clicked array
					s.clicked[len(s.clicked)-1], s.clicked[found] = s.clicked[found], s.clicked[len(s.clicked)-1]
					s.clicked = s.clicked[:len(s.clicked)-1]
				}
			}
		}
		for _, sel := range s.clicked {
			AddSelectUI(Move, sel.X, sel.Y)
		}
		if legal {
			AddSelectUI(MoveSolid, x, y)
		} else {
			AddSelectUI(Blank, x, y)
		}
	}
	if s.input.Cancel.JustPressed() {
		s.input.Cancel.Consume()
		s.clicked = []world.Coords{}
	}
}

func (s *EmptyHexSelect) IsDone() bool {
	return len(s.clicked) == s.count || s.isDone
}

func (s *EmptyHexSelect) Finish() []world.Coords {
	return s.clicked
}