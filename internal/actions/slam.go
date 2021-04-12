package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/sfx"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type SlamAction struct {
	values  selector.ActionValues
	area    []world.Coords
	landing world.Coords
	start   bool
	isDone  bool
}

func NewSlamAction(landing world.Coords, area []world.Coords, values selector.ActionValues) *SlamAction {
	if len(area) < 1 {
		return nil
	}
	return &SlamAction{
		values:  values,
		area:    area,
		landing: landing,
		start:   true,
		isDone:  false,
	}
}

func (a *SlamAction) Update() {
	if a.start {
		p := a.values.Source.GetPos()
		e := world.MapToWorld(a.landing)
		transform := animation.TransformBuilder{
			Transform: a.values.Source.GetTransform(),
			InterX:    gween.New(p.X, e.X, 0.25, ease.InOutQuad),
			InterY:    gween.New(p.Y, e.Y, 0.25, ease.InOutQuad),
		}
		a.values.Source.SetTransformEffect(transform.Build())
		a.start = false
	}
	if !a.values.Source.IsMoving() {
		a.isDone = true
		floor.CurrentFloor.MoveOccupant(a.values.Source, a.values.Source.Coords, a.landing)
		first := true
		for _, h := range a.area {
			// todo: add an effect

			// todo: this is where the damage modification happens
			if cha := floor.CurrentFloor.GetOccupant(h); cha != nil {
				if first {
					sfx.SoundPlayer.PlaySound("punch_hit")
					first = false
				}
				cha.Damage(a.values.Damage)
			}
		}
		if first {
			sfx.SoundPlayer.PlaySound("punch_miss")
		}
	}
}

func (a *SlamAction) IsDone() bool {
	return a.isDone
}