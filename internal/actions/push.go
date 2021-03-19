package actions

import (
	"github.com/faiface/pixel"
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
			next := world.NextHexLine(o, n)
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
	preDam  bool
	postDam bool
	isDone  bool
	interX  *gween.Tween
	interY  *gween.Tween
	targets []objects.Moveable
	ends    []world.Coords
	interXP []*gween.Tween
	interYP []*gween.Tween
}

func NewPushMultiAction(area []world.Coords, values selectors.ActionValues) *PushMultiAction {
	if len(area) < 1 {
		return nil
	}
	return &PushMultiAction{
		area:   area,
		values: values,
		start:  true,
		preDam: true,
	}
}

func (a *PushMultiAction) Update() {
	if a.start {
		SetAttackTransform(a.values.Source, a.area[0])
		a.start = false
	}
	if !a.values.Source.IsMoving() {
		a.preDam = false
	}
	if !a.preDam {
		if !a.postDam {
			SetResetTransform(a.values.Source)
		}
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
					next := world.NextHexLine(o, n)
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
					a.interXP = append(a.interXP, nil)
					a.interYP = append(a.interYP, nil)
				} else {
					end := nPath[len(nPath)-1]
					b := world.MapToWorld(end)
					a.ends = append(a.ends, end)
					p := target.GetPos()
					a.interXP = append(a.interXP, gween.New(p.X, b.X, 0.25, ease.OutCubic))
					a.interYP = append(a.interYP, gween.New(p.Y, b.Y, 0.25, ease.OutCubic))
				}
			}
		}

		done := !a.values.Source.IsMoving() && a.postDam
		for i, target := range a.targets {
			if !a.postDam {
				if t, ok := target.(objects.Targetable); ok {
					t.Damage(a.values.Damage)
				}
			}
			if a.interXP[i] != nil && a.interYP[i] != nil {
				x, finX := a.interXP[i].Update(timing.DT)
				y, finY := a.interYP[i].Update(timing.DT)
				target.SetPos(pixel.V(x, y))
				if finX && finY {
					floor.CurrentFloor.MoveOccupant(target, target.GetCoords(), a.ends[i])
					target.SetCoords(a.ends[i])
					a.interXP[i] = nil
					a.interYP[i] = nil
				} else {
					done = false
				}
			}
		}
		if !a.postDam {
			a.postDam = true
		}
		if done {
			a.isDone = true
		}
	}
}

func (a *PushMultiAction) IsDone() bool {
	return a.isDone
}