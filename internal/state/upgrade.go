package state

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/run"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"golang.org/x/image/colornames"
)

type Upgrade struct {
	Grid *player.Grid
	Done bool
}

func (s *Upgrade) Initialize() {
	s.Done = false
	run.CurrentRun.Level++
	s.Grid = player.NewGrid(player.Player1)
	s.Grid.Show = true
	s.Grid.Group = player.Player1.GetDeck()
	camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.White, 1.0)
}

func (s *Upgrade) TransitionIn() bool {
	return camera.Cam.Effect != nil
}

func (s *Upgrade) TransitionOut() bool {
	return camera.Cam.Effect != nil
}

func (s *Upgrade) Uninitialize() {
	s.Grid.ReturnCards()
	s.Grid = nil
}

func (s *Upgrade) Update(win *pixelgl.Window) {
	player.Player1.Input.Update(win)
	camera.Cam.Update(win)
	s.Grid.Update()
	if !s.Done {
		for _, card := range s.Grid.Group {
			card.Action.SetValues(card.Level)
		}
		if s.Grid.Hovered > -1 && s.Grid.Hovered < len(s.Grid.Group) {
			card := s.Grid.Group[s.Grid.Hovered]
			card.Action.SetValues(card.Level + 1)
		}
	}
	if !s.Done && len(s.Grid.Clicked) > 0 {
		for _, card := range s.Grid.Clicked {
			card.Upgrade()
		}
		s.Done = true
		camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.Black, 1.0)
		SwitchState(TheEncounter)
	}
}

func (s *Upgrade) Draw(win *pixelgl.Window) {
	win.SetSmooth(true)
	s.Grid.Draw(win)
	win.SetSmooth(false)
}

func (s *Upgrade) String() string {
	return "Upgrade"
}