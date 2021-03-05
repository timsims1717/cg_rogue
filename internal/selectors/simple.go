package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type NullSelect struct {}

func NewNullSelect() *NullSelect {
	return &NullSelect{}
}

func (s *NullSelect) Init(_ *input.Input) {}

func (s *NullSelect) SetValues(_ ActionValues) {}

func (s *NullSelect) Update() {}

func (s *NullSelect) IsCancelled() bool {
	return false
}

func (s *NullSelect) IsDone() bool {
	return true
}

func (s *NullSelect) Finish() []world.Coords {
	return []world.Coords{}
}