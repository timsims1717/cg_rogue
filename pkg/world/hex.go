package world

import (
	"github.com/faiface/pixel"
	"math"
	"math/rand"
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

func (a Coords) Equals(b Coords) bool {
	return a.X == b.X && a.Y == b.Y
}

func (a Coords) Neighbors(w, h int) []Coords {
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

func (a Coords) PathFrom(path []Coords) []Coords {
	if len(path) < 2 {
		return []Coords{a}
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

func (a Coords) InSextant(b Coords, s Sextant) bool {
	if sextant := GetSextant(b, a); s == sextant {
		return true
	}
	if sextant := GetSextantBias(b, a, true); s == sextant {
		return true
	}
	if sextant := GetSextantBias(b, a, false); s == sextant {
		return true
	}
	return false
}

func NextHexLine(orig, next Coords) Coords {
	var y int
	var x int
	if orig.X == next.X {
		x = next.X
		if orig.Y > next.Y {
			y = next.Y - 1
		} else {
			y = next.Y + 1
		}
	} else {
		if orig.X < next.X {
			x = next.X + 1
		} else {
			x = next.X - 1
		}
		if (orig.X % 2 == 0 && orig.Y > next.Y) || (orig.X % 2 != 0 && orig.Y < next.Y) {
			y = next.Y
		} else if orig.X % 2 == 0 {
			y = next.Y + 1
		} else {
			y = next.Y - 1
		}
	}
	return Coords{
		X: x,
		Y: y,
	}
}

func NextHexRot(orig, pivot Coords, right bool) Coords {
	sextant := GetSextant(orig, pivot)
	move := 0
	if right {
		if sextant == LineUpLeft || sextant == TopLeft {
			// move up right
			move = 1
		} else if sextant == LineUp || sextant == TopRight {
			// move down right
			move = 2
		} else if sextant == LineUpRight || sextant == Right {
			// move down
			move = 3
		} else if sextant == LineDownRight || sextant == BottomRight {
			// move down left
			move = 4
		} else if sextant == LineDown || sextant == BottomLeft {
			// move up left
			move = 5
		} else {
			// move up
			move = 0
		}
	} else {
		if sextant == TopLeft || sextant == LineUp {
			// move down left
			move = 4
		} else if sextant == TopRight || sextant == LineUpRight {
			// move up left
			move = 5
		} else if sextant == Right || sextant == LineDownRight {
			// move up
			move = 0
		} else if sextant == BottomRight || sextant == LineDown {
			// move up right
			move = 1
		} else if sextant == BottomLeft || sextant == LineDownLeft {
			// move down right
			move = 2
		} else {
			// move down
			move = 3
		}
	}
	even := orig.X % 2 == 0
	x := orig.X
	y := orig.Y
	switch move {
	case 0:
		// move up
		y++
	case 1:
		// move up right
		x++
		if !even {
			y++
		}
	case 2:
		// move down right
		x++
		if even {
			y--
		}
	case 3:
		// move down
		y--
	case 4:
		// move down left
		x--
		if even {
			y--
		}
	case 5:
		// move up left
		x--
		if !even {
			y++
		}
	}
	return Coords{
		X: x,
		Y: y,
	}
}

func ReverseList(s []Coords) []Coords {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func RandomizeList(s []Coords) []Coords {
	for i := len(s)-1; i > 0; i-- {
		j := rand.Intn(i)
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func Remove(c Coords, list []Coords) []Coords {
	in := -1
	for i, l := range list {
		if c.Equals(l) {
			in = i
		}
	}
	if in != -1 {
		return append(list[:in], list[in+1:]...)
	} else {
		return list
	}
}

func CoordsIn(c Coords, list []Coords) bool {
	for _, l := range list {
		if c.Equals(l) {
			return true
		}
	}
	return false
}

// AngleBetween returns the angle between the two coordinates.
// If a is above b, the angle will be positive.
func AngleBetween(a, b Coords) float64 {
	mag := MapToWorld(a).Sub(MapToWorld(b))
	return mag.Angle()
}

type Sextant int

const (
	LineUp = iota
	TopRight
	LineUpRight
	Right
	LineDownRight
	BottomRight
	LineDown
	BottomLeft
	LineDownLeft
	Left
	LineUpLeft
	TopLeft
)

func GetSextantBias(subject, pivot Coords, biasRight bool) Sextant {
	sextant := GetSextant(subject, pivot)
	if sextant == LineUp {
		if biasRight {
			return TopRight
		} else {
			return TopLeft
		}
	} else if sextant == LineUpRight {
		if biasRight {
			return Right
		} else {
			return TopRight
		}
	} else if sextant == LineDownRight {
		if biasRight {
			return BottomRight
		} else {
			return Right
		}
	} else if sextant == LineDown {
		if biasRight {
			return BottomLeft
		} else {
			return BottomRight
		}
	} else if sextant == LineDownLeft {
		if biasRight {
			return Left
		} else {
			return BottomLeft
		}
	} else if sextant == LineUpLeft {
		if biasRight {
			return TopLeft
		} else {
			return Left
		}
	}
	return sextant
}

func GetSextant(subject, pivot Coords) Sextant {
	angle := AngleBetween(subject, pivot)
	const (
		upl = 2.67
		upl1 = 2.68
		top = 1.57
		top1 = 1.58
		upr = 0.46
		upr1 = 0.47
		dnr = -upr
		dnr1 = -upr1
		dwn = -top
		dwn1 = -top1
		dnl = -upl
		dnl1 = -upl1
	)
	if angle <= upl1 && angle >= upl {
		return LineUpLeft
	} else if angle <= upl && angle >= top1 {
		return TopLeft
	} else if angle <= top1 && angle >= top {
		return LineUp
	} else if angle <= top && angle >= upr1 {
		return TopRight
	} else if angle <= upr1 && angle >= upr {
		return LineUpRight
	} else if angle <= upr && angle >= dnr {
		return Right
	} else if angle <= dnr && angle >= dnr1 {
		return LineDownRight
	} else if angle <= dnr1 && angle >= dwn {
		return BottomRight
	} else if angle <= dwn && angle >= dwn1 {
		return LineDown
	} else if angle <= dwn1 && angle >= dnl {
		return BottomLeft
	} else if angle <= dnl && angle >= dnl1 {
		return LineDownLeft
	} else {
		return Left
	}
}

func GetSextantWorld(subject, pivot pixel.Vec) Sextant {
	mag := subject.Sub(pivot)
	angle := mag.Angle()
	const (
		upl = (math.Pi/6.)*5.
		top = math.Pi/2.
		upr = math.Pi/6.
		dnr = -upr
		dwn = -top
		dnl = -upl
	)
	if angle <= upl && angle >= top {
		return TopLeft
	} else if angle <= top && angle >= upr {
		return TopRight
	} else if angle <= upr && angle >= dnr {
		return Right
	} else if angle <= dnr && angle >= dwn {
		return BottomRight
	} else if angle <= dwn && angle >= dnl {
		return BottomLeft
	} else {
		return Left
	}
}