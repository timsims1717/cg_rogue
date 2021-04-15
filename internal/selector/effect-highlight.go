package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HighlightEffect struct {
	*AbstractSelectionEffect
}

func (e *HighlightEffect) Update() {

}

func (e *HighlightEffect) Draw(target pixel.Target) {
	for _, c := range e.area {
		mat := pixel.IM.Moved(world.MapToWorld(c))
		SelectionSprites["default"].Draw(target, mat)
	}
}

func (e *HighlightEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
