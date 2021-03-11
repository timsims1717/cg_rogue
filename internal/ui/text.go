package ui

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/typeface"
	"image/color"
)

type ActionText struct {
	Text   *text.Text
	Raw    string
	Align  TextAlign
	VAlign TextAlign

	Mat             pixel.Matrix
	Pos             pixel.Vec
	pos             pixel.Vec
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

func NewActionText(raw string) *ActionText {
	return &ActionText{
		Text:      text.New(pixel.ZV, typeface.BasicAtlas),
		Raw:       raw,
		TextColor: color.RGBA{},
		Scalar:    pixel.V(1., 1.),
		Align:     Left,
		VAlign:    Left,
	}
}

func (t *ActionText) Update() {
	t.Text.Clear()
	if t.Align == Center {
		t.Text.Dot.X -= t.Text.BoundsOf(t.Raw).W() / 2
	} else if t.Align == Right {
		t.Text.Dot.X -= t.Text.BoundsOf(t.Raw).W()
	}
	t.Text.Color = t.TextColor
	fmt.Fprintf(t.Text, t.Raw)
	t.Mat = pixel.IM.ScaledXY(pixel.ZV, t.Scalar).Rotated(pixel.ZV, t.Rot).Moved(t.pos)
	if t.VAlign == Center {
		t.Mat = t.Mat.Moved(pixel.V(0., (t.Text.Orig.Y - t.Text.Dot.Y) / 2.))
	} else if t.VAlign == Left {
		t.Mat = t.Mat.Moved(pixel.V(0., t.Text.Orig.Y - t.Text.Dot.Y))
	}
	if t.TransformEffect != nil {
		t.TransformEffect.Update()
		if t.TransformEffect.IsDone() {
			t.TransformEffect = nil
		}
	}
	if t.ColorEffect != nil {
		t.ColorEffect.Update()
		if t.ColorEffect.IsDone() {
			t.ColorEffect = nil
		}
	}
}

func (t *ActionText) Draw(target pixel.Target) {
	t.Text.DrawColorMask(target, t.Mat, t.TextColor)
}

func (t *ActionText) GetPos() pixel.Vec {
	return t.Pos
}

func (t *ActionText) SetPos(v pixel.Vec) {
	t.Pos = v
}

func (t *ActionText) GetRot() float64 {
	return t.Rot
}

func (t *ActionText) SetRot(r float64) {
	t.Rot = r
}

func (t *ActionText) GetScaled() pixel.Vec {
	return t.Scalar
}

func (t *ActionText) SetScaled(v pixel.Vec) {
	t.Scalar = v
}

func (t *ActionText) GetColor() color.RGBA {
	return t.TextColor
}

func (t *ActionText) SetColor(mask color.RGBA) {
	t.TextColor = mask
}