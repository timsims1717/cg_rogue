package floor

import (
	"github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
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
	Empty    bool
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
	}
}

func (h *Hex) GetCoords() world.Coords {
	return world.Coords{
		X: h.X,
		Y: h.Y,
	}
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