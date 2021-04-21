package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type Selector interface {
	Update(*input.Input)
	SetValues(ActionValues)
	SetAbstract(*AbstractSelector)
}

func NewSelector(sel Selector, isMove bool) *AbstractSelector {
	nSel := &AbstractSelector{
		Selector: sel,
		IsMove:   isMove,
	}
	sel.SetAbstract(nSel)
	return nSel
}

// A Selector and AbstractSelector uses input.Input from the player
// to gather a certain area of hexes.
type AbstractSelector struct {
	Selector Selector
	Effect   *AbstractSelectionEffect
	IsMove   bool
	results  []*Result
	area     []world.Coords
	origin   world.Coords
	source   *floor.Character
	isDone   bool
	cancel   bool
}

func (s *AbstractSelector) SetSource(character *floor.Character) {
	s.source = character
}

func (s *AbstractSelector) Reset(origin world.Coords) {
	s.origin = origin
	s.area = []world.Coords{}
	s.results = []*Result{}
	s.cancel = false
	s.isDone = false
}

func (s *AbstractSelector) IsDone() bool {
	return s.isDone
}

func (s *AbstractSelector) Finish() []*Result {
	return s.results
}

func (s *AbstractSelector) Cancel() {
	s.cancel = true
	s.area = []world.Coords{}
}

func (s *AbstractSelector) IsCancelled() bool {
	return s.cancel
}

type Result struct {
	Area   []world.Coords
	Effect *AbstractSelectionEffect
	IsMove bool
}

func NewResult(area []world.Coords, effect *AbstractSelectionEffect, isMove bool) *Result {
	if effect != nil {
		return &Result{
			Area:   area,
			Effect: effect,
			IsMove: isMove,
		}
	}
	return &Result{
		Area:   area,
		Effect: nil,
		IsMove: isMove,
	}
}

// ActionValues holds the values that can be passed from cards
// to Selectors, Actions, etc.
type ActionValues struct {
	Source   *floor.Character
	Damage   int
	Move     int
	Range    int
	Targets  int
	Strength int
	Heal     int
	Defense  int
	Area     []world.Coords
}
