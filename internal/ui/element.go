package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"golang.org/x/image/colornames"
	"image/color"
)

//var DebugDraw = imdraw.New(nil)

// todo: add alignment and vertical alignment
type ActionEl struct {
	Text *ActionText
	Show bool
	IsUI bool

	Transform       *animation.Transform
	Spr             *pixel.Sprite
	canvas          *pixelgl.Canvas
	TransformEffect *animation.TransformEffect

	Mask        color.RGBA
	ColorEffect *animation.ColorEffect

	hovered      bool
	clicked      bool
	Disabled     bool
	disabled     bool
	hoverFn      func()
	onHoverFn    func()
	unHoverFn    func()
	clickFn      func()
	unClickFn    func()
	onDisabledFn func()
	disabledFn   func()
	onEnabledFn  func()
}

func NewActionEl(t *ActionText, rect pixel.Rect, isUI bool) *ActionEl {
	transform := animation.NewTransform(true)
	transform.Anchor = animation.Anchor{
		H: animation.Left,
		V: animation.Bottom,
	}
	transform.Rect = rect
	return &ActionEl{
		Text:      t,
		Transform: transform,
		canvas:    pixelgl.NewCanvas(rect),
		Mask:      colornames.White,
		IsUI:      isUI,
	}
}

func (e *ActionEl) Update(input *input.Input) {
	//DebugDraw.Clear()
	if e.Disabled {
		if e.disabled && e.disabledFn != nil {
			e.disabledFn()
		} else if !e.disabled && e.onDisabledFn != nil {
			e.onDisabledFn()
		}
		e.disabled = true
		if e.hovered {
			e.unHoverFn()
		}
		e.hovered = false
	} else {
		if e.disabled && e.onEnabledFn != nil {
			e.onEnabledFn()
		}
		e.disabled = false
		mouseOver := e.PointInside(input.World)
		if mouseOver && !e.hovered && e.onHoverFn != nil {
			e.onHoverFn()
		} else if !mouseOver && e.hovered && e.unHoverFn != nil {
			e.unHoverFn()
		} else if e.hovered && e.hoverFn != nil {
			e.hoverFn()
		}
		e.hovered = mouseOver
		if e.clicked || (e.hovered && input.Select.JustPressed() && e.clickFn != nil) {
			e.clickFn()
			e.clicked = false
		} else if e.hovered && input.Select.JustReleased() && e.unClickFn != nil {
			e.unClickFn()
			e.clicked = false
		}
	}
	if e.Text != nil {
		e.Text.Update(e.canvas.Bounds())
	}
	if e.TransformEffect != nil {
		e.TransformEffect.Update()
		if e.TransformEffect.IsDone() {
			e.TransformEffect = nil
		}
	}
	if e.ColorEffect != nil {
		e.ColorEffect.Update()
		if e.ColorEffect.IsDone() {
			e.ColorEffect = nil
		}
	}
	e.Transform.Rect = e.canvas.Bounds()
	e.Transform.Update(pixel.Rect{})
	if e.IsUI {
		e.Transform.Mat = camera.Cam.UITransform(e.Transform.RPos, e.Transform.Scalar, e.Transform.Rot)
	}
	//DebugDraw.Color = colornames.Red
	//DebugDraw.EndShape = imdraw.NoEndShape
	//for _, v := range e.canvas.Bounds().Vertices() {
	//	DebugDraw.Push(e.Mat.Project(v))
	//}
	//DebugDraw.Polygon(5.0)
}

func (e *ActionEl) Draw(target pixel.Target) {
	if e.Show {
		e.canvas.Clear(color.RGBA{})
		if e.Spr != nil {
			r := e.Spr.Frame()
			e.Spr.Draw(e.canvas, pixel.IM.Moved(pixel.V(r.W()*0.5, r.H()*0.5)))
		}
		if e.Text != nil {
			e.Text.Draw(e.canvas)
		}
		//e.canvas.DrawColorMask(target, e.Mat, e.Mask)
		e.canvas.Draw(target, e.Transform.Mat)
		//DebugDraw.Draw(target)
	}
}

func (e *ActionEl) GetColor() color.RGBA {
	return e.Mask
}

func (e *ActionEl) SetColor(mask color.RGBA) {
	e.Mask = mask
}

func (e *ActionEl) IsHovered() bool {
	return e.hovered
}

func (e *ActionEl) OnHover() {
	e.onHoverFn()
}

func (e *ActionEl) Hover() {
	e.hoverFn()
}

func (e *ActionEl) OnUnHover() {
	e.unHoverFn()
}

func (e *ActionEl) Click() {
	e.clicked = true
}

func (e *ActionEl) IsClicked() bool {
	return e.clicked
}

func (e *ActionEl) OnClick() {
	e.clickFn()
}

func (e *ActionEl) OnUnClick() {
	e.unClickFn()
}

func (e *ActionEl) SetHoverFn(fn func()) {
	e.hoverFn = fn
}

func (e *ActionEl) SetOnHoverFn(fn func()) {
	e.onHoverFn = fn
}

func (e *ActionEl) SetUnHoverFn(fn func()) {
	e.unHoverFn = fn
}

func (e *ActionEl) SetClickFn(fn func()) {
	e.clickFn = fn
}

func (e *ActionEl) SetUnClickFn(fn func()) {
	e.unClickFn = fn
}

func (e *ActionEl) SetOnDisabledFn(fn func()) {
	e.onDisabledFn = fn
}

func (e *ActionEl) SetDisabledFn(fn func()) {
	e.disabledFn = fn
}

func (e *ActionEl) SetEnabledFn(fn func()) {
	e.onEnabledFn = fn
}

func (e *ActionEl) PointInside(point pixel.Vec) bool {
	return util.PointInside(point, e.canvas.Bounds(), e.Transform.Mat)
}