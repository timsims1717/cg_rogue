package floor

import (
	"github.com/beefsack/go-astar"
	"github.com/phf/go-queue/queue"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

// IsLegal checks if a Hex is a legal candidate according to the current PathCheck
func (f *Floor) IsLegal(a world.Coords) *Hex {
	hex := f.Get(a)
	if hex != nil {
		if (a.X == f.checks.Orig.X && a.Y == f.checks.Orig.Y) || ((!f.checks.Unoccupied || hex.Occupant == nil) && (!f.checks.NonEmpty || !hex.Empty)) {
			return hex
		}
	}
	return nil
}

func (f *Floor) Line(orig, ref world.Coords, dist int, check PathChecks) []world.Coords {
	dist += 1
	path := make([]world.Coords, 0)
	path, d, found := f.FindPath(orig, ref, NoCheck)
	if !found || len(path) < 2 {
		return []world.Coords{}
	}
	for dist > d {
		nPath := path[len(path)-1].PathFrom(path)
		path = append(path[:len(path)-1], nPath...)
		d += len(nPath)
	}
	nPath := f.LegalPath(path, check)
	if len(nPath) > dist {
		return nPath[:dist]
	}
	return nPath
}

// AllWithin returns all legal coordinates within d tiles from o
func (f *Floor) AllWithin(o world.Coords, d int, check PathChecks) []world.Coords {
	f.checks = check
	width, height := f.Dimensions()
	type cont struct{
		c world.Coords
		w int
	}
	all := make([]world.Coords, 0)
	qu := queue.New()
	marked := make(map[world.Coords]bool)
	qu.PushFront(cont{ c: o, w: 0 })
	for n := qu.PopFront(); n != nil; {
		if c, ok := n.(cont); ok {
			if c.w+1 <= d {
				all = append(all, c.c)
				neighbors := c.c.Neighbors(width, height)
				for _, nb := range neighbors {
					if !marked[nb] {
						marked[nb] = true
						if f.IsLegal(nb) != nil {
							qu.PushBack(nb)
						}
					}
				}
			}
		}
	}
	f.checks = DefaultCheck
	return all
}

// LegalPath returns the longest section of the given path that is legal
// If the path is (0,0), (0,1), (0,2), (0,3), but (0,2) is not legal, it
// returns (0,0), (0,1)
func (f *Floor) LegalPath(path []world.Coords, check PathChecks) []world.Coords {
	f.checks = check
	for i, c := range path {
		if f.IsLegal(c) == nil {
			f.checks = DefaultCheck
			return path[:i]
		}
	}
	f.checks = DefaultCheck
	return path
}

// FindPath runs astar from a to b, returning just the world.Coords
func (f *Floor) FindPath(a, b world.Coords, check PathChecks) ([]world.Coords, int, bool) {
	f.checks = check
	pathA, distance, found := astar.Path(f.Get(b), f.Get(a))
	var path []*Hex
	for _, h := range pathA {
		path = append(path, h.(*Hex))
	}
	f.checks = DefaultCheck
	var cpath []world.Coords
	for _, p := range path {
		cpath = append(cpath, world.Coords{
			X: p.X,
			Y: p.Y,
		})
	}
	return cpath, int(distance), found
}

// FindPathHex runs astar from a to b returning a Hex array
func (f *Floor) FindPathHex(a, b world.Coords, check PathChecks) ([]*Hex, int, bool) {
	f.checks = check
	pathA, distance, found := astar.Path(f.Get(b), f.Get(a))
	var path []*Hex
	for _, h := range pathA {
		path = append(path, h.(*Hex))
	}
	f.checks = DefaultCheck
	return path, int(distance), found
}

// Neighbors returns each legal hex adjacent to the origin
func (f *Floor) Neighbors(hex *Hex) []*Hex {
	width, height := f.Dimensions()
	co := world.Coords{X: hex.X, Y: hex.Y}
	cNeighbors := co.Neighbors(width, height)
	neighbors := make([]*Hex, 0)
	for _, c := range cNeighbors {
		if n := f.IsLegal(c); n != nil {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

