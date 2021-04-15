package selector

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"image/color"
)

type AttackEffect struct {
	*AbstractSelectionEffect
	imd *imdraw.IMDraw
}

func (e *AttackEffect) Update() {
	if e.imd == nil {
		e.imd = imdraw.New(nil)
	}
	e.imd.Clear()
	e.imd.Color = color.RGBA{
		R: 255,
		G: 30,
		B: 10,
		A: 175,
	}
	e.imd.EndShape = imdraw.SharpEndShape
	for _, c := range e.area {
		e.imd.Push(world.MapToWorld(e.orig), world.MapToWorld(c))
		e.imd.Line(12)
	}
}

func (e *AttackEffect) Draw(target pixel.Target) {
	for _, c := range e.area {
		mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(world.MapToWorld(c))
		SelectionSprites["attack"].Draw(target, mat)
	}
	e.imd.Draw(target)
}

func (e *AttackEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
