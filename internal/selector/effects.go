package selector

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type SelectionEffect interface{
	Update()
	Draw(pixel.Target)
	SetAbstract(*AbstractSelectionEffect)
}

func NewSelectionEffect(effect SelectionEffect) *AbstractSelectionEffect {
	nEffect := &AbstractSelectionEffect{
		Effect: effect,
	}
	effect.SetAbstract(nEffect)
	return nEffect
}

type AbstractSelectionEffect struct {
	Effect SelectionEffect
	area   []world.Coords
}

func (e *AbstractSelectionEffect) SetArea(area []world.Coords) {
	e.area = area
}

var SelectionSprites map[string]*pixel.Sprite

var SelectionSet selectionSet

type selectionSet struct {
	sprites []*pixel.Sprite
	batch   *pixel.Batch
	set     []selectUI
	nset    []SelectionEffect
}

type SelectionType int

const (
	Default SelectionType = iota
	Move
	MoveSolid
	Attack
	Blank
)

type selectUI struct {
	t SelectionType
	x int
	y int
}

func (sel *selectUI) getCoords() world.Coords {
	return world.Coords{
		X: sel.x,
		Y: sel.y,
	}
}

func AddSelectUI(t SelectionType, x, y int) {
	SelectionSet.set = append(SelectionSet.set, selectUI{
		t: t,
		x: x,
		y: y,
	})
}

func AddSelectionEffect(effect *AbstractSelectionEffect) {
	SelectionSet.nset = append(SelectionSet.nset, effect.Effect)
}

func (s *selectionSet) SetSpriteSheet(sheet *img.SpriteSheet) {
	s.sprites = []*pixel.Sprite{}
	s.batch = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Img)
	for _, r := range sheet.Sprites {
		s.sprites = append(s.sprites, pixel.NewSprite(sheet.Img, r))
	}
	SelectionSprites = make(map[string]*pixel.Sprite )
	for k, r := range sheet.SpriteMap {
		SelectionSprites[k] = pixel.NewSprite(sheet.Img, r)
	}
}

func (s *selectionSet) Draw(win *pixelgl.Window) {
	s.batch.Clear()
	for _, sel := range s.set {
		if sel.t != Blank {
			mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(world.MapToWorld(sel.getCoords()))
			s.sprites[sel.t].Draw(s.batch, mat)
		}
	}
	for _, sel := range s.nset {
		sel.Draw(s.batch)
	}
	s.batch.Draw(win)
	s.set = []selectUI{}
	s.nset = []SelectionEffect{}
}