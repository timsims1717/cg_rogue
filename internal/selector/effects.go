package selector

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type SelectionEffect interface {
	Update()
	Draw(pixel.Target)
	SetAbstract(*AbstractSelectionEffect)
}

func NewSelectionEffect(effect SelectionEffect, values ActionValues) *AbstractSelectionEffect {
	nEffect := &AbstractSelectionEffect{
		Effect: effect,
		values: values,
	}
	effect.SetAbstract(nEffect)
	return nEffect
}

type AbstractSelectionEffect struct {
	Effect SelectionEffect
	area   []world.Coords
	values ActionValues
	orig   world.Coords
}

func (e *AbstractSelectionEffect) SetOrig(orig world.Coords) {
	e.orig = orig
}

func (e *AbstractSelectionEffect) SetValues(values ActionValues) {
	e.values = values
}

func (e *AbstractSelectionEffect) SetArea(area []world.Coords) {
	e.area = area
}

var SelectionSprites map[string]*pixel.Sprite

var SelectionSet selectionSet

type selectionSet struct {
	sprites []*pixel.Sprite
	batch   *pixel.Batch
	nset    []SelectionEffect
}

func AddSelectionEffect(effect *AbstractSelectionEffect) {
	if effect != nil {
		SelectionSet.nset = append(SelectionSet.nset, effect.Effect)
	}
}

func (s *selectionSet) SetSpriteSheet(sheet *img.SpriteSheet) {
	s.sprites = []*pixel.Sprite{}
	s.batch = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Img)
	for _, r := range sheet.Sprites {
		s.sprites = append(s.sprites, pixel.NewSprite(sheet.Img, r))
	}
	SelectionSprites = make(map[string]*pixel.Sprite)
	for k, r := range sheet.SpriteMap {
		SelectionSprites[k] = pixel.NewSprite(sheet.Img, r)
	}
}

func (s *selectionSet) Update() {
	for _, sel := range s.nset {
		sel.Update()
	}
}

func (s *selectionSet) Draw(win *pixelgl.Window) {
	s.batch.Clear()
	for _, sel := range s.nset {
		sel.Draw(s.batch)
	}
	s.batch.Draw(win)
	s.nset = []SelectionEffect{}
}

func attackAreas(area []world.Coords) ([][]pixel.Matrix, [][]*pixel.Sprite) {
	var areaMatrices [][]pixel.Matrix
	var areaSprites [][]*pixel.Sprite
	for i, c := range area {
		areaMatrices = append(areaMatrices, []pixel.Matrix{})
		areaSprites = append(areaSprites, []*pixel.Sprite{})
		mat := pixel.IM
		flip := pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
		flop := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.))
		flipflop := pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.))
		nbors := world.Intersection(c.Neighbors(floor.CurrentFloor.Dimensions()), area)
		if len(nbors) == 0 {
			areaMatrices[i] = append(areaMatrices[i], mat)
			areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_single"])
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
					areaMatrices[i] = append(areaMatrices[i], mat)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul"])
				}
				if !top[2] {
					areaMatrices[i] = append(areaMatrices[i], flip)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul"])
				}
			} else {
				if !top[0] && !top[2] {
					areaMatrices[i] = append(areaMatrices[i], mat)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul_u_ur"])
				} else if !top[0] && top[2] {
					areaMatrices[i] = append(areaMatrices[i], mat)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul_u"])
				} else if top[0] && !top[2] {
					areaMatrices[i] = append(areaMatrices[i], flip)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul_u"])
				} else {
					areaMatrices[i] = append(areaMatrices[i], mat)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_u"])
				}
			}
			if bot[1] {
				if !bot[0] {
					areaMatrices[i] = append(areaMatrices[i], flop)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul"])
				}
				if !bot[2] {
					areaMatrices[i] = append(areaMatrices[i], flipflop)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul"])
				}
			} else {
				if !bot[0] && !bot[2] {
					areaMatrices[i] = append(areaMatrices[i], flop)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul_u_ur"])
				} else if !bot[0] && bot[2] {
					areaMatrices[i] = append(areaMatrices[i], flop)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul_u"])
				} else if bot[0] && !bot[2] {
					areaMatrices[i] = append(areaMatrices[i], flipflop)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_ul_u"])
				} else {
					areaMatrices[i] = append(areaMatrices[i], flop)
					areaSprites[i] = append(areaSprites[i], SelectionSprites["attack_u"])
				}
			}
		}
	}
	return areaMatrices, areaSprites
}