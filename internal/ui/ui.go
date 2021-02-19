package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"sync"
)

var SelectionSet selectionSet

type selectionSet struct {
	mu      sync.Mutex
	sprites []*pixel.Sprite
	batch   *pixel.Batch
	set     []selectUI
}

type SelectionType int

const (
	Default SelectionType = iota
	Move
	MoveSolid
	Blank
)

type selectUI struct {
	t SelectionType
	x int
	y int
}

func AddSelectUI(t SelectionType, x, y int) {
	SelectionSet.mu.Lock()
	defer SelectionSet.mu.Unlock()
	SelectionSet.set = append(SelectionSet.set, selectUI{
		t: t,
		x: x,
		y: y,
	})
}

func (s *selectionSet) SetSpriteSheet(sheet *img.SpriteSheet) {
	s.sprites = []*pixel.Sprite{}
	s.batch = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Img)
	for _, r := range sheet.Sprites {
		s.sprites = append(s.sprites, pixel.NewSprite(sheet.Img, r))
	}

}

func (s *selectionSet) Draw(win *pixelgl.Window) {
	s.batch.Clear()
	for _, sel := range s.set {
		if sel.t != Blank {
			mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(pixel.V(world.MapToWorldHex(sel.x, sel.y)))
			s.sprites[sel.t].Draw(s.batch, mat)
		}
	}
	s.batch.Draw(win)
	s.set = []selectUI{}
}