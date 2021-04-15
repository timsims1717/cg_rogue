package world

var TileSize float64
var Origin = Coords{
	X: 0,
	Y: 0,
}

// Coords is a convenience struct for passing tile coordinates
type Coords struct {
	X int
	Y int
}

func (a Coords) Direction(b Coords) Sextant {
	return GetSextant(b, a)
}

// Up returns the Coords above
func (a Coords) Up() Coords {
	return Coords{X: a.X, Y: a.Y + 1}
}

// Down returns the Coords below
func (a Coords) Down() Coords {
	return Coords{X: a.X, Y: a.Y - 1}
}

func (a Coords) LeftUp() Coords {
	if a.X%2 == 0 {
		return Coords{X: a.X - 1, Y: a.Y}
	} else {
		return Coords{X: a.X - 1, Y: a.Y + 1}
	}
}

func (a Coords) LeftDown() Coords {
	if a.X%2 == 0 {
		return Coords{X: a.X - 1, Y: a.Y - 1}
	} else {
		return Coords{X: a.X - 1, Y: a.Y}
	}
}

func (a Coords) RightUp() Coords {
	if a.X%2 == 0 {
		return Coords{X: a.X + 1, Y: a.Y}
	} else {
		return Coords{X: a.X + 1, Y: a.Y + 1}
	}
}

func (a Coords) RightDown() Coords {
	if a.X%2 == 0 {
		return Coords{X: a.X + 1, Y: a.Y - 1}
	} else {
		return Coords{X: a.X + 1, Y: a.Y}
	}
}

// Eq checks if a and b are equal.
func (a Coords) Eq(b Coords) bool {
	return a.X == b.X && a.Y == b.Y
}

// Neighbors returns the six tiles surrounding the Coords, minus any
// outside the width and height provided.
func (a Coords) Neighbors(w, h int) []Coords {
	neighbors := make([]Coords, 0)
	if a.Y > 0 {
		neighbors = append(neighbors, Coords{a.X, a.Y - 1})
	}
	if a.X < w-1 {
		neighbors = append(neighbors, Coords{a.X + 1, a.Y})
	}
	if a.Y < h-1 {
		neighbors = append(neighbors, Coords{a.X, a.Y + 1})
	}
	if a.X > 0 {
		neighbors = append(neighbors, Coords{a.X - 1, a.Y})
	}
	if a.X%2 == 0 {
		if a.X < w-1 && a.Y > 0 {
			neighbors = append(neighbors, Coords{a.X + 1, a.Y - 1})
		}
		if a.X > 0 && a.Y > 0 {
			neighbors = append(neighbors, Coords{a.X - 1, a.Y - 1})
		}
	} else {
		if a.X < w-1 && a.Y < h-1 {
			neighbors = append(neighbors, Coords{a.X + 1, a.Y + 1})
		}
		if a.X > 0 && a.Y < h-1 {
			neighbors = append(neighbors, Coords{a.X - 1, a.Y + 1})
		}
	}
	return neighbors
}

// PathFrom constructs an new path following the relative hex changes
// of the input path, but starts at Coords a.
func (a Coords) PathFrom(path []Coords) []Coords {
	if len(path) < 2 {
		return []Coords{a}
	}
	orig := path[0]
	if a.Eq(orig) {
		return path
	}
	var np []Coords
	dx := a.X - orig.X
	dy := a.Y - orig.Y
	for _, p := range path {
		if a.X%2 == 0 && p.X%2 == 0 && a.X%2 != orig.X%2 {
			n := Coords{
				X: p.X + dx,
				Y: p.Y + dy - 1,
			}
			np = append(np, n)
		} else if a.X%2 != 0 && p.X%2 != 0 && a.X%2 != orig.X%2 {
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

// InSextant checks to see if Coords b is in Sextant s relative to Coords
// a. If Sextant s is a line Sextant, b must be on that line to return
// true. If b in on a line, it will still return true as long as b is
// still in the sextant when it is biased right or left.
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

// NextHexLine returns the Coords directly opposite the orig Coords relative
// to the next Coords.
// todo: change this function to use angle/sextant
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
		if (orig.X%2 == 0 && orig.Y > next.Y) || (orig.X%2 != 0 && orig.Y < next.Y) {
			y = next.Y
		} else if orig.X%2 == 0 {
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

// NextHexRot returns the Coords clockwise from the orig relative to the
// pivot if clockwise is true, or the Coords counter-clockwise from the orig
// relative to the pivot otherwise.
func NextHexRot(orig, pivot Coords, clockwise bool) Coords {
	sextant := GetSextant(orig, pivot)
	even := orig.X%2 == 0
	x := orig.X
	y := orig.Y
	if clockwise {
		if sextant == LineUpLeft || sextant == TopLeft {
			// move up right
			x++
			if !even {
				y++
			}
		} else if sextant == LineUp || sextant == TopRight {
			// move down right
			x++
			if even {
				y--
			}
		} else if sextant == LineUpRight || sextant == Right {
			// move down
			y--
		} else if sextant == LineDownRight || sextant == BottomRight {
			// move down left
			x--
			if even {
				y--
			}
		} else if sextant == LineDown || sextant == BottomLeft {
			// move up left
			x--
			if !even {
				y++
			}
		} else {
			// move up
			y++
		}
	} else {
		if sextant == TopLeft || sextant == LineUp {
			// move down left
			x--
			if even {
				y--
			}
		} else if sextant == TopRight || sextant == LineUpRight {
			// move up left
			x--
			if !even {
				y++
			}
		} else if sextant == Right || sextant == LineDownRight {
			// move up
			y++
		} else if sextant == BottomRight || sextant == LineDown {
			// move up right
			x++
			if !even {
				y++
			}
		} else if sextant == BottomLeft || sextant == LineDownLeft {
			// move down right
			x++
			if even {
				y--
			}
		} else {
			// move down
			y--
		}
	}
	return Coords{
		X: x,
		Y: y,
	}
}
