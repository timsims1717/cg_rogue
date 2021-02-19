package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
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

func (d *DamageAction) Update() {
	if d.start {
		// todo: add an effect

		// todo: this is where the damage modification happens
		d.target.Damage(d.dmg)
		d.start = false
		d.isDone = true
	}
}

func (d *DamageAction) IsDone() bool {
	return d.isDone
}