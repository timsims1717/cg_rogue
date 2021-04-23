package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
)

type HighlightEffect struct {
	*AbstractSelectionEffect
}

func (e *HighlightEffect) Update() {

}

func (e *HighlightEffect) Draw(target pixel.Target) {
	for _, c := range e.area {
		hex := floor.CurrentFloor.Get(c)
		hex.AddEffect([]img.Sprite{{
			S: gfx.SelectionSprites["default"],
			M: img.IM,
		}}, 6)
	}
}

func (e *HighlightEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
