package world

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"math"
)

func DistanceSimple(a, b Coords) int {
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

func Distance(a, b Coords) float64 {
	axf, ayf := MapToWorld(a.X, a.Y)
	bxf, byf := MapToWorld(b.X, b.Y)
	x := axf - bxf
	y := ayf - byf
	return math.Sqrt(x * x + y * y)
}

func DistanceWorld(a, b pixel.Vec) float64 {
	x := a.X - b.X
	y := a.Y - b.Y
	return math.Sqrt(x * x + y * y)
}

func OrderByDistSimple(orig Coords, ul []Coords) []Coords {
	ol := make([]Coords, 0)
	for len(ul) > 0 {
		near := 10000
		index := 0
		for i, c := range ul {
			dist := DistanceSimple(orig, c)
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

func OrderByDist(orig Coords, ul []Coords) []Coords {
	ol := make([]Coords, 0)
	for len(ul) > 0 {
		near := 10000.0
		index := 0
		for i, c := range ul {
			dist := Distance(orig, c)
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

func OrderByDistWorld(orig pixel.Vec, ul []Coords) []Coords {
	ol := make([]Coords, 0)
	for len(ul) > 0 {
		near := 10000.0
		index := 0
		for i, c := range ul {
			dist := DistanceWorld(orig, pixel.V(MapToWorld(c.X, c.Y)))
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