package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"golang.org/x/image/colornames"
	"image/color"
)

//var DebugDraw = imdraw.New(nil)

type ActionEl struct {
	T    *ActionText
	Show bool
	UI   bool

	Mat             pixel.Matrix
	Pos             pixel.Vec
	Rot             float64
	Scalar          pixel.Vec
	Spr             *pixel.Sprite
	canvas          *pixelgl.Canvas
	TransformEffect *animation.TransformEffect

	Mask        color.RGBA
	ColorEffect *animation.ColorEffect

	hovered   bool
	clicked   bool
	hoverFn   func()
	onHoverFn func()
	unHoverFn func()
	clickFn   func()
	unClickFn func()
}

func NewActionEl(t *ActionText, rect pixel.Rect) *ActionEl {
	return &ActionEl{
		T: t,
		canvas: pixelgl.NewCanvas(rect),
		Scalar: pixel.V(1., 1.),
		Mask:   colornames.White,
	}
}

func (e *ActionEl) Update(input *input.Input) {
	//DebugDraw.Clear()
	if e.Show {
		mouseOver := e.PointInside(input.World)
		if mouseOver && !e.hovered && e.onHoverFn != nil {
			e.onHoverFn()
		} else if !mouseOver && e.hovered && e.unHoverFn != nil {
			e.unHoverFn()
		} else if e.hovered && e.hoverFn != nil {
			e.hoverFn()
		}
		e.hovered = mouseOver
		if e.hovered && input.Select.JustPressed() && e.clickFn != nil {
			e.clickFn()
		} else if e.hovered && input.Select.JustReleased() && e.unClickFn != nil {
			e.unClickFn()
		}
		if e.T != nil {
			e.T.pos = e.T.Pos
			if e.T.Align == Center {
				e.T.pos.X += e.canvas.Bounds().W() / 2.
			} else if e.T.Align == Right {
				e.T.pos.X += e.canvas.Bounds().W()
			}
			if e.T.VAlign == Center {
				e.T.pos.Y += e.canvas.Bounds().H() / 2.
			} else if e.T.VAlign == Left {
				e.T.pos.Y += e.canvas.Bounds().H()
			}
			e.T.Update()
		}
		if e.UI {
			e.Mat = camera.UITransform(camera.Cam, e.Pos, e.Scalar, e.Rot, cfg.WindowWidthF, cfg.WindowHeightF)
		} else {
			e.Mat = pixel.IM.ScaledXY(pixel.ZV, e.Scalar).Rotated(pixel.ZV, e.Rot).Moved(e.Pos)
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
		//DebugDraw.Color = colornames.Red
		//DebugDraw.EndShape = imdraw.NoEndShape
		//for _, v := range e.canvas.Bounds().Vertices() {
		//	DebugDraw.Push(e.Mat.Project(v))
		//}
		//DebugDraw.Polygon(5.0)
	}
}

func (e *ActionEl) Draw(target pixel.Target) {
	if e.Show {
		e.canvas.Clear(color.RGBA{})
		if e.Spr != nil {
			r := e.Spr.Frame()
			e.Spr.Draw(e.canvas, pixel.IM.Moved(pixel.V(r.W()*0.5, r.H()*0.5)))
		}
		if e.T != nil {
			e.T.Draw(e.canvas)
		}
		//e.canvas.DrawColorMask(target, e.Mat, e.Mask)
		e.canvas.Draw(target, e.Mat)
		//DebugDraw.Draw(target)
	}
}

func (e *ActionEl) GetPos() pixel.Vec {
	return e.Pos
}

func (e *ActionEl) SetPos(v pixel.Vec) {
	e.Pos = v
}

func (e *ActionEl) GetRot() float64 {
	return e.Rot
}

func (e *ActionEl) SetRot(r float64) {
	e.Rot = r
}

func (e *ActionEl) GetScaled() pixel.Vec {
	return e.Scalar
}

func (e *ActionEl) SetScaled(v pixel.Vec) {
	e.Scalar = v
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

func (e *ActionEl) PointInside(point pixel.Vec) bool {
	return util.PointInside(point, e.canvas.Bounds(), e.Mat)
}