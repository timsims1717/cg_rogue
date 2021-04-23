package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

func attackAreas(area []world.Coords) [][]img.Sprite {
	var sprites [][]img.Sprite
	for i, c := range area {
		sprites = append(sprites, []img.Sprite{{
				S: gfx.SelectionSprites["attack_bg"],
				M: img.IM,
			}})
		nbors := world.Intersection(c.Neighbors(), area)
		if len(nbors) == 0 {
			sprites[i] = append(sprites[i], img.Sprite{
				S: gfx.SelectionSprites["attack_single"],
				M: img.IM,
			})
		} else {
			top := map[int]bool{
				0: false,
				1: false,
				2: false,
			}
			bot := map[int]bool{
				0: false,
				1: false,
				2: false,
			}
			for _, n := range nbors {
				dir := c.Direction(n)
				switch dir {
				case world.LineUpLeft:
					top[0] = true
				case world.LineUp:
					top[1] = true
				case world.LineUpRight:
					top[2] = true
				case world.LineDownLeft:
					bot[0] = true
				case world.LineDown:
					bot[1] = true
				case world.LineDownRight:
					bot[2] = true
				}
			}
			if top[1] {
				if !top[0] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul"],
						M: img.IM,
					})
				}
				if !top[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul"],
						M: img.Flip,
					})
				}
			} else {
				if !top[0] && !top[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul_u_ur"],
						M: img.IM,
					})
				} else if !top[0] && top[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul_u"],
						M: img.IM,
					})
				} else if top[0] && !top[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul_u"],
						M: img.Flip,
					})
				} else {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_u"],
						M: img.IM,
					})
				}
			}
			if bot[1] {
				if !bot[0] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul"],
						M: img.Flop,
					})
				}
				if !bot[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul"],
						M: img.FlipFlop,
					})
				}
			} else {
				if !bot[0] && !bot[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul_u_ur"],
						M: img.Flop,
					})
				} else if !bot[0] && bot[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul_u"],
						M: img.Flop,
					})
				} else if bot[0] && !bot[2] {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_ul_u"],
						M: img.FlipFlop,
					})
				} else {
					sprites[i] = append(sprites[i], img.Sprite{
						S: gfx.SelectionSprites["attack_u"],
						M: img.Flop,
					})
				}
			}
		}
	}
	return sprites
}
