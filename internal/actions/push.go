package actions

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/action"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/sfx"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type PushAction struct {
	*action.AbstractAction
	source *floor.Character
	target *floor.Character
	push   int
	start  bool
}

func NewPushAction(s, t *floor.Character, p int) *PushAction {
	if t == nil {
		return nil
	}
	return &PushAction{
		source: s,
		target: t,
		push:   p,
		start:  true,
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
			a.IsDone = true
			return
		}
		action.ActionManager.AddToTop(NewMoveAction(a.source, a.target, nPath[len(nPath)-1]), nil)
		a.start = false
		a.IsDone = true
	}
}

func (a *PushAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}

type PushMultiAction struct {
	*action.AbstractAction
	area    []world.Coords
	values  selector.ActionValues
	start   bool
	preDam  bool
	postDam bool
	interX  *gween.Tween
	interY  *gween.Tween
	targets []*floor.Character
	ends    []world.Coords
	interXP []*gween.Tween
	interYP []*gween.Tween
}

func NewPushMultiAction(area []world.Coords, values selector.ActionValues) *PushMultiAction {
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
		SetAttackTransform(a.values.Source, a.area)
		a.start = false
	}
	if !a.values.Source.IsMoving() {
		if a.preDam {
			orig := a.values.Source.GetCoords()
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
					o := orig
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
		}
		a.preDam = false
	}
	if !a.preDam {
		if !a.postDam {
			SetResetTransform(a.values.Source)
		}

		first := true
		done := !a.values.Source.IsMoving() && a.postDam
		for i, target := range a.targets {
			if !a.postDam {
				if first {
					sfx.SoundPlayer.PlaySound("punch_hit")
					first = false
				}
				target.Damage(a.values.Damage)
			}
			if a.interXP[i] != nil && a.interYP[i] != nil {
				x, finX := a.interXP[i].Update(timing.DT)
				y, finY := a.interYP[i].Update(timing.DT)
				target.SetPos(pixel.V(x, y))
				if finX && finY {
					floor.CurrentFloor.MoveOccupant(target, target.GetCoords(), a.ends[i])
					a.interXP[i] = nil
					a.interYP[i] = nil
				} else {
					done = false
				}
			}
		}
		if !a.postDam {
			if first {
				sfx.SoundPlayer.PlaySound("punch_miss")
			}
			a.postDam = true
		}
		if done {
			a.IsDone = true
		}
	}
}

func (a *PushMultiAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}
