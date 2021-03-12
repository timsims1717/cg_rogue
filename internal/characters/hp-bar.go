package characters

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Health struct {
	CurrHP  int
	MaxHP   int
	LastDmg int
	Alive   bool

	imd *imdraw.IMDraw
	pos pixel.Vec
}

func (h *Health) Update() {
	if h.Alive {
		h.imd.Clear()
		h.imd.Color = colornames.Darkgray
		h.imd.EndShape = imdraw.NoEndShape
		h.imd.Push(pixel.V(h.pos.X-20., h.pos.Y))
		h.imd.Push(pixel.V(h.pos.X-20., h.pos.Y+8.))
		h.imd.Push(pixel.V(h.pos.X+20., h.pos.Y+8.))
		h.imd.Push(pixel.V(h.pos.X+20., h.pos.Y))
		h.imd.Polygon(0.)
		h.imd.Color = colornames.Lightgray
		h.imd.EndShape = imdraw.NoEndShape
		h.imd.Push(pixel.V(h.pos.X-20., h.pos.Y))
		h.imd.Push(pixel.V(h.pos.X-20., h.pos.Y+8.))
		h.imd.Push(pixel.V(h.pos.X+20., h.pos.Y+8.))
		h.imd.Push(pixel.V(h.pos.X+20., h.pos.Y))
		h.imd.Polygon(2.)
		perc := int((38. / float64(h.MaxHP)) * float64(h.CurrHP))
		h.imd.Color = colornames.Darkred
		h.imd.EndShape = imdraw.NoEndShape
		h.imd.Push(pixel.V(h.pos.X-19., h.pos.Y+1.))
		h.imd.Push(pixel.V(h.pos.X-19., h.pos.Y+7.))
		h.imd.Push(pixel.V(h.pos.X+float64(perc)-19., h.pos.Y+7.))
		h.imd.Push(pixel.V(h.pos.X+float64(perc)-19., h.pos.Y+1.))
		h.imd.Polygon(0.)
	}
}

func (h *Health) Draw(win *pixelgl.Window) {
	if h.Alive {
		h.imd.Draw(win)
	}
}