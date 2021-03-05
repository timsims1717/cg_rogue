package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
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
		o := a.source.GetCoords()
		n := a.target.GetCoords()
		checks := floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    true,
			NonEmpty:      false,
			EndUnoccupied: true,
			Orig:          n,
		}
		nPath := []world.Coords{n}
		for i := 0; i < a.push; i++ {
			next := world.NextHex(o, n)
			if floor.CurrentFloor.IsLegal(next, checks) != nil {
				nPath = append(nPath, next)
				o = n
				n = next
			} else {
				break
			}
		}
		if len(nPath) < 2 {
			a.isDone = true
			return
		}
		AddToTop(NewMoveAction(a.source, a.target, nPath[len(nPath)-1]))
		a.start = false
		a.isDone = true
	}
}

func (a *PushAction) IsDone() bool {
	return a.isDone
}