package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/action"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/pkg/sfx"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type DamageAction struct {
	*action.AbstractAction
	values selector.ActionValues
	target *floor.Character
	coords world.Coords
	start  bool
	preDam bool
}

func NewDamageAction(area []world.Coords, values selector.ActionValues) *DamageAction {
	if len(area) < 1 {
		return nil
	}
	target := floor.CurrentFloor.Get(area[0]).GetOccupant()
	if target != nil {
		return &DamageAction{
			values: values,
			target: target,
			coords: area[0],
			start:  true,
			preDam: true,
		}
	}
	return nil
}

func (a *DamageAction) Update() {
	if a.start {
		SetAttackTransformSingle(a.values.Source, a.coords)
		sfx.SoundPlayer.PlaySound("punch_hit")
		a.start = false
	}
	if !a.values.Source.IsMoving() {
		if a.preDam {
			SetResetTransform(a.values.Source)
			// todo: add an effect

			// todo: this is where the damage modification happens
			a.target.Damage(a.values.Damage)
			a.preDam = false
		} else {
			a.IsDone = true
		}
	}
}

func (a *DamageAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}

type DamageHexAction struct {
	*action.AbstractAction
	values selector.ActionValues
	area   []world.Coords
	start  bool
	preDam bool
}

func NewDamageHexAction(area []world.Coords, values selector.ActionValues) *DamageHexAction {
	if len(area) > 0 {
		return &DamageHexAction{
			values: values,
			area:   area,
			start:  true,
			preDam: true,
		}
	}
	return nil
}

func (a *DamageHexAction) Update() {
	if a.start {
		SetAttackTransform(a.values.Source, a.area)
		a.start = false
	}
	if !a.values.Source.IsMoving() {
		if a.preDam {
			first := true
			SetResetTransform(a.values.Source)
			for _, h := range a.area {
				// todo: add an effect

				// todo: this is where the damage modification happens
				if cha := floor.CurrentFloor.Get(h).GetOccupant(); cha != nil {
					if first {
						sfx.SoundPlayer.PlaySound("punch_hit")
						first = false
					}
					cha.Damage(a.values.Damage)
				}
			}
			if first {
				sfx.SoundPlayer.PlaySound("punch_miss")
			}
			a.preDam = false
		} else {
			a.IsDone = true
		}
	}
}

func (a *DamageHexAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}
