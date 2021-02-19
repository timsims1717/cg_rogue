package input

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type Input struct {
	Cursor pixel.Vec
	World  pixel.Vec
	Coords world.Coords
	Select toggle
	Cancel toggle
	Move   toggle
	Attack toggle
}

func (i *Input) Update(win *pixelgl.Window) {
	i.Cursor = win.MousePosition()
	i.World = camera.Cam.Mat.Unproject(win.MousePosition())
	i.Coords.X, i.Coords.Y = world.WorldToMapHex(i.World.X, i.World.Y)
	i.Select.Set(win, pixelgl.MouseButtonLeft)
	i.Cancel.Set(win, pixelgl.MouseButtonRight)

	if win.Pressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}
	if win.Pressed(pixelgl.KeyLeft) {
		camera.Cam.Left()
	}
	if win.Pressed(pixelgl.KeyRight) {
		camera.Cam.Right()
	}
	if win.Pressed(pixelgl.KeyDown) {
		camera.Cam.Down()
	}
	if win.Pressed(pixelgl.KeyUp) {
		camera.Cam.Up()
	}
	camera.Cam.ZoomIn(win.MouseScroll().Y)
}

type toggle struct {
	justPressed  bool
	pressed      bool
	justReleased bool
	consumed     bool
}

func (t *toggle) JustPressed() bool {
	return t.justPressed && !t.consumed
}

func (t *toggle) Pressed() bool {
	return t.pressed && !t.consumed
}

func (t *toggle) JustReleased() bool {
	return t.justReleased && !t.consumed
}

func (t *toggle) Consume() {
	t.consumed = true
}

func (t *toggle) Set(win *pixelgl.Window, button pixelgl.Button) {
	t.justPressed = win.JustPressed(button)
	t.pressed = win.Pressed(button)
	t.justReleased = win.JustReleased(button)
	t.consumed = t.consumed && !t.justPressed && !t.pressed && !t.justReleased
}