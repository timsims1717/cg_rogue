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
	// todo: add VAlign

	Transform       *animation.Transform
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
	transform := animation.NewTransform(false)
	transform.Anchor = animation.Anchor{
		H: animation.Left,
		V: animation.Bottom,
	}
	return &ActionText{
		Text:      text.New(pixel.ZV, typeface.BasicAtlas),
		Raw:       raw,
		TextColor: color.RGBA{},
		Transform: transform,
		Align:     Left,
	}
}

func (t *ActionText) Update(r pixel.Rect) {
	t.Text.Clear()
	if t.Align == Center {
		t.Text.Dot.X -= t.Text.BoundsOf(t.Raw).W() / 2.
	} else if t.Align == Right {
		t.Text.Dot.X -= t.Text.BoundsOf(t.Raw).W()
	}
	t.Text.Color = t.TextColor
	fmt.Fprintf(t.Text, t.Raw)
	t.Transform.Rect = t.Text.Bounds()
	t.Transform.Update(r)
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
	t.Text.DrawColorMask(target, t.Transform.Mat, t.TextColor)
}

func (t *ActionText) GetColor() color.RGBA {
	return t.TextColor
}

func (t *ActionText) SetColor(mask color.RGBA) {
	t.TextColor = mask
}