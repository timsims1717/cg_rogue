package ui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"image/color"
)

type TextDisplay struct {
	Text  *text.Text
	Raw   string
	Show  bool
	UI    bool
	Align TextAlign

	Mat             pixel.Matrix
	Pos             pixel.Vec
	Rot             float64
	Scalar          pixel.Vec
	TransformEffect *animation.TransformEffect

	TextColor   color.RGBA
	ColorEffect *animation.ColorEffect
}

type TextAlign int

const (
	Left = iota
	Center
	Right
)

func (d *TextDisplay) Update() {
	if d.Show && len(d.Raw) > 0 {
		d.Text.Clear()
		if d.Align == Center {
			d.Text.Dot.X -= d.Text.BoundsOf(d.Raw).W() / 2
		} else if d.Align == Right {
			d.Text.Dot.X -= d.Text.BoundsOf(d.Raw).W()
		}
		d.Text.Color = d.TextColor
		fmt.Fprintf(d.Text, d.Raw)
		if d.UI {
			d.Mat = camera.UITransform(camera.Cam, d.Pos, d.Scalar, d.Rot, cfg.WindowWidthF, cfg.WindowHeightF)
		} else {
			d.Mat = pixel.IM.ScaledXY(pixel.ZV, d.Scalar).Rotated(pixel.ZV, d.Rot).Moved(d.Pos)
		}
		if d.TransformEffect != nil {
			d.TransformEffect.Update()
			if d.TransformEffect.IsDone() {
				d.TransformEffect = nil
			}
		}
		if d.ColorEffect != nil {
			d.ColorEffect.Update()
			if d.ColorEffect.IsDone() {
				d.ColorEffect = nil
			}
		}
	}
}

func (d *TextDisplay) Draw(target pixel.Target) {
	if d.Show && len(d.Raw) > 0 {
		d.Text.DrawColorMask(target, d.Mat, d.TextColor)
	}
}

func (d *TextDisplay) GetPos() pixel.Vec {
	return d.Pos
}

func (d *TextDisplay) SetPos(v pixel.Vec) {
	d.Pos = v
}

func (d *TextDisplay) GetRot() float64 {
	return d.Rot
}

func (d *TextDisplay) SetRot(r float64) {
	d.Rot = r
}

func (d *TextDisplay) GetScaled() pixel.Vec {
	return d.Scalar
}

func (d *TextDisplay) SetScaled(v pixel.Vec) {
	d.Scalar = v
}

func (d *TextDisplay) GetColor() color.RGBA {
	return d.TextColor
}

func (d *TextDisplay) SetColor(mask color.RGBA) {
	d.TextColor = mask
}