package selector

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"image/color"
)

type MoveEffect struct {
	*AbstractSelectionEffect
	imd *imdraw.IMDraw
}

func (e *MoveEffect) Update() {
	if e.imd == nil {
		e.imd = imdraw.New(nil)
	}
	e.imd.Clear()
	e.imd.Color = color.RGBA{
		R: 10,
		G: 30,
		B: 255,
		A: 175,
	}
	e.imd.EndShape = imdraw.SharpEndShape
	prev := e.orig
	for _, c := range e.area {
		if c != prev {
			e.imd.Push(world.MapToWorld(prev), world.MapToWorld(c))
		}
		e.imd.Line(12)
		prev = c
	}
}

func (e *MoveEffect) Draw(target pixel.Target) {
	for i, c := range e.area {
		mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(world.MapToWorld(c))
		if i != len(e.area)-1 {
			SelectionSprites["move"].Draw(target, mat)
		} else {
			SelectionSprites["move-solid"].Draw(target, mat)
		}
	}
	e.imd.Draw(target)
}

func (e *MoveEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
