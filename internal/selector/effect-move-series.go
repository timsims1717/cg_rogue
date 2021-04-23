package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type MoveSeriesEffect struct {
	*AbstractSelectionEffect
}

func (e *MoveSeriesEffect) Update() {
	if len(e.area) > 1 {
		for i, c := range e.area {
			spr := new(pixel.Sprite)
			var mat pixel.Matrix
			if i == 0 {
				next := e.area[i+1]
				dir := c.Direction(next)
				if dir == world.LineUp || dir == world.LineDown {
					spr = gfx.SelectionSprites["move_start_u"]
				} else {
					spr = gfx.SelectionSprites["move_start_ul"]
				}
				switch dir {
				case world.LineDown, world.LineDownLeft:
					mat = img.Flop
				case world.LineDownRight:
					mat = img.FlipFlop
				case world.LineUpRight:
					mat = img.Flip
				default:
					mat = img.IM
				}
			} else if len(e.area) - 1 == i {
				prev := e.area[i-1]
				dir := c.Direction(prev)
				if dir == world.LineUp || dir == world.LineDown {
					spr = gfx.SelectionSprites["move_end_u"]
				} else {
					spr = gfx.SelectionSprites["move_end_ul"]
				}
				switch dir {
				case world.LineDown, world.LineDownLeft:
					mat = img.Flop
				case world.LineDownRight:
					mat = img.FlipFlop
				case world.LineUpRight:
					mat = img.Flip
				default:
					mat = img.IM
				}
			} else {
				next := e.area[i+1]
				prev := e.area[i-1]
				dirN := c.Direction(next)
				dirP := c.Direction(prev)
				diff := util.Abs(int(dirN - dirP)) % world.TopLeft
				if diff == 2 || diff == 10 {
					if dirN == world.LineUp || dirN == world.LineDown || dirP == world.LineUp || dirP == world.LineDown {
						spr = gfx.SelectionSprites["move_ul_u"]
						if dirN == world.LineUpRight || dirP == world.LineUpRight {
							mat = img.Flip
						} else if dirN == world.LineDownRight || dirP == world.LineDownRight {
							mat = img.FlipFlop
						} else if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							mat = img.Flop
						} else {
							mat = img.IM
						}
					} else {
						spr = gfx.SelectionSprites["move_ul_dl"]
						if dirN == world.LineUpRight || dirP == world.LineUpRight {
							mat = img.Flip
						} else {
							mat = img.IM
						}
					}
				} else if diff == 4 || diff == 8 {
					if dirN == world.LineUp || dirN == world.LineDown || dirP == world.LineUp || dirP == world.LineDown {
						spr = gfx.SelectionSprites["move_ul_d"]
						if dirN == world.LineUpRight || dirP == world.LineUpRight {
							mat = img.Flip
						} else if dirN == world.LineDownRight || dirP == world.LineDownRight {
							mat = img.FlipFlop
						} else if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							mat = img.Flop
						} else {
							mat = img.IM
						}
					} else {
						spr = gfx.SelectionSprites["move_ul_ur"]
						if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							mat = img.Flop
						} else {
							mat = img.IM
						}
					}
				} else {
					if dirN == world.LineUp || dirP == world.LineUp {
						spr = gfx.SelectionSprites["move_u_d"]
						mat = img.IM
					} else {
						spr = gfx.SelectionSprites["move_ul_dr"]
						if dirN == world.LineDownLeft || dirP == world.LineDownLeft {
							mat = img.Flip
						} else {
							mat = img.IM
						}
					}
				}
			}
			hex := floor.CurrentFloor.Get(c)
			hex.AddEffect([]img.Sprite{{
				S: spr,
				M: mat,
			}}, 11)
		}
	}
}

func (e *MoveSeriesEffect) Draw(_ pixel.Target) {}

func (e *MoveSeriesEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
