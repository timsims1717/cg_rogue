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
	IsDone() bool
	Finish() []world.Coords
}

type ActionValues struct {
	Source   *characters.Character
	Damage   int
	Move     int
	Range    int
	Targets  int
	Strength int
	Heal     int
	Checks   floor.PathChecks
}