package objects

import "reflect"

type Occupant interface {
	TryToOccupy()
}

// notNil checks both if o is nil, and if the underlying
// Occupant is nil.
func NotNilOcc(o Occupant) bool {
	if o == nil {
		return false
	}
	switch reflect.TypeOf(o).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(o).IsNil() {
			return false
		}
	}
	return true
}

type Targetable interface {
	Damage(dmg int)
	IsTargeted()
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