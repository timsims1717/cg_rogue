package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/action"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HealAction struct {
	*action.AbstractAction
	values selector.ActionValues
	target *floor.Character
	coords world.Coords
	start  bool
}

func NewHealAction(area []world.Coords, values selector.ActionValues) *HealAction {
	if len(area) < 1 {
		return nil
	}
	target := floor.CurrentFloor.Get(area[0]).GetOccupant()
	if target != nil {
		return &HealAction{
			values: values,
			target: target,
			coords: area[0],
			start:  true,
		}
	}
	return nil
}

func (a *HealAction) Update() {
	a.target.Heal(a.values.Heal)
	a.IsDone = true
}

func (a *HealAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}
