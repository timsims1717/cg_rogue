package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type DamageAction struct {
	source *characters.Character
	target objects.Targetable
	dmg    int
	start  bool
	isDone bool
	interR *gween.Tween
}

func NewDamageAction(s *characters.Character, t objects.Targetable, d int) *DamageAction {
	if t == nil {
		return nil
	}
	return &DamageAction{
		source: s,
		target: t,
		dmg:    d,
		start:  true,
		isDone: false,
	}
}

func (a *DamageAction) Update() {
	if a.start {
		// todo: add an effect

		// todo: this is where the damage modification happens
		a.target.Damage(a.dmg)
		a.start = false
		a.isDone = true
	}
}

func (a *DamageAction) IsDone() bool {
	return a.isDone
}

type DamageHexAction struct {
	source *characters.Character
	area  []world.Coords
	dmg    int
	start  bool
	isDone bool
	interR *gween.Tween
}

func NewDamageHexAction(s *characters.Character, area []world.Coords, d int) *DamageHexAction {
	return &DamageHexAction{
		source: s,
		area:   area,
		dmg:    d,
		start:  true,
		isDone: false,
	}
}

func (a *DamageHexAction) Update() {
	if a.start {
		for _, h := range a.area {
			// todo: add an effect

			// todo: this is where the damage modification happens
			if o := floor.CurrentFloor.GetOccupant(h); o != nil {
				if t, ok := o.(objects.Targetable); ok {
					t.Damage(a.dmg)
				}
			}
		}
		a.start = false
		a.isDone = true
	}
}

func (a *DamageHexAction) IsDone() bool {
	return a.isDone
}