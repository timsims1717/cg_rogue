package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
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

type PushMultiAction struct{
	area    []world.Coords
	values  selectors.ActionValues
	start   bool
	isDone  bool
	targets []objects.Moveable
	ends    []world.Coords
	interX  []*gween.Tween
	interY  []*gween.Tween
}

func NewPushMultiAction(area []world.Coords, values selectors.ActionValues) *PushMultiAction {
	if len(area) < 1 {
		return nil
	}
	return &PushMultiAction{
		area:   area,
		values: values,
		start:  true,
		isDone: false,
	}
}

func (a *PushMultiAction) Update() {
	if a.start {
		o := a.values.Source.GetCoords()
		for _, n := range a.area {
			if target := floor.CurrentFloor.GetOccupant(n); target != nil {
				checks := floor.PathChecks{
					NotFilled:     true,
					Unoccupied:    true,
					NonEmpty:      false,
					EndUnoccupied: true,
					Orig:          n,
				}
				nPath := []world.Coords{n}
				for i := 0; i < a.values.Strength; i++ {
					next := world.NextHex(o, n)
					if floor.CurrentFloor.IsLegal(next, checks) != nil {
						nPath = append(nPath, next)
						o = n
						n = next
					} else {
						break
					}
				}
				a.targets = append(a.targets, target)
				if len(nPath) < 2 {
					a.ends = append(a.ends, target.GetCoords())
					a.interX = append(a.interX, nil)
					a.interY = append(a.interY, nil)
				} else {
					end := nPath[len(nPath)-1]
					b := world.MapToWorld(end)
					a.ends = append(a.ends, end)
					x, y := target.GetXY()
					a.interX = append(a.interX, gween.New(x, b.X, 0.25, ease.OutCubic))
					a.interY = append(a.interY, gween.New(y, b.Y, 0.25, ease.OutCubic))
				}
			}
		}
	}
	done := true
	for i, target := range a.targets {
		if a.start {
			if t, ok := target.(objects.Targetable); ok {
				t.Damage(a.values.Damage)
			}
		}
		if a.interX[i] != nil && a.interY[i] != nil {
			x, finX := a.interX[i].Update(timing.DT)
			y, finY := a.interY[i].Update(timing.DT)
			target.SetXY(x, y)
			if finX && finY {
				floor.CurrentFloor.MoveOccupant(target, target.GetCoords(), a.ends[i])
				target.SetCoords(a.ends[i])
				a.interX[i] = nil
				a.interY[i] = nil
			} else {
				done = false
			}
		}
	}
	if a.start {
		a.start = false
	}
	if done {
		a.isDone = true
	}
}

func (a *PushMultiAction) IsDone() bool {
	return a.isDone
}