package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HexSelect struct {
	*AbstractSelector
	Count    int
	MaxRange int
}

func NewHexSelect(isMove bool) *AbstractSelector {
	sel := &AbstractSelector{
		IsMove: isMove,
	}
	target := &HexSelect{
		sel,
		0,
		0,
	}
	sel.Selector = target
	return sel
}

func (s *HexSelect) SetValues(values ActionValues) {
	s.Count = values.Targets
	s.MaxRange = values.Range
}

func (s *HexSelect) Update(input *input.Input) {
	if !s.isDone {
		x := input.Coords.X
		y := input.Coords.Y
		checks := floor.PathChecks{
			NotFilled:  true,
			Unoccupied: false,
			NonEmpty:   false,
			Orig:       s.origin,
		}
		hex := floor.CurrentFloor.IsLegal(input.Coords, checks)
		legal := hex != nil && world.DistanceSimple(s.origin, input.Coords) <= s.MaxRange
		if legal {
			if input.Select.JustPressed() {
				input.Select.Consume()
				// add to or remove from the clicked array
				found := -1
				for i, hex := range s.area {
					if hex.X == x && hex.Y == y {
						found = i
						break
					}
				}
				if found == -1 {
					// add to clicked array
					s.area = append(s.area, input.Coords)
				} else {
					// remove from clicked array
					s.area[len(s.area)-1], s.area[found] = s.area[found], s.area[len(s.area)-1]
					s.area = s.area[:len(s.area)-1]
				}
			}
		}
		for _, sel := range s.area {
			if s.IsMove {
				AddSelectUI(Move, sel.X, sel.Y)
			} else {
				AddSelectUI(Attack, sel.X, sel.Y)
			}
		}
		if legal {
			if s.IsMove {
				AddSelectUI(MoveSolid, x, y)
			} else {
				AddSelectUI(Attack, x, y)
			}
		} else {
			AddSelectUI(Blank, x, y)
		}
	}
	if input.Cancel.JustPressed() {
		input.Cancel.Consume()
		s.area = []world.Coords{}
	}
	if len(s.area) >= s.Count {
		s.isDone = true
	}
}

type MoveHexSelect struct {
	*AbstractSelector
	Count    int
	MaxRange int
}

func NewMoveHexSelect() *AbstractSelector {
	sel := &AbstractSelector{
		IsMove: true,
	}
	target := &MoveHexSelect{
		sel,
		0,
		0,
	}
	sel.Selector = target
	return sel
}

func (s *MoveHexSelect) SetValues(values ActionValues) {
	s.Count = values.Targets
	s.MaxRange = values.Move
}

func (s *MoveHexSelect) Update(input *input.Input) {
	if !s.isDone {
		x := input.Coords.X
		y := input.Coords.Y
		hex := floor.CurrentFloor.Get(input.Coords)
		legal := hex != nil && hex.Occupant == nil && !hex.Empty && world.DistanceSimple(s.origin, input.Coords) <= s.MaxRange
		if legal {
			if input.Select.JustPressed() {
				input.Select.Consume()
				// add to or remove from the clicked array
				found := -1
				for i, hex := range s.area {
					if hex.X == x && hex.Y == y {
						found = i
						break
					}
				}
				if found == -1 {
					// add to clicked array
					s.area = append(s.area, input.Coords)
				} else {
					// remove from clicked array
					s.area[len(s.area)-1], s.area[found] = s.area[found], s.area[len(s.area)-1]
					s.area = s.area[:len(s.area)-1]
				}
			}
		}
		for _, sel := range s.area {
			AddSelectUI(Move, sel.X, sel.Y)
		}
		if legal {
			AddSelectUI(MoveSolid, x, y)
		} else {
			AddSelectUI(Blank, x, y)
		}
	}
	if input.Cancel.JustPressed() {
		input.Cancel.Consume()
		s.area = []world.Coords{}
	}
	if len(s.area) >= s.Count {
		s.isDone = true
	}
}