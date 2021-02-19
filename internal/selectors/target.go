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

func (t *TargetSelect) Init(input *input.Input) {
	t.isDone = false
	t.cancel = false
	t.input = input
	t.clicked = []world.Coords{}
}

func (t *TargetSelect) SetValues(values actions.ActionValues) {
	t.origin = values.Source.Coords
	t.Count = values.Targets
	t.MaxRange = values.Range
}

func (t *TargetSelect) Update() {
	if !t.isDone {
		x := t.input.Coords.X
		y := t.input.Coords.Y
		inRange := world.DistanceHex(t.origin.X, t.origin.Y, x, y) <= t.MaxRange
		occ := floor.CurrentFloor.GetOccupant(t.input.Coords)
		if occ != nil {
			if _, ok := occ.(objects.Targetable); ok {
				if !inRange {
					// todo: highlight
				} else if t.input.Select.JustPressed() {
					t.input.Select.Consume()
					t.clicked = append(t.clicked, t.input.Coords)
					t.isDone = true
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
	if t.input.Cancel.JustPressed() {
		t.input.Cancel.Consume()
		// cancel
		t.cancel = true
	}
}

func (t *TargetSelect) IsCancelled() bool {
	return t.cancel
}

func (t *TargetSelect) IsDone() bool {
	return len(t.clicked) > 0
}

func (t *TargetSelect) Finish() []world.Coords {
	return t.clicked
}