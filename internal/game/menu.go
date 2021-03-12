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
	s := "Start Game"
	text := ui2.NewActionText(s)
	text.VAlign = ui2.Right
	text.Scalar = pixel.V(4., 4.)
	text.TextColor = colornames.Purple
	r := pixel.R(0., 0., text.Text.BoundsOf(s).W() * 4., text.Text.BoundsOf(s).H() * 3.75)
	startMenuItem = ui2.NewActionEl(text, r)
	startMenuItem.Show = true
	startMenuItem.UI = true
	startMenuItem.Pos = pixel.V(20., 280.)
	startMenuItem.SetOnHoverFn(func() {
		startMenuItem.T.TextColor = colornames.Forestgreen
		startMenuItem.Scalar = pixel.V(1.05, 1.05)
		startMenuItem.Pos.X += 5.
	})
	startMenuItem.SetUnHoverFn(func() {
		startMenuItem.T.TextColor = colornames.Purple
		startMenuItem.Scalar = pixel.V(1., 1.)
		startMenuItem.Pos.X -= 5.
	})
	startMenuItem.SetClickFn(func() {
		SwitchState(state.InGame)
	})
	exitS := "Exit"
	exitText := ui2.NewActionText(exitS)
	exitText.VAlign = ui2.Right
	exitText.Scalar = pixel.V(4., 4.)
	exitText.TextColor = colornames.Purple
	exitR := pixel.R(0., 0., exitText.Text.BoundsOf(exitS).W() * 4., exitText.Text.BoundsOf(exitS).H() * 3.75)
	exitMenuItem = ui2.NewActionEl(exitText, exitR)
	exitMenuItem.Show = true
	exitMenuItem.UI = true
	exitMenuItem.Pos = pixel.V(20., 220.)
	exitMenuItem.SetOnHoverFn(func() {
		exitMenuItem.T.TextColor = colornames.Forestgreen
		exitMenuItem.Scalar = pixel.V(1.05, 1.05)
		exitMenuItem.Pos.X += 5.
	})
	exitMenuItem.SetUnHoverFn(func() {
		exitMenuItem.T.TextColor = colornames.Purple
		exitMenuItem.Scalar = pixel.V(1., 1.)
		exitMenuItem.Pos.X -= 5.
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