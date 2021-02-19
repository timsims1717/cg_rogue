package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

// A Selector can be updated, checked for completion, cancelled, and finished.
type Selector interface {
	Init(world.Coords, *input.Input, actions.ActionValues)
	Update()
	IsCancelled() bool
	IsDone() bool
	Finish() []world.Coords
}