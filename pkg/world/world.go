package world

import (
	"math"
)

func MapToWorld(x, y int) (float64, float64) {
	if x % 2 != 0 {
		return (float64(x) + 0.5) * ScaledTileSize, (float64(y) + 1.0) * ScaledTileSize
	} else {
		return (float64(x) + 0.5) * ScaledTileSize, (float64(y) + 0.5) * ScaledTileSize
	}
}

func WorldToMap(x, y float64) (int, int) {
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