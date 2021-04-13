package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type TargetSelect struct {
	*AbstractSelector
	count    int
	maxRange int
	Effect   *AbstractSelectionEffect
}

func (s *TargetSelect) SetValues(values ActionValues) {
	s.count = values.Targets
	s.maxRange = values.Range
}

func (s *TargetSelect) Update(input *input.Input) {
	if !s.isDone {
		inRange := world.DistanceSimple(s.origin, input.Coords) <= s.maxRange && input.Coords != s.origin
		occ := floor.CurrentFloor.GetOccupant(input.Coords)
		if occ != nil && inRange {
			if input.Select.JustPressed() {
				input.Select.Consume()
				// add to or remove from area
				var removed int
				s.area, removed = world.Remove(input.Coords, s.area)
				if removed == -1 {
					s.area = append(s.area, input.Coords)
				}
			}
		}

		if s.Effect != nil {
			if occ != nil && inRange {
				s.Effect.SetArea(append(s.area, input.Coords))
			} else {
				s.Effect.SetArea(s.area)
				eff := NewSelectionEffect(&HighlightEffect{})
				eff.SetArea([]world.Coords{input.Coords})
				AddSelectionEffect(eff)
			}
			AddSelectionEffect(s.Effect)
		}

		if input.Cancel.JustPressed() {
			input.Cancel.Consume()
			s.Cancel()
		}
		if len(s.area) >= s.count {
			s.isDone = true
			s.results = []*Result{
				NewResult(s.area, s.Effect, false),
			}
		}
	}
}

func (s *TargetSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}
