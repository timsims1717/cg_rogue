package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type LineTargetSelect struct {
	*AbstractSelector
	Count      int
	MaxRange   int
	PathChecks floor.PathChecks
	Effect     *AbstractSelectionEffect
	SecEffect  *AbstractSelectionEffect
}

func (s *LineTargetSelect) SetValues(values ActionValues) {
	s.Count = values.Targets
	if s.IsMove {
		s.MaxRange = values.Move
	} else {
		s.MaxRange = values.Range
	}
	s.PathChecks.Orig = s.origin
	s.Effect.SetValues(values)
	s.Effect.SetOrig(s.origin)
	s.SecEffect.SetValues(values)
	s.SecEffect.SetOrig(s.origin)
}

func (s *LineTargetSelect) Update(input *input.Input) {
	if !s.isDone {
		s.PathChecks.Orig = s.origin
		path := floor.CurrentFloor.LongestLegalPath(floor.CurrentFloor.Line(s.origin, input.Coords, s.MaxRange), 0, s.PathChecks)
		targets := make([]world.Coords, 0)
		if s.Count == 0 {
			targets = path
		} else {
			for _, p := range path {
				if p.Eq(s.origin) {
					continue
				}
				if len(targets) >= s.Count {
					break
				}
				if occ := floor.CurrentFloor.GetOccupant(p); occ != nil {
					targets = append(targets, p)
				}
			}
		}
		if len(path) > 0 {
			i := 0
			var primary []world.Coords
			var secondary []world.Coords
			for _, p := range path {
				if i < len(targets) {
					sel := targets[i]
					if sel.Eq(p) {
						primary = append(primary, p)
						i++
						continue
					}
				}
				secondary = append(secondary, p)
			}
			if len(primary) > 0 && s.Effect != nil {
				s.Effect.SetArea(primary)
				AddSelectionEffect(s.Effect)
			}
			if len(secondary) > 0 && s.SecEffect != nil {
				s.SecEffect.SetArea(secondary)
				AddSelectionEffect(s.SecEffect)
			}
			if input.Select.JustPressed() && len(targets) > 0 {
				input.Select.Consume()
				s.isDone = true
				s.results = []*Result{
					NewResult(primary, s.Effect, s.IsMove),
					NewResult(secondary, s.SecEffect, false),
				}
			}
		}
		if input.Cancel.JustPressed() {
			input.Cancel.Consume()
			s.Cancel()
		}
	}
}

func (s *LineTargetSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}
