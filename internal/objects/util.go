package objects

import (
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"reflect"
)

type Targetable interface {
	Damage(dmg int)
}

// NotNil checks both if i is nil, and if the underlying
// interface is nil.
func NotNil(i interface{}) bool {
	if i == nil {
		return false
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(i).IsNil() {
			return false
		}
	}
	return true
}

type Moveable interface {
	GetCoords() world.Coords
	SetCoords(world.Coords)
	GetXY() (float64, float64)
	SetXY(float64, float64)
}

// notNil checks both if o is nil, and if the underlying
// Occupant is nil.
func NotNilMov(m Moveable) bool {
	if m == nil {
		return false
	}
	switch reflect.TypeOf(m).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(m).IsNil() {
			return false
		}
	}
	return true
}