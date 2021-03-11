package util

import "github.com/faiface/pixel"

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Min returns the smaller number between a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger number between a and b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// PointInside returns true if the pixel.Vec is inside the pixel.Rect
// when unprojected by the pixel.Matrix
func PointInside(p pixel.Vec, r pixel.Rect, m pixel.Matrix) bool {
	return r.Moved(pixel.V(-(r.W() / 2.0), -(r.H() / 2.0))).Contains(m.Unproject(p))
}