package animation

import (
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"golang.org/x/image/colornames"
	"image/color"
)

type Colorable interface {
	GetColor() color.RGBA
	SetColor(color.RGBA)
}

type ColorEffect struct {
	target Colorable
	interR *gween.Tween
	interG *gween.Tween
	interB *gween.Tween
	interA *gween.Tween
	isDone bool
}

func (e *ColorEffect) Update() {
	isDone := true
	col := e.target.GetColor()
	if e.interR != nil {
		r, finR := e.interR.Update(timing.DT)
		col.R = uint8(r)
		if finR {
			e.interR = nil
		} else {
			isDone = false
		}
	}
	if e.interG != nil {
		g, finG := e.interG.Update(timing.DT)
		col.G = uint8(g)
		if finG {
			e.interG = nil
		} else {
			isDone = false
		}
	}
	if e.interB != nil {
		b, finB := e.interB.Update(timing.DT)
		col.B = uint8(b)
		if finB {
			e.interB = nil
		} else {
			isDone = false
		}
	}
	if e.interA != nil {
		a, finA := e.interA.Update(timing.DT)
		col.A = uint8(a)
		if finA {
			e.interA = nil
		} else {
			isDone = false
		}
	}
	e.target.SetColor(col)
	e.isDone = isDone
}

func (e *ColorEffect) IsDone() bool {
	return e.isDone
}

func FadeOut(target Colorable, dur float64) *ColorEffect {
	start := target.GetColor()
	end := colornames.Black
	end.A = 0
	return &ColorEffect{
		target: target,
		interR: gween.New(float64(start.R), float64(end.R), dur, ease.Linear),
		interG: gween.New(float64(start.G), float64(end.R), dur, ease.Linear),
		interB: gween.New(float64(start.B), float64(end.R), dur, ease.Linear),
		interA: gween.New(float64(start.A), float64(end.R), dur, ease.Linear),
		isDone: false,
	}
}

func FadeFrom(target Colorable, col color.RGBA, dur float64) *ColorEffect {
	start := col
	end := target.GetColor()
	return &ColorEffect{
		target: target,
		interR: gween.New(float64(start.R), float64(end.R), dur, ease.Linear),
		interG: gween.New(float64(start.G), float64(end.R), dur, ease.Linear),
		interB: gween.New(float64(start.B), float64(end.R), dur, ease.Linear),
		interA: gween.New(float64(start.A), float64(end.R), dur, ease.Linear),
		isDone: false,
	}
}

func FadeTo(target Colorable, col color.RGBA, dur float64) *ColorEffect {
	start := target.GetColor()
	end := col
	return &ColorEffect{
		target: target,
		interR: gween.New(float64(start.R), float64(end.R), dur, ease.Linear),
		interG: gween.New(float64(start.G), float64(end.R), dur, ease.Linear),
		interB: gween.New(float64(start.B), float64(end.R), dur, ease.Linear),
		interA: gween.New(float64(start.A), float64(end.R), dur, ease.Linear),
		isDone: false,
	}
}

func Reset(target Colorable, dur float64) *ColorEffect {
	return FadeTo(target, colornames.White, dur)
}