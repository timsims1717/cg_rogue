package floor

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"golang.org/x/image/colornames"
)

type Defense struct {
	CurrDef int
	MaxDef  int
	LastDmg int
	Alive   bool
	Display bool
	Hovered bool

	imd *imdraw.IMDraw
	pos pixel.Vec
}

func (d *Defense) Add(amt int) {
	if amt < 0 {
		return
	}
	d.CurrDef += amt
	if d.CurrDef > d.MaxDef {
		d.CurrDef = d.MaxDef
	}
}

func (d *Defense) RemoveAll() {
	d.CurrDef = 0
}

func (d *Defense) Damage(amt int) int {
	amount := amt
	if amount < 0 {
		amount = 0
	}
	if amount > 0 {
		d.LastDmg = util.Min(amount, d.CurrDef)
		d.CurrDef -= d.LastDmg
		return amount - d.LastDmg
	}
	return 0
}

func (d *Defense) Update() {
	if d.Alive && d.Display && d.CurrDef > 0 {
		d.imd.Clear()
		d.imd.Color = colornames.Darkgray
		d.imd.EndShape = imdraw.NoEndShape
		d.imd.Push(pixel.V(d.pos.X-10., d.pos.Y))
		d.imd.Push(pixel.V(d.pos.X-10., d.pos.Y+5.))
		d.imd.Push(pixel.V(d.pos.X+10., d.pos.Y+5.))
		d.imd.Push(pixel.V(d.pos.X+10., d.pos.Y))
		d.imd.Polygon(0.)
		d.imd.Color = colornames.Lightgray
		d.imd.EndShape = imdraw.NoEndShape
		d.imd.Push(pixel.V(d.pos.X-10., d.pos.Y))
		d.imd.Push(pixel.V(d.pos.X-10., d.pos.Y+5.))
		d.imd.Push(pixel.V(d.pos.X+10., d.pos.Y+5.))
		d.imd.Push(pixel.V(d.pos.X+10., d.pos.Y))
		d.imd.Polygon(2.)
		perc := int((18. / float64(d.MaxDef)) * float64(d.CurrDef))
		d.imd.Color = colornames.Mediumblue
		d.imd.EndShape = imdraw.NoEndShape
		d.imd.Push(pixel.V(d.pos.X-9., d.pos.Y+1.))
		d.imd.Push(pixel.V(d.pos.X-9., d.pos.Y+4.))
		d.imd.Push(pixel.V(d.pos.X+float64(perc)-9., d.pos.Y+4.))
		d.imd.Push(pixel.V(d.pos.X+float64(perc)-9., d.pos.Y+1.))
		d.imd.Polygon(0.)
	}
}

func (d *Defense) Draw(win *pixelgl.Window) {
	if d.Alive && d.Display && d.CurrDef > 0 {
		d.imd.Draw(win)
	}
}