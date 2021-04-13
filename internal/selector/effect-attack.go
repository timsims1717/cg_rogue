package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type AttackEffect struct {
	*AbstractSelectionEffect
}

func (e *AttackEffect) Update() {

}

func (e *AttackEffect) Draw(target pixel.Target) {
	for _, c := range e.area {
		mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(world.MapToWorld(c))
		SelectionSprites["attack"].Draw(target, mat)
	}
}

func (e *AttackEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
