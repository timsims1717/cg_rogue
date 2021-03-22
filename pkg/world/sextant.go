package world

import (
	"github.com/faiface/pixel"
	"math"
)

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

// AngleBetween returns the angle between the two coordinates.
// If a is above b, the angle will be positive.
func AngleBetween(a, b Coords) float64 {
	mag := MapToWorld(a).Sub(MapToWorld(b))
	return mag.Angle()
}

// GetSextant returns the sextant the subject Coords belongs to
// relative to the pivot.
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

// GetSextant returns the sextant the subject Coords belongs to
// relative to the pivot. When the subject Coords is aligned with
// the pivot, it will be changed to one of the six pies, based on
// the bias.
func GetSextantBias(subject, pivot Coords, biasClockwise bool) Sextant {
	sextant := GetSextant(subject, pivot)
	if sextant == LineUp {
		if biasClockwise {
			return TopRight
		} else {
			return TopLeft
		}
	} else if sextant == LineUpRight {
		if biasClockwise {
			return Right
		} else {
			return TopRight
		}
	} else if sextant == LineDownRight {
		if biasClockwise {
			return BottomRight
		} else {
			return Right
		}
	} else if sextant == LineDown {
		if biasClockwise {
			return BottomLeft
		} else {
			return BottomRight
		}
	} else if sextant == LineDownLeft {
		if biasClockwise {
			return Left
		} else {
			return BottomLeft
		}
	} else if sextant == LineUpLeft {
		if biasClockwise {
			return TopLeft
		} else {
			return Left
		}
	}
	return sextant
}

// GetSextantWorld returns the sextant the subject pixel.Vec belongs to
// relative to the pivot.
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