package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type ArcSelect struct {
	*AbstractSelector
	count      int
	maxRange   int
	PathChecks floor.PathChecks
	Effect     *AbstractSelectionEffect
}

func (s *ArcSelect) SetValues(values ActionValues) {
	s.count = values.Targets
	s.maxRange = values.Range
	s.PathChecks.Orig = s.origin
}

func (s *ArcSelect) Update(input *input.Input) {
	if !s.isDone {
		var neighbors []world.Coords
		dist := world.DistanceSimple(input.Coords, s.origin)
		if dist < 2 || dist%2 == 0 {
			neighbors = world.OrderByDistWorld(input.World, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		} else {
			neighbors = world.OrderByDist(input.Coords, s.origin.Neighbors(floor.CurrentFloor.Dimensions()))
		}
		var closest []world.Coords
		for i, n := range neighbors {
			if i < s.count {
				hex := floor.CurrentFloor.IsLegal(n, s.PathChecks)
				legal := hex != nil && world.DistanceSimple(s.origin, n) <= s.maxRange
				if legal {
					closest = append(closest, n)
				}
			} else {
				break
			}
		}
		if len(closest) > 0 {
			s.area = closest
			if s.Effect != nil {
				s.Effect.SetArea(s.area)
				AddSelectionEffect(s.Effect)
			}
			if input.Select.JustPressed() {
				input.Select.Consume()
				s.isDone = true
				s.results = []*Result{
					NewResult(s.area, s.Effect, s.IsMove),
				}
			}
		}
		if input.Cancel.JustPressed() {
			input.Cancel.Consume()
			s.Cancel()
		}
	}
}

func (s *ArcSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}
