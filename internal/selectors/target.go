package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type TargetSelect struct {
	input      *input.Input
	clicked    []world.Coords
	Count      int
	MaxRange   int
	origin     world.Coords
	isDone     bool
	cancel     bool
}

func NewTargetSelect() *TargetSelect {
	return &TargetSelect{}
}

func (s *TargetSelect) Init(input *input.Input) {
	s.isDone = false
	s.cancel = false
	s.input = input
	s.clicked = []world.Coords{}
}

func (s *TargetSelect) SetValues(values actions.ActionValues) {
	s.origin = values.Source.Coords
	s.Count = values.Targets
	s.MaxRange = values.Range
}

func (s *TargetSelect) Update() {
	if !s.isDone {
		x := s.input.Coords.X
		y := s.input.Coords.Y
		inRange := world.DistanceHex(s.origin, s.input.Coords) <= s.MaxRange
		occ := floor.CurrentFloor.GetOccupant(s.input.Coords)
		if occ != nil {
			if _, ok := occ.(objects.Targetable); ok {
				if !inRange {
					// todo: highlight
				} else if s.input.Select.JustPressed() {
					s.input.Select.Consume()
					s.clicked = append(s.clicked, s.input.Coords)
					s.isDone = true
				} else {
					ui.AddSelectUI(ui.MoveSolid, x, y)
				}
			}
		} else {
			if inRange {
				ui.AddSelectUI(ui.Default, x, y)
			} else {
				ui.AddSelectUI(ui.Blank, x, y)
			}
		}
	}
	if s.input.Cancel.JustPressed() {
		s.input.Cancel.Consume()
		// cancel
		s.cancel = true
	}
}

func (s *TargetSelect) IsCancelled() bool {
	return s.cancel
}

func (s *TargetSelect) IsDone() bool {
	return len(s.clicked) > 0
}

func (s *TargetSelect) Finish() []world.Coords {
	return s.clicked
}