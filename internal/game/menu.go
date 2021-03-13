package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/state"
	ui2 "github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
)

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

func InitializeMenu(win *pixelgl.Window) {
	spritesheet, err := img.LoadSpriteSheet("assets/img/testfloor.json")
	if err != nil {
		panic(err)
	}
	floor.DefaultFloor(100, 100, spritesheet)
	player.Player1 = player.NewPlayer(nil)
	player.Player1.Hand = nil
	player.Player1.PlayCard = nil
	player.Player1.Discard = nil
	camera.Cam.CenterOn([]pixel.Vec{world.MapToWorld(cameraPan[0])})
	start := "Start Game"
	startText := ui2.NewActionText(start)
	startText.Transform.Scalar = pixel.V(4., 4.)
	startText.TextColor = colornames.Purple
	startR := pixel.R(0., 0., startText.Text.BoundsOf(start).W() * 4., startText.Text.BoundsOf(start).H() * 3.75)
	startMenuItem = ui2.NewActionEl(startText, startR, camera.Cam)
	startMenuItem.Show = true
	startMenuItem.Transform.Pos = pixel.V(20., 280.)
	startMenuItem.SetOnHoverFn(func() {
		startMenuItem.Text.TextColor = colornames.Forestgreen
		startMenuItem.Transform.Scalar = pixel.V(1.05, 1.05)
		startMenuItem.Transform.Pos.X -= 5.
	})
	startMenuItem.SetUnHoverFn(func() {
		startMenuItem.Text.TextColor = colornames.Purple
		startMenuItem.Transform.Scalar = pixel.V(1., 1.)
		startMenuItem.Transform.Pos.X += 5.
	})
	startMenuItem.SetClickFn(func() {
		SwitchState(state.InGame)
	})
	exitS := "Exit"
	exitText := ui2.NewActionText(exitS)
	exitText.Transform.Scalar = pixel.V(4., 4.)
	exitText.TextColor = colornames.Purple
	exitR := pixel.R(0., 0., exitText.Text.BoundsOf(exitS).W() * 4., exitText.Text.BoundsOf(exitS).H() * 3.75)
	exitMenuItem = ui2.NewActionEl(exitText, exitR, camera.Cam)
	exitMenuItem.Show = true
	exitMenuItem.Transform.Pos = pixel.V(20., 220.)
	exitMenuItem.SetOnHoverFn(func() {
		exitMenuItem.Text.TextColor = colornames.Forestgreen
		exitMenuItem.Transform.Scalar = pixel.V(1.05, 1.05)
		exitMenuItem.Transform.Pos.X -= 2.
	})
	exitMenuItem.SetUnHoverFn(func() {
		exitMenuItem.Text.TextColor = colornames.Purple
		exitMenuItem.Transform.Scalar = pixel.V(1., 1.)
		exitMenuItem.Transform.Pos.X += 2.
	})
	exitMenuItem.SetClickFn(func() {
		win.SetClosed(true)
	})
}

func TransitionInMenu() bool {
	return false
}

func TransitionOutMenu() bool {
	return false
}

func UninitializeMenu() {
	camera.Cam.Stop()
	floor.CurrentFloor = nil
}

func UpdateMenu(win *pixelgl.Window) {
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

func DrawMenu(win *pixelgl.Window) {
	floor.CurrentFloor.Draw(win)
	startMenuItem.Draw(win)
	exitMenuItem.Draw(win)
}