package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type MoveAction struct{
	source *characters.Character
	target objects.Moveable
	start  world.Coords
	end    world.Coords
	isDone bool
	interX *gween.Tween
	interY *gween.Tween
}

func NewMoveAction(source *characters.Character, target objects.Moveable, end world.Coords) *MoveAction {
	if end.Equals(target.GetCoords()) {
		return nil
	}
	b := world.MapToWorld(end)
	px, py := target.GetXY()
	return &MoveAction{
		source: source,
		target: target,
		end:    end,
		start:  target.GetCoords(),
		isDone: false,
		interX: gween.New(px, b.X, 0.25, ease.InOutQuad),
		interY: gween.New(py, b.Y, 0.25, ease.InOutQuad),
	}
}

func (a *MoveAction) Update() {
	x, finX := a.interX.Update(timing.DT)
	y, finY := a.interY.Update(timing.DT)
	a.target.SetXY(x, y)
	if finX && finY {
		a.target.SetCoords(a.end)
		a.isDone = true
		floor.CurrentFloor.MoveOccupant(a.target, a.start, a.end)
	}
}

func (a *MoveAction) IsDone() bool {
	return a.isDone
}

type MoveSeriesAction struct{
	source *characters.Character
	target objects.Moveable
	series []world.Coords
	step   int
	start  world.Coords
	isDone bool
	interX *gween.Tween
	interY *gween.Tween
}

func NewMoveSeriesAction(source *characters.Character, target objects.Moveable, series []world.Coords) *MoveSeriesAction {
	if len(series) == 0 {
		return nil
	} else {
		first := series[0]
		b := world.MapToWorld(first)
		px, py := target.GetXY()

		return &MoveSeriesAction{
			target: target,
			series: series,
			step:   0,
			start:  target.GetCoords(),
			isDone: false,
			interX: gween.New(px, b.X, 0.25, ease.InQuad),
			interY: gween.New(py, b.Y, 0.25, ease.InQuad),
		}
	}
}

func (m *MoveSeriesAction) Update() {
	x, finX := m.interX.Update(timing.DT)
	y, finY := m.interY.Update(timing.DT)
	m.target.SetXY(x, y)
	if finX && finY {
		if m.step >= len(m.series) - 1 {
			next := m.series[m.step]
			m.target.SetCoords(next)
			m.isDone = true
			floor.CurrentFloor.MoveOccupant(m.target, m.start, next)
		} else {
			next := m.series[m.step + 1]
			b := world.MapToWorld(next)
			if m.step >= len(m.series) - 2 {
				m.interX = gween.New(x, b.X, 0.25, ease.OutQuad)
				m.interY = gween.New(y, b.Y, 0.25, ease.OutQuad)
			} else {
				m.interX = gween.New(x, b.X, 0.15, ease.Linear)
				m.interY = gween.New(y, b.Y, 0.15, ease.Linear)
			}
		}
		m.step++
	}
}

func (m *MoveSeriesAction) IsDone() bool {
	return m.isDone
}