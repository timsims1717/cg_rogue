package cfg

import "github.com/faiface/pixel"

const (
	Title = "Card Game Rogue!"

	// Tile Constants
	TileSize       = 32.0
	Scalar         = 4.0
	ScaledTileSize = TileSize * Scalar
)

var (
	OffsetVector = pixel.V(0.0, ScaledTileSize * 0.5)
)