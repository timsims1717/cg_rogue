package world

func Spiral(l int) []Coords {
	list := []Coords{Origin}
	start := Origin
	first := Origin.Up()
	curr := first
	for i := 0; i < l; i++ {
		list = append(list, curr)
		next := NextHexRot(curr, start, true)
		if next == first {
			next = curr.Up()
		}
		curr = next
	}
	return list
}
