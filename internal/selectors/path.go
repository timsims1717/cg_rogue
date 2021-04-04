package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type PathSelect struct {
	*AbstractSelector
	MaxRange   int
	PathChecks floor.PathChecks
}

func NewPathSelect(isMove bool, checks floor.PathChecks) *AbstractSelector {
	sel := &AbstractSelector{
		IsMove: isMove,
	}
	target := &PathSelect{
		sel,
		0,
		checks,
	}
	sel.Selector = target
	return sel
}

func (s *PathSelect) SetValues(values ActionValues) {
	s.MaxRange = values.Move
	s.PathChecks.Orig = s.origin
}

func (s *PathSelect) Update(input *input.Input) {
	if !s.isDone {
		x, y := input.Coords.X, input.Coords.Y
		s.PathChecks.Orig = s.origin
		hex := floor.CurrentFloor.IsLegal(input.Coords, s.PathChecks)
		legal := hex != nil && world.DistanceSimple(s.origin, input.Coords) <= s.MaxRange
		if legal {
			path, dist, found := floor.CurrentFloor.FindPath(s.origin, input.Coords, s.PathChecks)
			if found && dist <= s.MaxRange {
				if input.Select.JustPressed() {
					input.Select.Consume()
					for _, h := range path {
						if h.X != s.origin.X || h.Y != s.origin.Y {
							s.area = append(s.area, h)
						}
					}
					s.isDone = true
				}
				for _, h := range path {
					if h.X == x && h.Y == y {
						AddSelectUI(MoveSolid, h.X, h.Y)
					} else {
						AddSelectUI(Move, h.X, h.Y)
					}
				}
			}
		}
	}
}