package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type ArcSelect struct {
	input      *input.Input
	clicked    []world.Coords
	count      int
	maxRange   int
	origin     world.Coords
	isDone     bool
	isAtk      bool
}

func NewArcSelect() *ArcSelect {
	return &ArcSelect{}
}

func (s *ArcSelect) Init(input *input.Input) {
	s.isDone = false
	s.input = input
	s.clicked = []world.Coords{}
}

func (s *ArcSelect) SetValues(values ActionValues) {
	s.origin = values.Source.Coords
	s.maxRange = values.Range
	s.count = values.Targets
	s.isAtk = values.Damage > 0
}

func (s *ArcSelect) Update() {
	if !s.isDone {
		//x := s.input.Coords.X
		//y := s.input.Coords.Y
		checks := floor.PathChecks{
			NotFilled:  true,
			Unoccupied: false,
			NonEmpty:   false,
			Orig:       s.origin,
		}
		var neighbors []world.Coords
		dist := world.DistanceSimpleHex(s.input.Coords, s.origin)
		if dist < 2 || dist % 2 == 0 {
			neighbors = world.OrderByDistWorld(s.input.World, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		} else {
			neighbors = world.OrderByDist(s.input.Coords, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		}
		var closest []world.Coords
		for i, n := range neighbors {
			if i < s.count {
				hex := floor.CurrentFloor.IsLegal(n, checks)
				legal := hex != nil && world.DistanceSimpleHex(s.origin, n) <= s.maxRange
				if legal {
					closest = append(closest, n)
				}
			} else {
				break
			}
		}
		if len(closest) > 0 {
			if s.input.Select.JustPressed() {
				s.input.Select.Consume()
				// add to or remove from the clicked array
				s.clicked = closest
			}
			for _, sel := range closest {
				if s.isAtk {
					ui.AddSelectUI(ui.Attack, sel.X, sel.Y)
				} else {
					ui.AddSelectUI(ui.Move, sel.X, sel.Y)
				}
			}
		}
	}
}

func (s *ArcSelect) IsDone() bool {
	return len(s.clicked) > 0
}

func (s *ArcSelect) Finish() []world.Coords {
	return s.clicked
}