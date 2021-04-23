package selector

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type SelectionEffect interface {
	Update()
	Draw(pixel.Target)
	SetAbstract(*AbstractSelectionEffect)
}

func NewSelectionEffect(effect SelectionEffect, values ActionValues) *AbstractSelectionEffect {
	nEffect := &AbstractSelectionEffect{
		Effect: effect,
		values: values,
	}
	effect.SetAbstract(nEffect)
	return nEffect
}

type AbstractSelectionEffect struct {
	Effect SelectionEffect
	area   []world.Coords
	values ActionValues
	orig   world.Coords
}

func (e *AbstractSelectionEffect) SetOrig(orig world.Coords) {
	e.orig = orig
}

func (e *AbstractSelectionEffect) SetValues(values ActionValues) {
	e.values = values
}

func (e *AbstractSelectionEffect) SetArea(area []world.Coords) {
	e.area = area
}

var SelectionSet selectionSet

type selectionSet struct {
	nset    []SelectionEffect
}

func AddSelectionEffect(effect *AbstractSelectionEffect) {
	if effect != nil {
		SelectionSet.nset = append(SelectionSet.nset, effect.Effect)
	}
}

func (s *selectionSet) Update() {
	for _, sel := range s.nset {
		sel.Update()
	}
}

func (s *selectionSet) Draw(win *pixelgl.Window) {
	gfx.Clear()
	for _, sel := range s.nset {
		sel.Draw(gfx.Batch())
	}
	gfx.Draw(win)
	s.nset = []SelectionEffect{}
}