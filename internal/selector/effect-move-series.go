package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type MoveSeriesEffect struct {
	*AbstractSelectionEffect
	sprites  []*pixel.Sprite
	matrices []pixel.Matrix
}

func (e *MoveSeriesEffect) Update() {
	e.sprites = []*pixel.Sprite{}
	e.matrices = []pixel.Matrix{}
	if len(e.area) > 1 {
		for i, c := range e.area {
			if i == 0 {
				next := e.area[i+1]
				dir := c.Direction(next)
				if dir == world.LineUp || dir == world.LineDown {
					e.sprites = append(e.sprites, SelectionSprites["move_start_u"])
				} else {
					e.sprites = append(e.sprites, SelectionSprites["move_start_ul"])
				}
				switch dir {
				case world.LineDown, world.LineDownLeft:
					e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.)))
				case world.LineDownRight:
					e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.)))
				case world.LineUpRight:
					e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.)))
				default:
					e.matrices = append(e.matrices, pixel.IM)
				}
			} else if len(e.area) - 1 == i {
				prev := e.area[i-1]
				dir := c.Direction(prev)
				if dir == world.LineUp || dir == world.LineDown {
					e.sprites = append(e.sprites, SelectionSprites["move_end_u"])
				} else {
					e.sprites = append(e.sprites, SelectionSprites["move_end_ul"])
				}
				switch dir {
				case world.LineDown, world.LineDownLeft:
					e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.)))
				case world.LineDownRight:
					e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.)))
				case world.LineUpRight:
					e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.)))
				default:
					e.matrices = append(e.matrices, pixel.IM)
				}
			} else {
				next := e.area[i+1]
				prev := e.area[i-1]
				dirN := c.Direction(next)
				dirP := c.Direction(prev)
				diff := util.Abs(int(dirN - dirP)) % world.TopLeft
				if diff == 2 || diff == 10 {
					if dirN == world.LineUp || dirN == world.LineDown || dirP == world.LineUp || dirP == world.LineDown {
						e.sprites = append(e.sprites, SelectionSprites["move_ul_u"])
						if dirN == world.LineUpRight || dirP == world.LineUpRight {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.)))
						} else if dirN == world.LineDownRight || dirP == world.LineDownRight {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.)))
						} else if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.)))
						} else {
							e.matrices = append(e.matrices, pixel.IM)
						}
					} else {
						e.sprites = append(e.sprites, SelectionSprites["move_ul_dl"])
						if dirN == world.LineUpRight || dirP == world.LineUpRight {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.)))
						} else {
							e.matrices = append(e.matrices, pixel.IM)
						}
					}
					e.matrices = append(e.matrices, pixel.IM)
				} else if diff == 4 || diff == 8 {
					if dirN == world.LineUp || dirN == world.LineDown || dirP == world.LineUp || dirP == world.LineDown {
						e.sprites = append(e.sprites, SelectionSprites["move_ul_d"])
						if dirN == world.LineUpRight || dirP == world.LineUpRight {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.)))
						} else if dirN == world.LineDownRight || dirP == world.LineDownRight {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.)))
						} else if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.)))
						} else {
							e.matrices = append(e.matrices, pixel.IM)
						}
					} else {
						e.sprites = append(e.sprites, SelectionSprites["move_ul_ur"])
						if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.)))
						} else {
							e.matrices = append(e.matrices, pixel.IM)
						}
					}
				} else {
					if dirN == world.LineUp || dirP == world.LineUp {
						e.sprites = append(e.sprites, SelectionSprites["move_u_d"])
						e.matrices = append(e.matrices, pixel.IM)
					} else {
						e.sprites = append(e.sprites, SelectionSprites["move_ul_dr"])
						if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							e.matrices = append(e.matrices, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.)))
						} else {
							e.matrices = append(e.matrices, pixel.IM)
						}
					}
				}
			}
		}
	}
}

func (e *MoveSeriesEffect) Draw(target pixel.Target) {
	if len(e.area) > 1 {
		for i, c := range e.area {
			mat := e.matrices[i].Moved(world.MapToWorld(c))
			e.sprites[i].Draw(target, mat)
		}
	}
}

func (e *MoveSeriesEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
