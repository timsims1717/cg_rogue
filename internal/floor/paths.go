package floor

import (
	"github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
	"github.com/phf/go-queue/queue"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)


// IsLegal checks if a Hex is a legal candidate according to the current PathCheck
func (f *Floor) IsLegal(a world.Coords, checks PathChecks) *Hex {
	hex := f.Get(a)
	if hex != nil {
		if (a.X == checks.Orig.X && a.Y == checks.Orig.Y) || ((!checks.Unoccupied || hex.Occupant == nil) && (!checks.NonEmpty || !hex.Empty)) {
			return hex
		}
	}
	return nil
}

// isLegal checks if a Hex is a legal candidate according to the current PathCheck
func (f *Floor) isLegal(a world.Coords) *Hex {
	hex := f.Get(a)
	if hex != nil {
		if (a.X == f.checks.Orig.X && a.Y == f.checks.Orig.Y) || ((!f.checks.Unoccupied || hex.Occupant == nil) && (!f.checks.NonEmpty || !hex.Empty)) {
			return hex
		}
	}
	return nil
}

// IsSetLegal filters out the world.Coords that are not legal according to the PathChecks
func (f *Floor) IsSetLegal(set []world.Coords, checks PathChecks) []world.Coords {
	legal := []world.Coords{}
	for _, c := range set {
		if f.IsLegal(c, checks) != nil {
			legal = append(legal, c)
		}
	}
	return legal
}

// Line returns a straight path from orig through ref out to dist tiles.
// The path returned can be reduced by using LongestLegalPath.
func (f *Floor) Line(orig, ref world.Coords, dist int) []world.Coords {
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
	if len(path) > dist {
		return path[:dist]
	}
	return path
}

// AllWithin returns all legal coordinates within d tiles from o that can be reached
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
	for qu.Front() != nil {
		n := qu.PopFront()
		if c, ok := n.(cont); ok {
			if c.w+1 <= d {
				all = append(all, c.c)
				neighbors := c.c.Neighbors(width, height)
				for _, nb := range neighbors {
					if !marked[nb] {
						marked[nb] = true
						if f.isLegal(nb) != nil {
							qu.PushBack(cont{ c: nb, w: c.w+1 })
						}
					}
				}
			}
		}
	}
	f.checks = DefaultCheck
	return all
}

// AllWithinNoPath returns all legal coordinates within d tiles from o
func (f *Floor) AllWithinNoPath(o world.Coords, d int, check PathChecks) []world.Coords {
	f.checks = check
	defer f.SetDefaultChecks()
	width, height := f.Dimensions()
	type cont struct{
		c world.Coords
		w int
	}
	all := make([]world.Coords, 0)
	qu := queue.New()
	marked := make(map[world.Coords]bool)
	qu.PushFront(cont{ c: o, w: 0 })
	for qu.Front() != nil {
		n := qu.PopFront()
		if c, ok := n.(cont); ok {
			if c.w+1 <= d {
				if f.isLegal(c.c) != nil {
					all = append(all, c.c)
				}
				neighbors := c.c.Neighbors(width, height)
				for _, nb := range neighbors {
					if !marked[nb] {
						marked[nb] = true
						qu.PushBack(cont{ c: nb, w: c.w+1 })
					}
				}
			}
		}
	}
	return all
}

// LongestLegalPath returns the longest section of the given path that is legal
// If the path is (0,0), (0,1), (0,2), (0,3), but (0,2) is not legal, it
// returns (0,0), (0,1)
func (f *Floor) LongestLegalPath(path []world.Coords, check PathChecks) []world.Coords {
	f.checks = check
	defer f.SetDefaultChecks()
	lastLegal := 0
	for i, c := range path {
		if h := f.isLegal(c); h == nil {
			return path[:lastLegal+1]
		}
		if !check.EndUnoccupied || !f.HasOccupant(c) {
			lastLegal = i
		}
	}
	if lastLegal + 1 >= len(path) {
		return path
	}
	return path[:lastLegal+1]
}

// FindPathWithinOne runs astar from a to one within b, returning a world.Coords array
func (f *Floor) FindPathWithinOne(a, b world.Coords, check PathChecks) ([]world.Coords, int, bool) {
	f.checks = check
	defer f.SetDefaultChecks()
	for _, n := range world.OrderByDist(a, b.Neighbors(f.Dimensions())) {
		if  !f.Exists(a) || !f.Exists(n) || (check.EndUnoccupied && f.HasOccupant(n)) {
			return nil, 0, false
		}
		f.SetLine(a, b)
		pathA, distance, found := astar.Path(f.Get(n), f.Get(a))
		if !found {
			continue
		}
		var path []*Hex
		for _, h := range pathA {
			path = append(path, h.(*Hex))
		}
		var cpath []world.Coords
		for _, p := range path {
			cpath = append(cpath, world.Coords{
				X: p.X,
				Y: p.Y,
			})
		}
		return cpath, int(distance), found
	}
	return nil, 0, false
}

// FindPathWithinOneHex runs astar from a to one within b returning a Hex array
func (f *Floor) FindPathWithinOneHex(a, b world.Coords, check PathChecks) ([]*Hex, int, bool) {
	f.checks = check
	defer f.SetDefaultChecks()
	for _, n := range world.OrderByDist(a, b.Neighbors(f.Dimensions())) {
		if  !f.Exists(a) || !f.Exists(n) || (check.EndUnoccupied && f.HasOccupant(n)) {
			return nil, 0, false
		}
		f.SetLine(a, b)
		pathA, distance, found := astar.Path(f.Get(n), f.Get(a))
		if !found {
			continue
		}
		var path []*Hex
		for _, h := range pathA {
			path = append(path, h.(*Hex))
		}
		return path, int(distance), found
	}
	return nil, 0, false
}

// FindPath runs astar from a to b, returning a world.Coords array
func (f *Floor) FindPath(a, b world.Coords, check PathChecks) ([]world.Coords, int, bool) {
	if  !f.Exists(a) || !f.Exists(b) || (check.EndUnoccupied && f.HasOccupant(b)) {
		return nil, 0, false
	}
	f.checks = check
	defer f.SetDefaultChecks()
	f.SetLine(a, b)
	pathA, distance, found := astar.Path(f.Get(b), f.Get(a))
	var path []*Hex
	for _, h := range pathA {
		path = append(path, h.(*Hex))
	}
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
	if  !f.Exists(a) || !f.Exists(b) || (check.EndUnoccupied && f.HasOccupant(b)) {
		return nil, 0, false
	}
	f.checks = check
	defer f.SetDefaultChecks()
	f.SetLine(a, b)
	pathA, distance, found := astar.Path(f.Get(b), f.Get(a))
	var path []*Hex
	for _, h := range pathA {
		path = append(path, h.(*Hex))
	}
	return path, int(distance), found
}

// Neighbors returns each legal hex adjacent to the origin
func (f *Floor) Neighbors(hex *Hex) []*Hex {
	width, height := f.Dimensions()
	co := world.Coords{X: hex.X, Y: hex.Y}
	cNeighbors := co.Neighbors(width, height)
	neighbors := make([]*Hex, 0)
	for _, c := range cNeighbors {
		if n := f.isLegal(c); n != nil {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (f *Floor) SetLine(a, b world.Coords) {
	ax, ay := world.MapToWorld(a.X, a.Y)
	bx, by := world.MapToWorld(b.X, b.Y)
	f.PathLine = pixel.L(pixel.V(ax, ay), pixel.V(bx, by))
}