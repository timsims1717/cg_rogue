package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type MoveAction struct{
	c      *characters.Character
	start  world.Coords
	end    world.Coords
	isDone bool
	interX *gween.Tween
	interY *gween.Tween
}

func NewMoveAction(c *characters.Character, end world.Coords) *MoveAction {
	if end.Equals(c.Coords) {
		return nil
	}
	bx, by := world.MapToWorldHex(end.X, end.Y)
	return &MoveAction{
		c:      c,
		end:    end,
		start:  c.Coords,
		isDone: false,
		interX: gween.New(c.Pos.X, bx, 0.25, ease.InOutQuad),
		interY: gween.New(c.Pos.Y, by, 0.25, ease.InOutQuad),
	}
}

func (m *MoveAction) Update() {
	x, finX := m.interX.Update(timing.DT)
	y, finY := m.interY.Update(timing.DT)
	m.c.Pos.X = x
	m.c.Pos.Y = y
	if finX && finY {
		m.c.Coords = m.end
		m.isDone = true
		floor.CurrentFloor.MoveOccupant(m.c, m.start, m.end)
	}
}

func (m *MoveAction) IsDone() bool {
	return m.isDone
}

type MoveSeriesAction struct{
	c      *characters.Character
	series []world.Coords
	step   int
	start  world.Coords
	isDone bool
	interX *gween.Tween
	interY *gween.Tween
}

func NewMoveSeriesAction(c *characters.Character, series []world.Coords) *MoveSeriesAction {
	if len(series) == 0 {
		return nil
	} else {
		first := series[0]
		bx, by := world.MapToWorldHex(first.X, first.Y)

		return &MoveSeriesAction{
			c:      c,
			series: series,
			step:   0,
			start: c.Coords,
			isDone: false,
			interX: gween.New(c.Pos.X, bx, 0.25, ease.InQuad),
			interY: gween.New(c.Pos.Y, by, 0.25, ease.InQuad),
		}
	}
}

func (m *MoveSeriesAction) Update() {
	x, finX := m.interX.Update(timing.DT)
	y, finY := m.interY.Update(timing.DT)
	m.c.Pos.X = x
	m.c.Pos.Y = y
	if finX && finY {
		if m.step >= len(m.series) - 1 {
			next := m.series[m.step]
			m.c.Coords = next
			m.isDone = true
			floor.CurrentFloor.MoveOccupant(m.c, m.start, next)
		} else {
			next := m.series[m.step + 1]
			bx, by := world.MapToWorldHex(next.X, next.Y)
			if m.step >= len(m.series) - 2 {
				m.interX = gween.New(m.c.Pos.X, bx, 0.25, ease.OutQuad)
				m.interY = gween.New(m.c.Pos.Y, by, 0.25, ease.OutQuad)
			} else {
				m.interX = gween.New(m.c.Pos.X, bx, 0.15, ease.Linear)
				m.interY = gween.New(m.c.Pos.Y, by, 0.15, ease.Linear)
			}
		}
		m.step++
	}
}

func (m *MoveSeriesAction) IsDone() bool {
	return m.isDone
}