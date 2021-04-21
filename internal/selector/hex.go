package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HexSelect struct {
	*AbstractSelector
	count      int
	maxRange   int
	PathChecks floor.PathChecks
	Effect     *AbstractSelectionEffect
}

func (s *HexSelect) SetValues(values ActionValues) {
	s.count = values.Targets
	if s.IsMove {
		s.maxRange = values.Move
	} else {
		s.maxRange = values.Range
	}
	if s.Effect != nil {
		s.Effect.SetValues(values)
		s.Effect.SetOrig(s.origin)
	}
}

func (s *HexSelect) Update(input *input.Input) {
	if !s.isDone {
		if s.IsMove {
			s.source.RemoveClaim()
		}
		hex := floor.CurrentFloor.IsLegal(input.Coords, s.PathChecks)
		legal := hex != nil && world.DistanceSimple(s.origin, input.Coords) <= s.maxRange
		if legal {
			if input.Select.JustPressed() {
				input.Select.Consume()
				// add to or remove from area
				var removed int
				s.area, removed = world.Remove(input.Coords, s.area)
				if removed == -1 {
					s.area = append(s.area, input.Coords)
				}
			}
			if s.IsMove {
				s.source.MakeClaim(input.Coords)
			}
		}
		if s.Effect != nil {
			if legal && !world.CoordsIn(input.Coords, s.area) {
				s.Effect.SetArea(append(s.area, input.Coords))
			} else {
				s.Effect.SetArea(s.area)
				eff := NewSelectionEffect(&HighlightEffect{}, s.Effect.values)
				eff.SetArea([]world.Coords{input.Coords})
				AddSelectionEffect(eff)
			}
			AddSelectionEffect(s.Effect)
		}
		if input.Cancel.JustPressed() {
			input.Cancel.Consume()
			var removed int
			s.area, removed = world.Remove(input.Coords, s.area)
			if removed == -1 {
				s.Cancel()
			}
		}
		if len(s.area) >= s.count {
			s.isDone = true
			s.results = []*Result{
				NewResult(s.area, s.Effect, s.IsMove),
			}
		}
	}
}

func (s *HexSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}
