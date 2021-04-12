package state

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	ui2 "github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/sfx"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
)

type MainMenu struct {}

var (
	cameraPan = []world.Coords{
		{
			33, 40,
		},
		{
			28, 65,
		},
		{
			35, 68,
		},
		{
			48, 72,
		},
		{
			60, 65,
		},
		{
			69, 40,
		},
		{
			66, 33,
		},
		{
			42, 32,
		},
	}
	panIndex = 0
)

var (
	startMenuItem *ui2.ActionEl
	exitMenuItem  *ui2.ActionEl
)

func (s *MainMenu) Initialize() {
	spritesheet, err := img.LoadSpriteSheet("assets/img/testfloor.json")
	if err != nil {
		panic(err)
	}
	floor.DefaultFloor(100, 100, spritesheet)
	player.Player1.Hand = nil
	player.Player1.PlayCard = nil
	player.Player1.Discard = nil
	player.Player1.Grid = nil
	player.Player1.Input.RemoveHotKeys()
	camera.Cam.CenterOn([]pixel.Vec{world.MapToWorld(cameraPan[0])})
	start := "Start Game"
	startText := ui2.NewActionText(start)
	startText.Transform.Scalar = pixel.V(4., 4.)
	startText.TextColor = colornames.Purple
	startR := pixel.R(0., 0., startText.Text.BoundsOf(start).W() * 4., startText.Text.BoundsOf(start).H() * 3.75)
	startMenuItem = ui2.NewActionEl(startText, startR, true)
	startMenuItem.Show = true
	startMenuItem.Transform.Pos = pixel.V(20., 280.)
	startMenuItem.SetOnHoverFn(func() {
		startMenuItem.Text.TextColor = colornames.Forestgreen
		startMenuItem.Transform.Scalar = pixel.V(1.05, 1.05)
		startMenuItem.Transform.Pos.X -= 5.
		sfx.SoundPlayer.PlaySound("click")
	})
	startMenuItem.SetUnHoverFn(func() {
		startMenuItem.Text.TextColor = colornames.Purple
		startMenuItem.Transform.Scalar = pixel.V(1., 1.)
		startMenuItem.Transform.Pos.X += 5.
	})
	startMenuItem.SetClickFn(func() {
		camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.Black, 1.0)
		sfx.MusicPlayer.FadeOut(1.0)
		SwitchState(TheStartRun)
	})
	exitS := "Exit"
	exitText := ui2.NewActionText(exitS)
	exitText.Transform.Scalar = pixel.V(4., 4.)
	exitText.TextColor = colornames.Purple
	exitR := pixel.R(0., 0., exitText.Text.BoundsOf(exitS).W() * 4., exitText.Text.BoundsOf(exitS).H() * 3.75)
	exitMenuItem = ui2.NewActionEl(exitText, exitR, true)
	exitMenuItem.Show = true
	exitMenuItem.Transform.Pos = pixel.V(20., 220.)
	exitMenuItem.SetOnHoverFn(func() {
		exitMenuItem.Text.TextColor = colornames.Forestgreen
		exitMenuItem.Transform.Scalar = pixel.V(1.05, 1.05)
		exitMenuItem.Transform.Pos.X -= 2.
		sfx.SoundPlayer.PlaySound("click")
	})
	exitMenuItem.SetUnHoverFn(func() {
		exitMenuItem.Text.TextColor = colornames.Purple
		exitMenuItem.Transform.Scalar = pixel.V(1., 1.)
		exitMenuItem.Transform.Pos.X += 2.
	})
	exitMenuItem.SetClickFn(func() {
		camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.Black, 0.2)
		sfx.MusicPlayer.FadeOut(0.5)
		SwitchState(&Exiting{})
	})
	camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.White, 1.0)
	sfx.MusicPlayer.SetCurrentTracks([]string{"main_menu"})
	sfx.MusicPlayer.PlayNextTrack(0.0, 0.5, 0., false)
}

func (s *MainMenu) TransitionIn() bool {
	return camera.Cam.Effect != nil
}

func (s *MainMenu) TransitionOut() bool {
	return camera.Cam.Effect != nil
}

func (s *MainMenu) Uninitialize() {
	camera.Cam.Stop()
	floor.CurrentFloor = nil
}

func (s *MainMenu) Update(win *pixelgl.Window) {
	player.Player1.Input.Update(win)
	if !camera.Cam.Moving() {
		panIndex += 1
		panIndex = panIndex % 8
		camera.Cam.MoveTo(world.MapToWorld(cameraPan[panIndex]), 10., true)
	}
	camera.Cam.Update(win)
	startMenuItem.Update(player.Player1.Input)
	exitMenuItem.Update(player.Player1.Input)
}

func (s *MainMenu) Draw(win *pixelgl.Window) {
	floor.CurrentFloor.Draw(win)
	startMenuItem.Draw(win)
	exitMenuItem.Draw(win)
}

func (s *MainMenu) String() string {
	return "MainMenu"
}