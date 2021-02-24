package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

// A Selector can be updated, checked for completion, cancelled, and finished.
type Selector interface {
	Init(*input.Input)
	SetValues(actions.ActionValues)
	Update()
	IsCancelled() bool
	IsDone() bool
	Finish() []world.Coords
}

type TargetArea struct {
	SetArea func(int, int, world.Coords) []world.Coords
	area    []world.Coords
}

func (t *TargetArea) GetArea() []world.Coords {
	return t.area
}

func BasicHexArea(rad, _ int, orig world.Coords) []world.Coords {
	return floor.CurrentFloor.AllWithin(orig, rad, floor.DefaultCheck)
}