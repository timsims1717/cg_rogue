package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type ArcSelect struct {
	input      *input.Input
	clicked    []world.Coords
	count      int
	maxRange   int
	origin     world.Coords
	pathChecks floor.PathChecks
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
	s.pathChecks = values.Checks
}

func (s *ArcSelect) Update() {
	if !s.isDone {
		var neighbors []world.Coords
		dist := world.DistanceSimple(s.input.Coords, s.origin)
		if dist < 2 || dist % 2 == 0 {
			neighbors = world.OrderByDistWorld(s.input.World, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		} else {
			neighbors = world.OrderByDist(s.input.Coords, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		}
		var closest []world.Coords
		for i, n := range neighbors {
			if i < s.count {
				hex := floor.CurrentFloor.IsLegal(n, s.pathChecks)
				legal := hex != nil && world.DistanceSimple(s.origin, n) <= s.maxRange
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
					AddSelectUI(Attack, sel.X, sel.Y)
				} else {
					AddSelectUI(Move, sel.X, sel.Y)
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