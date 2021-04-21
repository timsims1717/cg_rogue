package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/action"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
)

type DefenseAction struct {
	*action.AbstractAction
	values selector.ActionValues
	target *floor.Character
	start  bool
}

func NewDefenseAction(target *floor.Character, values selector.ActionValues) *DefenseAction {
	if target == nil {
		return nil
	}
	return &DefenseAction{
		values:         values,
		target:         target,
		start:          true,
	}
}

func (a *DefenseAction) Update() {
	a.target.Defense.Add(a.values.Defense)
	a.IsDone = true
}

func (a *DefenseAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}