package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type DefenseEffect struct {
	*AbstractSelectionEffect
}

func (e *DefenseEffect) Update() {

}

func (e *DefenseEffect) Draw(target pixel.Target) {
	for _, c := range e.area {
		mat := pixel.IM.Moved(world.MapToWorld(c))
		SelectionSprites["defense"].Draw(target, mat)
	}
}

func (e *DefenseEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
