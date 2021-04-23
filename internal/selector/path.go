package selector

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type PathSelect struct {
	*AbstractSelector
	MaxRange   int
	PathChecks floor.PathChecks
	Effect     *AbstractSelectionEffect
	selecting  bool
}

func (s *PathSelect) SetValues(values ActionValues) {
	s.MaxRange = values.Move
	s.PathChecks.Orig = s.origin
	if s.Effect != nil {
		s.Effect.SetValues(values)
		s.Effect.SetOrig(s.origin)
	}
}

func (s *PathSelect) Update(input *input.Input) {
	if !s.isDone {
		if s.IsMove {
			s.source.RemoveClaim()
		}
		if !s.selecting {
			//x, y := input.Coords.X, input.Coords.Y
			s.area = []world.Coords{}
			s.PathChecks.Orig = s.origin
			hex := floor.CurrentFloor.IsLegal(input.Coords, s.PathChecks)
			legal := hex != nil && world.DistanceSimple(s.origin, input.Coords) <= s.MaxRange
			if legal {
				path, dist, found := floor.CurrentFloor.FindPath(s.origin, input.Coords, s.PathChecks)
				if found && dist <= s.MaxRange {
					for _, h := range path {
						s.area = append(s.area, h)
					}
					if input.Select.JustPressed() {
						input.Select.Consume()
						s.selecting = true
					}
				}
			}
			if input.Cancel.JustPressed() {
				input.Cancel.Consume()
				s.Cancel()
			}
		} else {
			if input.Cancel.JustPressed() {
				input.Cancel.Consume()
				s.selecting = false
			} else if input.Select.JustReleased() {
				if len(s.area) > 0 {
					s.isDone = true
					s.results = []*Result{
						NewResult(s.area, s.Effect, s.IsMove),
					}
				}
				s.selecting = false
			} else {
				if len(s.area) > 0 {
					curr := s.area[len(s.area)-1]
					if curr != input.Coords {
						if len(s.area) > 1 && s.area[len(s.area)-2] == input.Coords {
							s.area, _ = world.Remove(s.area[len(s.area)-1], s.area)
						} else {
							path, _, found := floor.CurrentFloor.FindPath(curr, input.Coords, s.PathChecks)
							if found {
								s.area = world.Combine(s.area, path)
							}
						}
					}
				}
				s.area = floor.CurrentFloor.LongestLegalPath(s.area, s.MaxRange, s.PathChecks)
			}
		}
		if s.IsMove && len(s.area) > 0 {
			s.source.MakeClaim(s.area[len(s.area)-1])
		}
		if s.Effect != nil && !input.Cancel.Pressed() {
			s.Effect.SetArea(s.area)
			AddSelectionEffect(s.Effect)
		}
	}
}

func (s *PathSelect) SetAbstract(sel *AbstractSelector) {
	s.AbstractSelector = sel
}
