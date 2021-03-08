package text

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"image/color"
)

type TextDisplay struct {
	Text *text.Text
	Raw  string

	Mask   color.RGBA
	interR *gween.Tween
	interG *gween.Tween
	interB *gween.Tween
	interA *gween.Tween

	Pos    pixel.Vec
	interX *gween.Tween
	interY *gween.Tween
	interS *gween.Tween
}

func (c TextDisplay) Update() {
	if c.interR != nil {
		r, finR := c.interR.Update(timing.DT)
		c.Mask.R = uint8(r)
		if finR {
			c.interR = nil
		}
	}
	if c.interG != nil {
		g, finG := c.interG.Update(timing.DT)
		c.Mask.G = uint8(g)
		if finG {
			c.interG = nil
		}
	}
	if c.interB != nil {
		b, finB := c.interB.Update(timing.DT)
		c.Mask.B = uint8(b)
		if finB {
			c.interB = nil
		}
	}
	if c.interA != nil {
		a, finA := c.interA.Update(timing.DT)
		c.Mask.A = uint8(a)
		if finA {
			c.interA = nil
		}
	}
}