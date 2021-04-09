package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HealAction struct {
	values selector.ActionValues
	target *floor.Character
	coords world.Coords
	start  bool
	isDone bool
}

func NewHealAction(area []world.Coords, values selector.ActionValues) *HealAction {
	if len(area) < 1 {
		return nil
	}
	target := floor.CurrentFloor.GetOccupant(area[0])
	if target != nil {
		return &HealAction{
			values: values,
			target: target,
			coords: area[0],
			start:  true,
			isDone: false,
		}
	}
	return nil
}

func (a *HealAction) Update() {
	a.target.Heal(a.values.Heal)
	a.isDone = true
}

func (a *HealAction) IsDone() bool {
	return a.isDone
}