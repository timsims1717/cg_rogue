package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
)

type DefenseEffect struct {
	*AbstractSelectionEffect
}

func (e *DefenseEffect) Update() {
	for _, c := range e.area {
		hex := floor.CurrentFloor.Get(c)
		hex.AddEffect([]img.Sprite{{
			S: gfx.SelectionSprites["defense"],
			M: img.IM,
		}}, 0)
	}
}

func (e *DefenseEffect) Draw(_ pixel.Target) {}

func (e *DefenseEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
