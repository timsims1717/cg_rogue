package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type TargetSelect struct {
	*AbstractSelector
	Count    int
	MaxRange int
}

func NewTargetSelect() *AbstractSelector {
	sel := &AbstractSelector{}
	target := &TargetSelect{
		sel,
		0,
		0,
	}
	sel.Selector = target
	return sel
}

func (s *TargetSelect) SetValues(values ActionValues) {
	s.Count = values.Targets
	s.MaxRange = values.Range
}

func (s *TargetSelect) Update(input *input.Input) {
	if !s.isDone {
		x := input.Coords.X
		y := input.Coords.Y
		inRange := world.DistanceSimple(s.origin, input.Coords) <= s.MaxRange
		occ := floor.CurrentFloor.GetOccupant(input.Coords)
		if occ != nil {
			if _, ok := occ.(objects.Targetable); ok {
				if !inRange {
					// todo: highlight
				} else if input.Select.JustPressed() {
					input.Select.Consume()
					s.area = append(s.area, input.Coords)
					s.isDone = true
				} else {
					AddSelectUI(Attack, x, y)
				}
			}
		} else {
			if inRange {
				AddSelectUI(Default, x, y)
			} else {
				AddSelectUI(Blank, x, y)
			}
		}
	}
}