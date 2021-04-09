package cards

import "github.com/timsims1717/cg_rogue_go/internal/manager"

func AddToTop(a manager.Action) {
	manager.ActionManager.AddToTop(a)
}

func AddToBot(a manager.Action) {
	manager.ActionManager.AddToBot(a)
}