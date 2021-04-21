package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type HexAreaSplitSelect struct {
	*AbstractSelector
	secArea       [][]world.Coords
	count         int
	maxRange      int
	targetArea    []world.Coords
	PathChecks    floor.PathChecks
	SecPathChecks floor.PathChecks
	Effect        *AbstractSelectionEffect
	SecEffect     *AbstractSelectionEffect
}

func (s *HexAreaSplitSelect) SetValues(values ActionValues) {
	s.count = values.Targets
	if s.count < 1 {
		s.count = 1
	}
	if s.IsMove {
		s.maxRange = values.Move
	} else {
		s.maxRange = values.Range
	}
	s.targetArea = values.Area
	if s.Effect != nil {
		s.Effect.SetValues(values)
		s.Effect.SetOrig(s.origin)
	}
	if s.SecEffect != nil {
		s.SecEffect.SetValues(values)
		s.SecEffect.SetOrig(s.origin)
	}
}

func (s *HexAreaSplitSelect) Update(input *input.Input) {
	if !s.isDone {
		hex := floor.CurrentFloor.IsLegal(input.Coords, s.PathChecks)
		legal := hex != nil && world.DistanceSimple(s.origin, input.Coords) <= s.maxRange
		if legal {
			if input.Select.JustPressed() {
				input.Select.Consume()
				// add to or remove from area
				var removed int
				s.area, removed = world.Remove(input.Coords, s.area)
				if removed == -1 {
					s.area = append(s.area, input.Coords)
					area, _ := world.Remove(input.Coords, floor.CurrentFloor.IsSetLegal(input.Coords.PathFrom(s.targetArea), s.SecPathChecks))
					s.secArea = append(s.secArea, area)
				} else {
					s.secArea = append(s.secArea[:removed], s.secArea[removed+1:]...)
				}
			}
		}
		var secArea []world.Coords
		for _, area := range s.secArea {
			secArea = world.Combine(secArea, area)
		}
		if s.Effect != nil {
			if legal && !world.CoordsIn(input.Coords, s.area) {
				s.Effect.SetArea(append(s.area, input.Coords))
			} else {
				s.Effect.SetArea(s.area)
				eff := NewSelectionEffect(&HighlightEffect{}, s.Effect.values)
				eff.SetArea([]world.Coords{input.Coords})
				AddSelectionEffect(eff)
			}
			AddSelectionEffect(s.Effect)
		}
		if s.SecEffect != nil {
			if legal && !world.CoordsIn(input.Coords, s.area) {
				temp, _ := world.Remove(input.Coords, input.Coords.PathFrom(s.targetArea))
				s.SecEffect.SetArea(world.Combine(secArea, temp))
			} else {
				s.SecEffect.SetArea(secArea)
			}
			s.SecEffect.SetOrig(input.Coords)
			AddSelectionEffect(s.SecEffect)
		}
		if input.Cancel.JustPressed() {
			input.Cancel.Consume()
			var removed int
			s.area, removed = world.Remove(input.Coords, s.area)
			if removed == -1 {
				s.Cancel()
			}
		}
		if len(s.area) >= s.count {
			s.isDone = true
			s.results = []*Result{
				NewResult(s.area, s.Effect, s.IsMove),
				NewResult(secArea, s.SecEffect, false),
			}
			s.secArea = [][]world.Coords{}
		}
	}
}

func (s *HexAreaSplitSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}
