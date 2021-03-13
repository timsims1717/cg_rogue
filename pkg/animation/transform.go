package animation

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
)

//type Transformable interface {
//	GetPos() pixel.Vec
//	SetPos(pixel.Vec)
//	GetRot() float64
//	SetRot(float64)
//	GetScaled() pixel.Vec
//	SetScaled(pixel.Vec)
//}

type Alignment int

const (
	Left = iota
	Center
	Right
	Top = Left
	Bottom = Right
)

type Anchor struct {
	H Alignment
	V Alignment
}

type Transform struct {
	Cam    *camera.Camera
	Anchor Anchor
	Rect   pixel.Rect
	Mat    pixel.Matrix
	Pos    pixel.Vec
	Offset pixel.Vec
	pos    pixel.Vec
	Rot    float64
	Scalar pixel.Vec
	ocent  bool
}

func NewTransform(isOrigCentered bool) *Transform {
	return &Transform{
		Scalar: pixel.Vec{
			X: 1.,
			Y: 1.,
		},
		ocent: isOrigCentered,
	}
}

func (t *Transform) Update(r pixel.Rect) {
	t.pos = t.Pos
	if t.ocent {
		if t.Anchor.H == Left {
			t.pos.X += t.Rect.W() * t.Scalar.X / 2.
		} else if t.Anchor.H == Center {
			t.pos.X += r.W() / 2.
		} else if t.Anchor.H == Right {
			t.pos.X += r.W()
			t.pos.X -= t.Rect.W() * t.Scalar.X / 2.
		}
		if t.Anchor.V == Bottom {
			t.pos.Y += t.Rect.H() * t.Scalar.Y / 2.
		} else if t.Anchor.V == Center {
			t.pos.Y += r.H() / 2.
		} else if t.Anchor.V == Top {
			t.pos.Y += r.H()
			t.pos.Y -= t.Rect.H() * t.Scalar.Y / 2.
		}
	} else {
		if t.Anchor.H == Center {
			t.pos.X += r.W() / 2.
		} else if t.Anchor.H == Right {
			t.pos.X += r.W()
		}
		if t.Anchor.V == Center {
			t.pos.Y += r.H() / 2.
		} else if t.Anchor.V == Top {
			t.pos.Y += r.H()
		}
	}
	//if t.Anchor.V == Bottom {
	//	t.pos.Y += t.Rect.H() / 2.
	//} else if t.Anchor.V == Top {
	//	t.pos.Y -= t.Rect.H() / 2.
	//}
	t.pos.X += t.Offset.X
	t.pos.Y += t.Offset.Y
	if t.Cam != nil {
		t.Mat = t.Cam.UITransform(t.pos, t.Scalar, t.Rot)
	} else {
		t.Mat = pixel.IM
		t.Mat = t.Mat.ScaledXY(pixel.ZV, t.Scalar)
		t.Mat = t.Mat.Rotated(pixel.ZV, t.Rot)
		t.Mat = t.Mat.Moved(t.pos)
	}
}

type TransformEffect struct {
	target  *Transform
	interX  *gween.Tween
	interY  *gween.Tween
	interR  *gween.Tween
	interSX *gween.Tween
	interSY *gween.Tween
	isDone  bool
}

func (e *TransformEffect) Update() {
	isDone := true
	pos := e.target.Pos
	rot := e.target.Rot
	sca := e.target.Scalar
	if e.interX != nil {
		x, fin := e.interX.Update(timing.DT)
		pos.X = x
		if fin {
			e.interX = nil
		} else {
			isDone = false
		}
	}
	if e.interY != nil {
		y, fin := e.interY.Update(timing.DT)
		pos.Y = y
		if fin {
			e.interY = nil
		} else {
			isDone = false
		}
	}
	if e.interR != nil {
		r, fin := e.interR.Update(timing.DT)
		rot = r
		if fin {
			e.interR = nil
		} else {
			isDone = false
		}
	}
	if e.interSX != nil {
		x, fin := e.interSX.Update(timing.DT)
		sca.X = x
		if fin {
			e.interSX = nil
		} else {
			isDone = false
		}
	}
	if e.interSY != nil {
		y, fin := e.interSY.Update(timing.DT)
		sca.Y = y
		if fin {
			e.interSY = nil
		} else {
			isDone = false
		}
	}
	e.target.Pos = pos
	e.target.Rot = rot
	e.target.Scalar = sca
	e.isDone = isDone
}

func (e *TransformEffect) IsDone() bool {
	return e.isDone
}

type TransformBuilder struct {
	Target  *Transform
	InterX  *gween.Tween
	InterY  *gween.Tween
	InterR  *gween.Tween
	InterSX *gween.Tween
	InterSY *gween.Tween
}

func (b *TransformBuilder) Build() *TransformEffect {
	return &TransformEffect{
		target:  b.Target,
		interX:  b.InterX,
		interY:  b.InterY,
		interR:  b.InterR,
		interSX: b.InterSX,
		interSY: b.InterSY,
	}
}