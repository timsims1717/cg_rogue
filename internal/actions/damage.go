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