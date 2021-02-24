package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
)

type PushAction struct{
	source *characters.Character
	target objects.Moveable
	push   int
	start  bool
	isDone bool
}

func NewPushAction(s *characters.Character, t objects.Moveable, p int) *PushAction {
	if t == nil {
		return nil
	}
	return &PushAction{
		source: s,
		target: t,
		push:   p,
		start:  true,
		isDone: false,
	}
}

func (a *PushAction) Update() {
	if a.start {
		tOrig := a.target.GetCoords()
		oPath, _, found := floor.CurrentFloor.FindPath(a.source.GetCoords(), tOrig, floor.NoCheck)
		if !found {
			a.isDone = true
			return
		}
		nPath := tOrig.PathFrom(oPath)
		if len(nPath) == 0 {
			a.isDone = true
			return
		}
		path := floor.CurrentFloor.Line(a.target.GetCoords(), nPath[len(nPath)-1], a.push, floor.PathChecks{
			NotFilled:  true,
			Unoccupied: true,
			NonEmpty:   false,
			Orig:       a.target.GetCoords(),
		})
		if len(path) == 0 {
			a.isDone = true
			return
		}
		AddToTop(NewMoveAction(a.source, a.target, path[len(path)-1]))
		a.start = false
		a.isDone = true
	}
}

func (a *PushAction) IsDone() bool {
	return a.isDone
}