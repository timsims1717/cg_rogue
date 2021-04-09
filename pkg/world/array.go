package world

import "math/rand"

// ReverseList reverses the order of the Coords array.
func ReverseList(s []Coords) []Coords {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// RandomizeList shuffles the elements of the Coords array.
func RandomizeList(s []Coords) []Coords {
	for i := len(s)-1; i > 0; i-- {
		j := rand.Intn(i)
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Remove attempts to remove the Coords c from the list.
// If c is not in the list, no change is made.
func Remove(c Coords, list []Coords) ([]Coords, int) {
	in := -1
	for i, l := range list {
		if c.Eq(l) {
			in = i
		}
	}
	if in != -1 {
		return append(list[:in], list[in+1:]...), in
	} else {
		return list, in
	}
}

// CoordsIn returns true if Coords c are in the list.
func CoordsIn(c Coords, list []Coords) bool {
	for _, l := range list {
		if c.Eq(l) {
			return true
		}
	}
	return false
}


func Combine(a, b []Coords) []Coords {
	n := a
	for _, c := range b {
		found := false
		for _, d := range a {
			if c == d {
				found = true
				break
			}
		}
		if !found {
			n = append(n, c)
		}
	}
	return n
}