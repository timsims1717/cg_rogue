package world

import (
	"github.com/faiface/pixel"
	"math"
)

func MapToWorld(a Coords) pixel.Vec {
	if a.X%2 != 0 {
		return pixel.V((float64(a.X)+0.5)*TileSize, (float64(a.Y)+1.0)*TileSize)
	} else {
		return pixel.V((float64(a.X)+0.5)*TileSize, (float64(a.Y)+0.5)*TileSize)
	}
}

func WorldToMap(x, y float64) (int, int) {
	mapX := x / TileSize
	mapXf := math.Floor(mapX)
	x1 := mapXf - 0.125
	x2 := mapXf + 0.125
	x3 := mapXf + 0.875
	x4 := mapXf + 1.125
	mapY := y / TileSize
	mapYO := math.Floor(mapY - 0.5)
	mapYE := math.Floor(mapY)
	mapYf := mapYE
	odd := int(mapXf)%2 != 0
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
