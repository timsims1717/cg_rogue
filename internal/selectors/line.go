package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type LineSelect struct {
	input      *input.Input
	clicked    []world.Coords
	count      int
	maxRange   int
	origin     world.Coords
	pathChecks floor.PathChecks
	isDone     bool
	isAtk      bool
}

func NewLineSelect() *LineSelect {
	return &LineSelect{}
}

func (s *LineSelect) Init(input *input.Input) {
	s.isDone = false
	s.input = input
	s.clicked = []world.Coords{}
}

func (s *LineSelect) SetValues(values ActionValues) {
	s.origin = values.Source.Coords
	s.maxRange = values.Range
	s.count = values.Targets
	s.isAtk = values.Damage > 0
	s.pathChecks = values.Checks
}

func (s *LineSelect) Update() {
	if !s.isDone {
		s.pathChecks.Orig = s.origin
		path := floor.CurrentFloor.LongestLegalPath(floor.CurrentFloor.Line(s.origin, s.input.Coords, s.maxRange), 0, s.pathChecks)
		targets := make([]world.Coords, 0)
		if s.count == 0 {
			targets = path
		} else {
			for _, p := range path {
				if p.Equals(s.origin) {
					continue
				}
				if len(targets) >= s.count {
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
			if s.input.Select.JustPressed() {
				s.input.Select.Consume()
				// add to or remove from the clicked array
				s.clicked = targets
			}
			i := 0
			for _, p := range path {
				if i < len(targets) {
					sel := targets[i]
					if sel.Equals(p) {
						if s.isAtk {
							AddSelectUI(Attack, sel.X, sel.Y)
						} else {
							AddSelectUI(Move, sel.X, sel.Y)
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

func (s *LineSelect) IsDone() bool {
	return len(s.clicked) > 0
}

func (s *LineSelect) Finish() []world.Coords {
	return s.clicked
}