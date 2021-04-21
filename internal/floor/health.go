package floor

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"golang.org/x/image/colornames"
)

type Health struct {
	CurrHP  int
	MaxHP   int
	LastDmg int
	Alive   bool
	Display bool
	Hovered bool

	imd *imdraw.IMDraw
	pos pixel.Vec
}

func (h *Health) Damage(amt int) int {
	amount := amt
	if amount < 0 {
		amount = 0
	}
	if amount > 0 {
		h.LastDmg = util.Min(amount, h.CurrHP)
		h.CurrHP -= h.LastDmg
		return amount - h.LastDmg
	}
	return 0
}

func (h *Health) Update() {
	if h.Alive && h.Display && (h.Hovered || h.CurrHP < h.MaxHP) {
		h.imd.Clear()
		h.imd.Color = colornames.Darkgray
		h.imd.EndShape = imdraw.NoEndShape
		h.imd.Push(pixel.V(h.pos.X-10., h.pos.Y))
		h.imd.Push(pixel.V(h.pos.X-10., h.pos.Y+5.))
		h.imd.Push(pixel.V(h.pos.X+10., h.pos.Y+5.))
		h.imd.Push(pixel.V(h.pos.X+10., h.pos.Y))
		h.imd.Polygon(0.)
		h.imd.Color = colornames.Lightgray
		h.imd.EndShape = imdraw.NoEndShape
		h.imd.Push(pixel.V(h.pos.X-10., h.pos.Y))
		h.imd.Push(pixel.V(h.pos.X-10., h.pos.Y+5.))
		h.imd.Push(pixel.V(h.pos.X+10., h.pos.Y+5.))
		h.imd.Push(pixel.V(h.pos.X+10., h.pos.Y))
		h.imd.Polygon(2.)
		perc := int((18. / float64(h.MaxHP)) * float64(h.CurrHP))
		h.imd.Color = colornames.Darkred
		h.imd.EndShape = imdraw.NoEndShape
		h.imd.Push(pixel.V(h.pos.X-9., h.pos.Y+1.))
		h.imd.Push(pixel.V(h.pos.X-9., h.pos.Y+4.))
		h.imd.Push(pixel.V(h.pos.X+float64(perc)-9., h.pos.Y+4.))
		h.imd.Push(pixel.V(h.pos.X+float64(perc)-9., h.pos.Y+1.))
		h.imd.Polygon(0.)
	}
}

func (h *Health) Draw(win *pixelgl.Window) {
	if h.Alive && h.Display && (h.Hovered || h.CurrHP < h.MaxHP) {
		h.imd.Draw(win)
	}
}