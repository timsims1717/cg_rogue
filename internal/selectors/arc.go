package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type ArcSelect struct {
	*AbstractSelector
	Count      int
	MaxRange   int
	PathChecks floor.PathChecks
}

func NewArcSelect() *AbstractSelector {
	sel := &AbstractSelector{}
	target := &ArcSelect{
		sel,
		0,
		0,
		floor.PathChecks{},
	}
	sel.Selector = target
	return sel
}

func (s *ArcSelect) SetValues(values ActionValues) {
	s.Count = values.Targets
	s.MaxRange = values.Range
}

func (s *ArcSelect) Update(input *input.Input) {
	if !s.isDone {
		var neighbors []world.Coords
		dist := world.DistanceSimple(input.Coords, s.origin)
		if dist < 2 || dist % 2 == 0 {
			neighbors = world.OrderByDistWorld(input.World, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		} else {
			neighbors = world.OrderByDist(input.Coords, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		}
		var closest []world.Coords
		for i, n := range neighbors {
			if i < s.Count {
				hex := floor.CurrentFloor.IsLegal(n, s.PathChecks)
				legal := hex != nil && world.DistanceSimple(s.origin, n) <= s.MaxRange
				if legal {
					closest = append(closest, n)
				}
			} else {
				break
			}
		}
		if len(closest) > 0 {
			if input.Select.JustPressed() {
				input.Select.Consume()
				// add to or remove from the clicked array
				s.area = closest
				s.isDone = true
			}
			for _, sel := range closest {
				if s.IsMove {
					AddSelectUI(Move, sel.X, sel.Y)
				} else {
					AddSelectUI(Attack, sel.X, sel.Y)
				}
			}
		}
	}
}