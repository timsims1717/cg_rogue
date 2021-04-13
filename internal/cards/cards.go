package cards

import (
	"github.com/timsims1717/cg_rogue_go/internal/action"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
)

func AddToTop(a action.Action, effects... *selector.AbstractSelectionEffect) {
	action.ActionManager.AddToTop(a, effects)
}

func AddToBot(a action.Action, effects... *selector.AbstractSelectionEffect) {
	action.ActionManager.AddToBot(a, effects)
}
