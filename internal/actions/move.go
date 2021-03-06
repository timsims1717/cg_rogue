package actions

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/action"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/sfx"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type MoveAction struct {
	*action.AbstractAction
	source *floor.Character
	target *floor.Character
	start  world.Coords
	end    world.Coords
	interX *gween.Tween
	interY *gween.Tween
	put    bool
}

func NewMoveAction(source *floor.Character, target *floor.Character, end world.Coords) *MoveAction {
	if end.Eq(target.GetCoords()) {
		return nil
	}
	b := world.MapToWorld(end)
	p := target.GetPos()
	return &MoveAction{
		source: source,
		target: target,
		end:    end,
		start:  target.GetCoords(),
		interX: gween.New(p.X, b.X, 0.25, ease.InOutQuad),
		interY: gween.New(p.Y, b.Y, 0.25, ease.InOutQuad),
	}
}

func (a *MoveAction) Update() {
	if a.target.Coords.Above(a.end) && !a.put {
		floor.CurrentFloor.PutOccupant(a.target, a.end)
		a.put = true
	}
	x, finX := a.interX.Update(timing.DT)
	y, finY := a.interY.Update(timing.DT)
	a.target.SetPos(pixel.V(x, y))
	if finX && finY {
		if !a.put {
			floor.CurrentFloor.PutOccupant(a.target, a.end)
		}
		sfx.SoundPlayer.PlaySound("step1")
		a.IsDone = true
	}
}

func (a *MoveAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}

type MoveSeriesAction struct {
	*action.AbstractAction
	source *floor.Character
	target *floor.Character
	series []world.Coords
	step   int
	put    bool
	start  world.Coords
	interX *gween.Tween
	interY *gween.Tween
}

func NewMoveSeriesAction(source *floor.Character, target *floor.Character, series []world.Coords) *MoveSeriesAction {
	if len(series) == 0 {
		return nil
	} else {
		first := series[0]
		b := world.MapToWorld(first)
		p := target.GetPos()

		return &MoveSeriesAction{
			target: target,
			series: series,
			step:   0,
			start:  target.GetCoords(),
			interX: gween.New(p.X, b.X, 0.25, ease.InQuad),
			interY: gween.New(p.Y, b.Y, 0.25, ease.InQuad),
		}
	}
}

func (a *MoveSeriesAction) Update() {
	next := a.series[a.step]
	if a.target.Coords.Above(next) && !a.put {
		floor.CurrentFloor.PutOccupant(a.target, a.series[a.step])
		a.put = true
	}
	x, finX := a.interX.Update(timing.DT)
	y, finY := a.interY.Update(timing.DT)
	a.target.SetPos(pixel.V(x, y))
	if finX && finY {
		if a.step >= len(a.series)-1 {
			a.IsDone = true
			if !a.put {
				floor.CurrentFloor.PutOccupant(a.target, next)
			}
		} else {
			sfx.SoundPlayer.PlaySound("step1")
			if !a.put {
				floor.CurrentFloor.PutOccupant(a.target, next)
			}
			a.put = false
			b := world.MapToWorld(a.series[a.step+1])
			if a.step >= len(a.series)-2 {
				a.interX = gween.New(x, b.X, 0.25, ease.OutQuad)
				a.interY = gween.New(y, b.Y, 0.25, ease.OutQuad)
			} else {
				a.interX = gween.New(x, b.X, 0.15, ease.Linear)
				a.interY = gween.New(y, b.Y, 0.15, ease.Linear)
			}
		}
		a.step++
	}
}

func (a *MoveSeriesAction) SetAbstract(abstractAction *action.AbstractAction) {
	a.AbstractAction = abstractAction
}
