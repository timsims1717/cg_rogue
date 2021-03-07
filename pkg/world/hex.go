package world

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