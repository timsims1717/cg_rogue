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

func (h *EmptyHexSelect) Init(input *input.Input) {
	h.isDone = false
	h.cancel = false
	h.input = input
	h.clicked = []world.Coords{}
}

func (h *EmptyHexSelect) SetValues(values actions.ActionValues) {
	h.origin = values.Source.Coords
	h.maxRange = util.Max(values.Range, values.Move)
	h.count = values.Targets
}

func (h *EmptyHexSelect) Update() {
	if !h.isDone {
		x := h.input.Coords.X
		y := h.input.Coords.Y
		hex := floor.CurrentFloor.Get(h.input.Coords)
		legal := hex != nil && hex.Occupant == nil && !hex.Empty && world.DistanceHex(h.origin.X, h.origin.Y, x, y) <= h.maxRange
		if legal {
			if h.input.Select.JustPressed() {
				h.input.Select.Consume()
				// add to or remove from the clicked array
				found := -1
				for i, hex := range h.clicked {
					if hex.X == x && hex.Y == y {
						found = i
						break
					}
				}
				if found == -1 {
					// add to clicked array
					h.clicked = append(h.clicked, h.input.Coords)
				} else {
					// remove from clicked array
					h.clicked[len(h.clicked)-1], h.clicked[found] = h.clicked[found], h.clicked[len(h.clicked)-1]
					h.clicked = h.clicked[:len(h.clicked)-1]
				}
			}
		}
		for _, sel := range h.clicked {
			ui.AddSelectUI(ui.Move, sel.X, sel.Y)
		}
		if legal {
			ui.AddSelectUI(ui.MoveSolid, x, y)
		} else {
			ui.AddSelectUI(ui.Blank, x, y)
		}
	}
	if h.input.Cancel.JustPressed() {
		h.input.Cancel.Consume()
		// cancel
		h.cancel = true
	}
}

func (h *EmptyHexSelect) IsCancelled() bool {
	return h.cancel
}

func (h *EmptyHexSelect) IsDone() bool {
	return len(h.clicked) == h.count || h.isDone
}

func (h *EmptyHexSelect) Finish() []world.Coords {
	return h.clicked
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

func (p *PathSelect) Init(input *input.Input) {
	p.isDone = false
	p.cancel = false
	p.input = input
	p.picked = []world.Coords{}
}

func (p *PathSelect) SetValues(values actions.ActionValues) {
	p.origin = values.Source.Coords
	p.maxRange = util.Max(values.Range, values.Move)
}

func (p *PathSelect) Update() {
	if !p.isDone {
		x, y := p.input.Coords.X, p.input.Coords.Y
		hex := floor.CurrentFloor.Get(p.input.Coords)
		legal := hex != nil && (p.EndUnocc || hex.Occupant == nil) && (p.EndNonemp || !hex.Empty) && world.DistanceHex(p.origin.X, p.origin.Y, x, y) <= p.maxRange
		if legal {
			path, dist, found := floor.CurrentFloor.FindPath(p.origin, p.input.Coords, p.Unoccupied, p.Nonempty)
			if found && dist <= p.maxRange {
				if p.input.Select.JustPressed() {
					p.input.Select.Consume()
					for _, h := range path {
						if h.X != p.origin.X || h.Y != p.origin.Y {
							p.picked = append(p.picked, h.GetCoords())
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
		if p.input.Cancel.JustPressed() {
			p.input.Cancel.Consume()
			// cancel
			p.cancel = true
		}
	}
}

func (p *PathSelect) IsCancelled() bool {
	return p.cancel
}

func (p *PathSelect) IsDone() bool {
	return len(p.picked) > 0
}

func (p *PathSelect) Finish() []world.Coords {
	return p.picked
}