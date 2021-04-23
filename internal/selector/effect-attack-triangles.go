package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type AttackTriangleEffect struct {
	*AbstractSelectionEffect
	areaMatrices   [][]pixel.Matrix
	areaSprites    [][]*pixel.Sprite
}

func (e *AttackTriangleEffect) Update() {
	e.areaMatrices = [][]pixel.Matrix{}
	e.areaSprites = [][]*pixel.Sprite{}
	if len(e.area) > 0 {
		sprites := attackAreas(e.area)
		for i, c := range e.area {
			if world.CoordsIn(c, e.orig.Neighbors()) {
				dir := c.Direction(e.orig)
				mat := img.IM
				spr := new(pixel.Sprite)
				switch dir {
				case world.LineUpLeft:
					spr = gfx.SelectionSprites["attack_arrow_ul"]
				case world.LineUp:
					spr = gfx.SelectionSprites["attack_arrow_u"]
				case world.LineUpRight:
					spr = gfx.SelectionSprites["attack_arrow_ul"]
					mat = img.Flip
				case world.LineDownLeft:
					spr = gfx.SelectionSprites["attack_arrow_ul"]
					mat = img.Flop
				case world.LineDown:
					spr = gfx.SelectionSprites["attack_arrow_u"]
					mat = img.Flop
				case world.LineDownRight:
					spr = gfx.SelectionSprites["attack_arrow_ul"]
					mat = img.FlipFlop
				}
				if spr != nil {
					sprites[i] = append(sprites[i], img.Sprite{
						S: spr,
						M: mat,
					})
				}
			}
			hex := floor.CurrentFloor.Get(c)
			hex.AddEffect(sprites[i], 10)
		}
	}
}

func (e *AttackTriangleEffect) Draw(_ pixel.Target) {}

func (e *AttackTriangleEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
