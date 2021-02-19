package world

import (
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"math"
)

var ScaledTileSize float64

// Coords is a convenience struct for passing tile coordinates
type Coords struct {
	X int
	Y int
}

func (a *Coords) Equals(b Coords) bool {
	return a.X == b.X && a.Y == b.Y
}

func MapToWorldHex(x, y int) (float64, float64) {
	if x % 2 != 0 {
		return (float64(x) + 0.5) * ScaledTileSize, (float64(y) + 1.0) * ScaledTileSize
	} else {
		return (float64(x) + 0.5) * ScaledTileSize, (float64(y) + 0.5) * ScaledTileSize
	}
}

func WorldToMapSquare(x, y float64) (int, int) {
	mapX := math.Floor(x / ScaledTileSize)
	if int(mapX) % 2 != 0 {
		return int(mapX), int(math.Floor(y / ScaledTileSize - 0.5))
	} else {
		return int(mapX), int(math.Floor(y / ScaledTileSize))
	}
}

func WorldToMapHex(x, y float64) (int, int) {
	mapX := x / ScaledTileSize
	mapXf := math.Floor(mapX)
	x1 := mapXf - 0.125
	x2 := mapXf + 0.125
	x3 := mapXf + 0.875
	x4 := mapXf + 1.125
	mapY := y / ScaledTileSize
	mapYO := math.Floor(mapY - 0.5)
	mapYE := math.Floor(mapY)
	mapYf := mapYE
	odd := int(mapXf) % 2 != 0
	if odd {
		mapYf = mapYO
	}
	_ = mapYf - 0.5
	y2 := mapYf + 0.5
	if odd {
		y2 += 0.5
	}
	_ = mapYf + 1.0
	if mapX < x2 {
		// In the left section of the hex
		if mapY > y2 {
			// In the top half of the hex
			s := 0.25*(mapY-y2) - 0.5*(mapX-x1)
			if s > 0 {
				if odd {
					return int(mapXf - 1.0), int(mapYf + 1.0)
				} else {
					return int(mapXf - 1.0), int(mapYf)
				}
			}
		} else {
			// In the bottom half of the hex
			s := 0.25*(mapY-y2) + 0.5*(mapX-x1)
			if s < 0 {
				if odd {
					return int(mapXf - 1.0), int(mapYf)
				} else {
					return int(mapXf - 1.0), int(mapYf - 1.0)
				}
			}
		}
	} else if mapX > x3 {
		// In the right section of the hex
		if mapY > y2 {
			// In the top half of the hex
			s := 0.25*(mapY-y2) + 0.5*(mapX-x4)
			if s > 0 {
				if odd {
					return int(mapXf + 1.0), int(mapYf + 1.0)
				} else {
					return int(mapXf + 1.0), int(mapYf)
				}
			}
		} else {
			// In the bottom half of the hex
			s := 0.25*(mapY-y2) - 0.5*(mapX-x4)
			if s < 0 {
				if odd {
					return int(mapXf + 1.0), int(mapYf)
				} else {
					return int(mapXf + 1.0), int(mapYf - 1.0)
				}
			}
		}
	}
	// Normal
	return int(mapXf), int(mapYf)
}

func DistanceHex(ax, ay, bx, by int) int {
	dist := 0
	x, y := ax, ay
	for x != bx {
		if x % 2 == 0 && y > by {
			y -= 1
		} else if x % 2 != 0 && y < by {
			y += 1
		}
		if x > bx {
			x -= 1
		} else {
			x += 1
		}
		dist += 1
	}
	return dist + util.Abs(y - by)
}

func DistanceWorldHex(ax, ay, bx, by int) float64 {
	axf, ayf := MapToWorldHex(ax, ay)
	bxf, byf := MapToWorldHex(bx, by)
	x := axf - bxf
	y := ayf - byf
	return x * x + y * y
}