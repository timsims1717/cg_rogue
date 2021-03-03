package world

import (
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"math"
)

var ScaledTileSize float64
var Origin = Coords{
	X: 0,
	Y: 0,
}

// Coords is a convenience struct for passing tile coordinates
type Coords struct {
	X int
	Y int
}

func (a *Coords) Equals(b Coords) bool {
	return a.X == b.X && a.Y == b.Y
}

func (a *Coords) Neighbors(w, h int) []Coords {
	neighbors := make([]Coords, 0)
	if a.Y > 0 {
		neighbors = append(neighbors, Coords{a.X, a.Y-1})
	}
	if a.X < w - 1 {
		neighbors = append(neighbors, Coords{a.X+1, a.Y})
	}
	if a.Y < h - 1 {
		neighbors = append(neighbors, Coords{a.X, a.Y+1})
	}
	if a.X > 0 {
		neighbors = append(neighbors, Coords{a.X-1, a.Y})
	}
	if a.X % 2 == 0 {
		if a.X < w - 1 && a.Y > 0 {
			neighbors = append(neighbors, Coords{a.X+1, a.Y-1})
		}
		if a.X > 0 && a.Y > 0 {
			neighbors = append(neighbors, Coords{a.X-1, a.Y-1})
		}
	} else {
		if a.X < w - 1 && a.Y < h - 1 {
			neighbors = append(neighbors, Coords{a.X+1, a.Y+1})
		}
		if a.X > 0 && a.Y < h - 1 {
			neighbors = append(neighbors, Coords{a.X-1, a.Y+1})
		}
	}
	return neighbors
}

func (a *Coords) PathFrom(path []Coords) []Coords {
	if len(path) < 2 {
		return []Coords{*a}
	}
	orig := path[0]
	if a.Equals(orig) {
		return path
	}
	var np []Coords
	dx := a.X - orig.X
	dy := a.Y - orig.Y
	for _, p := range path {
		if a.X % 2 == 0 && p.X % 2 == 0 && a.X % 2 != orig.X % 2 {
			n := Coords{
				X: p.X + dx,
				Y: p.Y + dy - 1,
			}
			np = append(np, n)
		} else if a.X % 2 != 0 && p.X % 2 != 0 && a.X % 2 != orig.X % 2 {
			n := Coords{
				X: p.X + dx,
				Y: p.Y + dy + 1,
			}
			np = append(np, n)
		} else {
			n := Coords{
				X: p.X + dx,
				Y: p.Y + dy,
			}
			np = append(np, n)
		}
	}
	return np
}

func NextHex(a, b Coords) Coords {
	var y int
	var x int
	if a.X == b.X {
		x = b.X
		if a.Y > b.Y {
			y = b.Y - 1
		} else {
			y = b.Y + 1
		}
	} else {
		if a.X < b.X {
			x = b.X + 1
		} else {
			x = b.X - 1
		}
		if (a.X % 2 == 0 && a.Y > b.Y) || (a.X % 2 != 0 && a.Y < b.Y) {
			y = b.Y
		} else if a.X % 2 == 0 {
			y = b.Y + 1
		} else {
			y = b.Y - 1
		}
	}
	return Coords{
		X: x,
		Y: y,
	}
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

func DistanceHex(a, b Coords) int {
	dist := 0
	x, y := a.X, a.Y
	for x != b.X {
		if x % 2 == 0 && y > b.Y {
			y -= 1
		} else if x % 2 != 0 && y < b.Y {
			y += 1
		}
		if x > b.X {
			x -= 1
		} else {
			x += 1
		}
		dist += 1
	}
	return dist + util.Abs(y - b.Y)
}

func DistanceWorldHex(ax, ay, bx, by int) float64 {
	axf, ayf := MapToWorldHex(ax, ay)
	bxf, byf := MapToWorldHex(bx, by)
	x := axf - bxf
	y := ayf - byf
	return x * x + y * y
}

func OrderByDist(orig Coords, ul []Coords) []Coords {
	ol := make([]Coords, 0)
	for len(ul) > 0 {
		near := 10000
		index := 0
		for i, c := range ul {
			dist := DistanceHex(orig, c)
			if dist < near {
				index = i
				near = dist
			}
		}
		ol = append(ol, ul[index])
		ul = append(ul[:index], ul[index+1:]...)
	}
	return ol
}