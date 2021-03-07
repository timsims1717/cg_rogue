package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type PathSelect struct {
	input      *input.Input
	picked     []world.Coords
	maxRange   int
	origin     world.Coords
	isDone     bool
	pathChecks floor.PathChecks
}

func NewPathSelect() *PathSelect {
	return &PathSelect{}
}

func (s *PathSelect) Init(input *input.Input) {
	s.isDone = false
	s.input = input
	s.picked = []world.Coords{}
}

func (s *PathSelect) SetValues(values ActionValues) {
	s.origin = values.Source.Coords
	s.maxRange = util.Max(values.Range, values.Move)
	s.pathChecks = values.Checks
}

func (s *PathSelect) Update() {
	if !s.isDone {
		x, y := s.input.Coords.X, s.input.Coords.Y
		s.pathChecks.Orig = s.origin
		hex := floor.CurrentFloor.IsLegal(s.input.Coords, s.pathChecks)
		legal := hex != nil && world.DistanceSimple(s.origin, s.input.Coords) <= s.maxRange
		if legal {
			path, dist, found := floor.CurrentFloor.FindPath(s.origin, s.input.Coords, s.pathChecks)
			if found && dist <= s.maxRange {
				if s.input.Select.JustPressed() {
					s.input.Select.Consume()
					for _, h := range path {
						if h.X != s.origin.X || h.Y != s.origin.Y {
							s.picked = append(s.picked, h)
						}
					}
				}
				for _, h := range path {
					if h.X == x && h.Y == y {
						ui.AddSelectUI(ui.MoveSolid, h.X, h.Y)
					} else {
						ui.AddSelectUI(ui.Move, h.X, h.Y)
					}
				}
			}
		}
	}
}

func (s *PathSelect) IsDone() bool {
	return len(s.picked) > 0
}

func (s *PathSelect) Finish() []world.Coords {
	return s.picked
}