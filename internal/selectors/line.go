package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type LineSelect struct {
	*AbstractSelector
	Count      int
	MaxRange   int
	PathChecks floor.PathChecks
}

func NewLineSelect(isMove bool) *AbstractSelector {
	sel := &AbstractSelector{
		IsMove: isMove,
	}
	target := &LineSelect{
		sel,
		0,
		0,
		floor.PathChecks{},
	}
	sel.Selector = target
	return sel
}

func (s *LineSelect) SetValues(values ActionValues) {
	s.Count = values.Targets
	if s.IsMove {
		s.MaxRange = values.Move
	} else {
		s.MaxRange = values.Range
	}
	s.PathChecks = values.Checks
}

func (s *LineSelect) Update(input *input.Input) {
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
					if _, ok := occ.(objects.Targetable); ok {
						targets = append(targets, p)
					}
				}
			}
		}
		if len(path) > 0 {
			if input.Select.JustPressed() {
				input.Select.Consume()
				// add to or remove from the clicked array
				s.area = targets
				s.isDone = true
			}
			i := 0
			for _, p := range path {
				if i < len(targets) {
					sel := targets[i]
					if sel.Eq(p) {
						if s.IsMove {
							AddSelectUI(Move, sel.X, sel.Y)
						} else {
							AddSelectUI(Attack, sel.X, sel.Y)
						}
						i++
						continue
					}
				}
				AddSelectUI(Default, p.X, p.Y)
			}
		}
	}
}