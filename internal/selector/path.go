package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type PathSelect struct {
	*AbstractSelector
	MaxRange   int
	PathChecks floor.PathChecks
	Effect     *AbstractSelectionEffect
}

func (s *PathSelect) SetValues(values ActionValues) {
	s.MaxRange = values.Move
	s.PathChecks.Orig = s.origin
}

func (s *PathSelect) Update(input *input.Input) {
	if !s.isDone {
		//x, y := input.Coords.X, input.Coords.Y
		s.area = []world.Coords{}
		s.PathChecks.Orig = s.origin
		hex := floor.CurrentFloor.IsLegal(input.Coords, s.PathChecks)
		legal := hex != nil && world.DistanceSimple(s.origin, input.Coords) <= s.MaxRange
		if legal {
			path, dist, found := floor.CurrentFloor.FindPath(s.origin, input.Coords, s.PathChecks)
			if found && dist <= s.MaxRange {
				for _, h := range path {
					s.area = append(s.area, h)
				}
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
		}
		if input.Cancel.JustPressed() {
			input.Cancel.Consume()
			s.Cancel()
		}
	}
}

func (s *PathSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}