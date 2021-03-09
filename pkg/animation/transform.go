package animation

import (
	"github.com/faiface/pixel"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
)

type Transformable interface {
	GetPos() pixel.Vec
	SetPos(pixel.Vec)
	GetRot() float64
	SetRot(float64)
	GetScaled() pixel.Vec
	SetScaled(pixel.Vec)
}

type TransformEffect struct {
	target  Transformable
	interX  *gween.Tween
	interY  *gween.Tween
	interR  *gween.Tween
	interSX *gween.Tween
	interSY *gween.Tween
	isDone  bool
}

func (e *TransformEffect) Update() {
	isDone := true
	pos := e.target.GetPos()
	rot := e.target.GetRot()
	sca := e.target.GetScaled()
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
	e.target.SetPos(pos)
	e.target.SetRot(rot)
	e.target.SetScaled(sca)
	e.isDone = isDone
}

func (e *TransformEffect) IsDone() bool {
	return e.isDone
}

type TransformBuilder struct {
	Target  Transformable
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