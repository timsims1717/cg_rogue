package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type SelfSelect struct {
	*AbstractSelector
	Effect     *AbstractSelectionEffect
}

func (s *SelfSelect) SetValues(values ActionValues) {
	if s.Effect != nil {
		s.Effect.SetValues(values)
		s.Effect.SetOrig(s.origin)
	}
}

func (s *SelfSelect) Update(_ *input.Input) {
	s.area = []world.Coords{s.source.GetCoords()}
	s.Effect.SetArea(s.area)
	s.results = []*Result{
		{
			Area:   s.area,
			Effect: s.Effect,
		},
	}
	s.isDone = true
}

func (s *SelfSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}
