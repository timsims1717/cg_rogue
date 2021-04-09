package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type MoveEffect struct{
	*AbstractSelectionEffect
}

func (e *MoveEffect) Update() {

}

func (e *MoveEffect) Draw(target pixel.Target) {
	for i, c := range e.area {
		mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(world.MapToWorld(c))
		if i != len(e.area)-1 {
			SelectionSprites["move"].Draw(target, mat)
		} else {
			SelectionSprites["move-solid"].Draw(target, mat)
		}
	}
}

func (e *MoveEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}