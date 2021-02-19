package objects

//type Entity interface {
//	Update()
//	Draw(win *pixelgl.Window)
//	ID() uuid.UUID
//	Damage(dmg int)
//	IsDestroyed() bool
//}
//
//var AllEntities entities
//
//type entities struct{
//	set [][]Pile
//}
//
//type Pile struct {
//	Target Targetable
//	Pile   []Entity
//}
//
//// Initializes the Entity set.
//func (e *entities) Initialize(w, h int) {
//	e.set = make([][]Pile, 0)
//	for x := 0; x < w; x++ {
//		e.set = append(e.set, make([]Pile, 0))
//		for y := 0; y < h; y++ {
//			e.set[x] = append(e.set[x], Pile{
//				Target: nil,
//				Pile:   []Entity{},
//			})
//		}
//	}
//}
//
//// Dimensions returns the width and height of the Entity set.
//func (e *entities) Dimensions() (int, int) {
//	width := len(e.set)
//	height := len(e.set[0])
//	return width, height
//}
//
//// Exists returns true if the world.Coords exists inside the set.
//func (e *entities) Exists(a world.Coords) bool {
//	w, h := e.Dimensions()
//	return a.X >= 0 && a.Y >= 0 && a.X < w && a.Y < h
//}
//
//// GetPile returns the Entity array located at world.Coords a
//func (e *entities) GetPile(a world.Coords) []Entity {
//	if e.Exists(a) {
//		return e.set[a.X][a.Y].Pile
//	}
//	return nil
//}
//
//// PutEntity an entity in the Pile at world.Coords a.
//// Returns true if successfully put.
//func (e *entities) PutEntity(entity Entity, a world.Coords) bool {
//	if e.Exists(a) {
//		e.set[a.X][a.Y].Pile = append(e.set[a.X][a.Y].Pile, entity)
//		return true
//	}
//	return false
//}
//
//// RemoveEntity an entity from the world.Coords. Returns true if the entity
//// was found and removed.
//func (e *entities) RemoveEntity(entity Entity, a world.Coords) bool {
//	if e.Exists(a) {
//		for i, ent := range e.set[a.X][a.Y].Pile {
//			if uuid.Equal(ent.ID(), entity.ID()) {
//				e.set[a.X][a.Y].Pile = append(e.set[a.X][a.Y].Pile[:i], e.set[a.X][a.Y].Pile[i+1:]...)
//				return true
//			}
//		}
//	}
//	return false
//}
//
//// MoveEntity an entity from world.Coords a to world.Coords b.
//// If the entity does not exist in a, no entity is moved.
//// Returns true if the entity was successfully moved.
//func (e *entities) MoveEntity(entity Entity, a, b world.Coords) bool {
//	if !e.Exists(a) || !e.Exists(b) {
//		return false
//	}
//	success := e.RemoveEntity(entity, a)
//	if success {
//		success = e.PutEntity(entity, b)
//	}
//	return success
//}
//
//// notNil checks both if e is nil, and if the underlying
//// Entity is nil.
//func notNil(e Entity) bool {
//	if e == nil {
//		return false
//	}
//	switch reflect.TypeOf(e).Kind() {
//	case reflect.Ptr:
//		if reflect.ValueOf(e).IsNil() {
//			return false
//		}
//	}
//	return true
//}