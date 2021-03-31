package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

// A Selector can be updated, checked for completion, cancelled, and finished.
type Selector interface {
	Update(*input.Input)
	SetValues(ActionValues)
}

type AbstractSelector struct {
	Selector Selector
	IsMove   bool
	area     []world.Coords
	origin   world.Coords
	isDone   bool
	cancel   bool
}

func (s *AbstractSelector) Reset(origin world.Coords) {
	s.origin = origin
	s.area = []world.Coords{}
	s.cancel = false
	s.isDone = false
}

func (s *AbstractSelector) IsDone() bool {
	return s.isDone
}

func (s *AbstractSelector) Finish() []world.Coords {
	return s.area
}

func (s *AbstractSelector) Cancel() {
	s.cancel = true
}

func (s *AbstractSelector) IsCancelled() bool {
	return s.cancel
}

type ActionValues struct {
	Source   *characters.Character
	Damage   int
	Move     int
	Range    int
	Targets  int
	Strength int
	Heal     int
	Checks   floor.PathChecks
}