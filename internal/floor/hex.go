package floor

import (
	"github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

// Hex represents a single hexagonal tile on a Floor
// X and Y are its coordinates on the Floor
type Hex struct {
	X        int
	Y        int
	f        *Floor
	Tile     *pixel.Sprite
	Occupant *Character
	Claimant *Character
	Empty    bool
	storage  []*Character
	effects  map[int][]img.Sprite
	matrix   pixel.Matrix
}

// NewHex creates a Hex with a reference to its Floor
func NewHex(floor *Floor, x, y int, tile *pixel.Sprite) Hex {
	return Hex{
		X:        x,
		Y:        y,
		f:        floor,
		Tile:     tile,
		Occupant: nil,
		Empty:    false,
		storage:  []*Character{},
		effects:  make(map[int][]img.Sprite),
		matrix:   pixel.IM.Moved(world.MapToWorld(world.Coords{ X: x, Y: y })),
	}
}

func (h *Hex) GetCoords() world.Coords {
	return world.Coords{
		X: h.X,
		Y: h.Y,
	}
}

func (h *Hex) GetWorlds() pixel.Vec {
	return world.MapToWorld(h.GetCoords())
}

func (h *Hex) IsClaimed() bool {
	return h == nil || h.Claimant != nil || h.Occupant != nil
}

func (h *Hex) GetClaimant() *Character {
	if h == nil {
		return nil
	}
	return h.Claimant
}

func (h *Hex) RemoveClaim() {
	if h != nil {
		h.Claimant = nil
	}
}

func (h *Hex) Claim(c *Character) {
	if h != nil {
		if h.Claimant != nil {
			h.Claimant.RemoveClaim()
		}
		h.Claimant = c
		c.Claim = h.GetCoords()
	}
}

func (h *Hex) IsOccupied() bool {
	return h == nil || h.Occupant != nil
}

func (h *Hex) GetOccupant() *Character {
	if h == nil {
		return nil
	}
	return h.Occupant
}

func (h *Hex) RemoveOccupant() *Character {
	if h != nil && h.Occupant != nil {
		former := h.Occupant
		h.Occupant = nil
		return former
	}
	return nil
}

func (h *Hex) IsStoring(c *Character) bool {
	if h != nil {
		for _, s := range h.storage {
			if s.ID() == c.ID() {
				return true
			}
		}
	}
	return false
}

func (h *Hex) Store(c *Character) {
	if h != nil {
		if c.Floor != nil {
			hex := c.Floor.Get(c.GetCoords())
			if !h.IsStoring(c) {
				hex.RemoveOccupant()
				hex.UnStore(c)
				c.Floor.UnStore(c)
			}
		}
		c.Floor = h.f
		h.storage = append(h.storage, c)
		c.OnMap = true
	}
}

func (h *Hex) UnStore(c *Character) {
	if h != nil {
		for i, s := range h.storage {
			if s.ID() == c.ID() {
				h.storage = append(h.storage[:i], h.storage[i+1:]...)
				break
			}
		}
	}
}

func (h *Hex) Clear() {
	h.RemoveOccupant()
	h.RemoveClaim()
	h.storage = []*Character{}
}

func (h *Hex) AddEffect(spr []img.Sprite, p int) {
	if h == nil {
		return
	}
	if _, ok := h.effects[p]; !ok {
		h.effects[p] = spr
	}
}

func (h *Hex) ClearEffects() {
	h.effects = make(map[int][]img.Sprite)
}

// PathNeighbors is part of the astar implementation and returns
// legal moves to the Hex' neighbors.
func (h *Hex) PathNeighbors() []astar.Pather {
	n := h.f.Neighbors(h)
	var neighbors []astar.Pather
	for _, hex := range n {
		neighbors = append(neighbors, hex)
	}
	return neighbors
}

func (h *Hex) PathNeighborCost(to astar.Pather) float64 {
	return 1.0
}

func (h *Hex) PathEstimatedCost(to astar.Pather) float64 {
	v := world.MapToWorld(h.GetCoords())
	return pixel.L(h.f.PathLine.Closest(v), v).Len()
}
