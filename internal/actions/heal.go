package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HealAction struct {
	values selectors.ActionValues
	target *characters.Character
	coords world.Coords
	start  bool
	isDone bool
}

func NewHealAction(area []world.Coords, values selectors.ActionValues) *HealAction {
	if len(area) < 1 {
		return nil
	}
	occ := floor.CurrentFloor.GetOccupant(area[0])
	if occ != nil {
		if target, ok := occ.(*characters.Character); ok {
			return &HealAction{
				values: values,
				target: target,
				coords: area[0],
				start:  true,
				isDone: false,
			}
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