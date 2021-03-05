package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

// A Selector can be updated, checked for completion, cancelled, and finished.
type Selector interface {
	Init(*input.Input)
	SetValues(ActionValues)
	Update()
	IsCancelled() bool
	IsDone() bool
	Finish() []world.Coords
}

type ActionValues struct {
	Source  *characters.Character
	Damage  int
	Move    int
	Range   int
	Targets int
}

type TargetArea struct {
	SetArea func(int, int, world.Coords, floor.PathChecks) []world.Coords
	area    []world.Coords
}

func (t *TargetArea) GetArea() []world.Coords {
	return t.area
}

func SingleTile(_, _ int, orig world.Coords, check floor.PathChecks) []world.Coords {
	if floor.CurrentFloor.IsLegal(orig, check) != nil {
		return []world.Coords{orig}
	}
	return []world.Coords{}
}

func BasicHexArea(rad, _ int, orig world.Coords, check floor.PathChecks) []world.Coords {
	return floor.CurrentFloor.AllWithin(orig, rad, check)
}