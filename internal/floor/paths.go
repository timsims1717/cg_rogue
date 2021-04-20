package floor

import (
	"github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
	"github.com/phf/go-queue/queue"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

// IsLegal checks if a Hex is a legal candidate according to the current PathCheck
func (f *Floor) IsLegal(a world.Coords, checks PathChecks) *Hex {
	hex := f.Get(a)
	if hex != nil && ((a.X == checks.Orig.X && a.Y == checks.Orig.Y) ||
		((!checks.Unoccupied || (!f.IsOccupied(a) && (!checks.HonorClaim || !f.IsClaimed(a)))) &&
		(!checks.NonEmpty || !hex.Empty))) {
		return hex
	}
	return nil
}

// isLegal checks if a Hex is a legal candidate according to the current PathCheck
func (f *Floor) isLegal(a world.Coords) *Hex {
	hex := f.Get(a)
	if hex != nil && ((a.X == f.checks.Orig.X && a.Y == f.checks.Orig.Y) ||
		((!f.checks.Unoccupied || (!f.IsOccupied(a) && (!f.checks.HonorClaim || !f.IsClaimed(a)))) &&
		(!f.checks.NonEmpty || !hex.Empty))) {
		return hex
	}
	return nil
}

// IsSetLegal filters out the world.Coords that are not legal according to the PathChecks
func (f *Floor) IsSetLegal(set []world.Coords, checks PathChecks) []world.Coords {
	var legal []world.Coords
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

// AllInSextant
func (f *Floor) AllInSextant(orig, ref world.Coords, d int, check PathChecks) []world.Coords {
	f.checks = check
	sextant := world.GetSextantBias(ref, orig, rand.Intn(2)%2 == 0)
	width, height := f.Dimensions()
	type cont struct {
		c world.Coords
		w int
	}
	all := make([]world.Coords, 0)
	qu := queue.New()
	marked := make(map[world.Coords]bool)
	marked[orig] = true
	qu.PushFront(cont{c: orig, w: 0})
	for qu.Front() != nil {
		n := qu.PopFront()
		if c, ok := n.(cont); ok {
			all = append(all, c.c)
			if c.w < d {
				neighbors := c.c.Neighbors(width, height)
				for _, nb := range neighbors {
					if !marked[nb] {
						marked[nb] = true
						if f.isLegal(nb) != nil && orig.InSextant(nb, sextant) {
							qu.PushBack(cont{c: nb, w: c.w + 1})
						}
					}
				}
			}
		}
	}
	f.checks = DefaultCheck
	return all
}

// AllWithin returns all legal coordinates within d tiles from o that can be reached
func (f *Floor) AllWithin(orig world.Coords, d int, check PathChecks) []world.Coords {
	f.checks = check
	width, height := f.Dimensions()
	type cont struct {
		c world.Coords
		w int
	}
	all := make([]world.Coords, 0)
	qu := queue.New()
	marked := make(map[world.Coords]bool)
	marked[orig] = true
	qu.PushFront(cont{c: orig, w: 0})
	for qu.Front() != nil {
		n := qu.PopFront()
		if c, ok := n.(cont); ok {
			all = append(all, c.c)
			if c.w < d {
				neighbors := c.c.Neighbors(width, height)
				for _, nb := range neighbors {
					if !marked[nb] {
						marked[nb] = true
						if f.isLegal(nb) != nil {
							qu.PushBack(cont{c: nb, w: c.w + 1})
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
func (f *Floor) AllWithinNoPath(orig world.Coords, d int, check PathChecks) []world.Coords {
	f.checks = check
	defer f.SetDefaultChecks()
	width, height := f.Dimensions()
	type cont struct {
		c world.Coords
		w int
	}
	all := make([]world.Coords, 0)
	qu := queue.New()
	marked := make(map[world.Coords]bool)
	marked[orig] = true
	qu.PushFront(cont{c: orig, w: 0})
	for qu.Front() != nil {
		n := qu.PopFront()
		if c, ok := n.(cont); ok {
			if f.isLegal(c.c) != nil {
				all = append(all, c.c)
			}
			if c.w < d {
				neighbors := c.c.Neighbors(width, height)
				for _, nb := range neighbors {
					if !marked[nb] {
						marked[nb] = true
						qu.PushBack(cont{c: nb, w: c.w + 1})
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
func (f *Floor) LongestLegalPath(path []world.Coords, max int, check PathChecks) []world.Coords {
	f.checks = check
	defer f.SetDefaultChecks()
	lastLegal := 0
	for i, c := range path {
		if (!check.EndUnoccupied || !f.IsOccupied(c)) && (!check.HonorClaim || !f.IsClaimed(c)) {
			lastLegal = i
		}
		if h := f.isLegal(c); h == nil {
			break
		}
		if max > 0 && i >= max {
			break
		}
	}
	if lastLegal+1 >= len(path) {
		return path
	}
	return path[:lastLegal+1]
}

// FindPathPerpendicularTo finds a semi-random path perpendicular to the specified world.Coords.
func (f *Floor) FindPathPerpendicularTo(orig, to world.Coords, within, dist int, check, endCheck PathChecks) ([]world.Coords, int, bool) {
	distTo := world.DistanceSimple(orig, to)
	allPossible, _ := world.Remove(orig, f.AllWithin(orig, within, check))
	var possible []world.Coords
	for _, c := range allPossible {
		if world.DistanceSimple(to, c) <= dist {
			possible = append(possible, c)
		}
	}
	if len(possible) > 0 {
		ordered := world.OrderByDistDiff(to, possible, distTo)
		if len(ordered) > 4 {
			l := len(ordered) / 4
			choice := rand.Intn(l)
			return f.FindPath(orig, ordered[choice], check)
		} else {
			return f.FindPath(orig, ordered[0], check)
		}
	}
	return []world.Coords{}, 0, false
}

// FindPathAwayFrom finds a semi-random path away from the specified world.Coords.
func (f *Floor) FindPathAwayFrom(orig, from world.Coords, dist int, check PathChecks) ([]world.Coords, int, bool) {
	possible, _ := world.Remove(orig, f.AllWithin(orig, dist, check))
	if len(possible) > 0 {
		ordered := world.ReverseList(world.OrderByDistSimple(from, possible))
		if len(ordered) > 6 {
			choice := rand.Intn(len(ordered) / 6)
			return f.FindPath(orig, ordered[choice], check)
		} else {
			return f.FindPath(orig, ordered[0], check)
		}
	}
	return []world.Coords{}, 0, false
}

// FindPathWithinOne runs astar from a to one within b, returning a world.Coords array
func (f *Floor) FindPathWithinOne(a, b world.Coords, check PathChecks) ([]world.Coords, int, bool) {
	f.checks = check
	defer f.SetDefaultChecks()
	if !f.Exists(a) {
		return nil, 0, false
	}
	for _, n := range world.OrderByDist(a, b.Neighbors(f.Dimensions())) {
		if a.Eq(n) {
			return []world.Coords{a}, 0, true
		}
		if !f.Exists(n) || (check.EndUnoccupied && f.IsOccupied(n)) || (check.HonorClaim && f.IsClaimed(n)) {
			continue
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
	if !f.Exists(a) {
		return nil, 0, false
	}
	for _, n := range world.OrderByDist(a, b.Neighbors(f.Dimensions())) {
		if a.Eq(n) {
			return []*Hex{f.Get(a)}, 0, true
		}
		if !f.Exists(n) || (check.EndUnoccupied && f.IsOccupied(n)) || (check.HonorClaim && f.IsClaimed(n)) {
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
	if !f.Exists(a) || !f.Exists(b) {
		return nil, 0, false
	}
	if a.Eq(b) {
		return []world.Coords{a}, 0, true
	}
	if (check.EndUnoccupied && f.IsOccupied(b)) || (check.HonorClaim && f.IsClaimed(b)) {
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
	if !f.Exists(a) || !f.Exists(b) {
		return nil, 0, false
	}
	if a.Eq(b) {
		return []*Hex{f.Get(a)}, 0, true
	}
	if (check.EndUnoccupied && f.IsOccupied(b)) || (check.HonorClaim && f.IsClaimed(b)) {
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
	if hex == nil {
		return []*Hex{}
	}
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
	f.PathLine = pixel.L(world.MapToWorld(a), world.MapToWorld(b))
}
