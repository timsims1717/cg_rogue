package selectors

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type EmptyHexSelect struct {
	input      *input.Input
	clicked    []world.Coords
	count      int
	maxRange   int
	origin     world.Coords
	isDone     bool
	cancel     bool
}

func NewEmptyHexSelect() *EmptyHexSelect {
	return &EmptyHexSelect{}
}

func (s *EmptyHexSelect) Init(input *input.Input) {
	s.isDone = false
	s.cancel = false
	s.input = input
	s.clicked = []world.Coords{}
}

func (s *EmptyHexSelect) SetValues(values actions.ActionValues) {
	s.origin = values.Source.Coords
	s.maxRange = util.Max(values.Range, values.Move)
	s.count = values.Targets
}

func (s *EmptyHexSelect) Update() {
	if !s.isDone {
		x := s.input.Coords.X
		y := s.input.Coords.Y
		hex := floor.CurrentFloor.Get(s.input.Coords)
		legal := hex != nil && hex.Occupant == nil && !hex.Empty && world.DistanceHex(s.origin, s.input.Coords) <= s.maxRange
		if legal {
			if s.input.Select.JustPressed() {
				s.input.Select.Consume()
				// add to or remove from the clicked array
				found := -1
				for i, hex := range s.clicked {
					if hex.X == x && hex.Y == y {
						found = i
						break
					}
				}
				if found == -1 {
					// add to clicked array
					s.clicked = append(s.clicked, s.input.Coords)
				} else {
					// remove from clicked array
					s.clicked[len(s.clicked)-1], s.clicked[found] = s.clicked[found], s.clicked[len(s.clicked)-1]
					s.clicked = s.clicked[:len(s.clicked)-1]
				}
			}
		}
		for _, sel := range s.clicked {
			ui.AddSelectUI(ui.Move, sel.X, sel.Y)
		}
		if legal {
			ui.AddSelectUI(ui.MoveSolid, x, y)
		} else {
			ui.AddSelectUI(ui.Blank, x, y)
		}
	}
	if s.input.Cancel.JustPressed() {
		s.input.Cancel.Consume()
		// cancel
		s.cancel = true
	}
}

func (s *EmptyHexSelect) IsCancelled() bool {
	return s.cancel
}

func (s *EmptyHexSelect) IsDone() bool {
	return len(s.clicked) == s.count || s.isDone
}

func (s *EmptyHexSelect) Finish() []world.Coords {
	return s.clicked
}

type PathSelect struct {
	input      *input.Input
	picked     []world.Coords
	maxRange   int
	origin     world.Coords
	isDone     bool
	cancel     bool
	Unoccupied bool
	Nonempty   bool
	EndUnocc   bool
	EndNonemp  bool
}

func NewPathSelect() *PathSelect {
	return &PathSelect{}
}

func (s *PathSelect) Init(input *input.Input) {
	s.isDone = false
	s.cancel = false
	s.input = input
	s.picked = []world.Coords{}
}

func (s *PathSelect) SetValues(values actions.ActionValues) {
	s.origin = values.Source.Coords
	s.maxRange = util.Max(values.Range, values.Move)
}

func (s *PathSelect) Update() {
	if !s.isDone {
		x, y := s.input.Coords.X, s.input.Coords.Y
		hex := floor.CurrentFloor.Get(s.input.Coords)
		legal := hex != nil && (s.EndUnocc || hex.Occupant == nil) && (s.EndNonemp || !hex.Empty) && world.DistanceHex(s.origin, s.input.Coords) <= s.maxRange
		if legal {
			checks := floor.PathChecks{
				NotFilled:  true,
				Unoccupied: s.Unoccupied,
				NonEmpty:   s.Nonempty,
				Orig:       s.origin,
			}
			path, dist, found := floor.CurrentFloor.FindPath(s.origin, s.input.Coords, checks)
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
		if s.input.Cancel.JustPressed() {
			s.input.Cancel.Consume()
			// cancel
			s.cancel = true
		}
	}
}

func (s *PathSelect) IsCancelled() bool {
	return s.cancel
}

func (s *PathSelect) IsDone() bool {
	return len(s.picked) > 0
}

func (s *PathSelect) Finish() []world.Coords {
	return s.picked
}