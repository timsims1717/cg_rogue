package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
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
		e.areaMatrices, e.areaSprites = attackAreas(e.area)
		for i, c := range e.area {
			if world.CoordsIn(c, e.orig.Neighbors(floor.CurrentFloor.Dimensions())) {
				dir := c.Direction(e.orig)
				mat := pixel.IM
				spr := new(pixel.Sprite)
				switch dir {
				case world.LineUpLeft:
					spr = SelectionSprites["attack_arrow_ul"]
				case world.LineUp:
					spr = SelectionSprites["attack_arrow_u"]
				case world.LineUpRight:
					spr = SelectionSprites["attack_arrow_ul"]
					mat = pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
				case world.LineDownLeft:
					spr = SelectionSprites["attack_arrow_ul"]
					mat = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.))
				case world.LineDown:
					spr = SelectionSprites["attack_arrow_u"]
					mat = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.))
				case world.LineDownRight:
					spr = SelectionSprites["attack_arrow_ul"]
					mat = pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.))
				}
				if spr != nil {
					e.areaMatrices[i] = append(e.areaMatrices[i], mat)
					e.areaSprites[i] = append(e.areaSprites[i], spr)
				}
			}
		}
	}
}

func (e *AttackTriangleEffect) Draw(target pixel.Target) {
	if len(e.area) > 0 {
		for i, c := range e.area {
			SelectionSprites["attack_bg"].Draw(target, pixel.IM.Moved(world.MapToWorld(c)))
			for j, m := range e.areaMatrices[i] {
				e.areaSprites[i][j].Draw(target, m.Moved(world.MapToWorld(c)))
			}
		}
	}
}

func (e *AttackTriangleEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
