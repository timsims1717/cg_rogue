package debug

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
)

var (
	imd         *imdraw.IMDraw
)

func InitializeLines() {
	imd = imdraw.New(nil)
}

func UpdateLines() {
	imd.Clear()
	// pc to coords
	imd.Color = colornames.Red
	imd.EndShape = imdraw.NoEndShape
	imd.Push(pixel.V(player.Player1.Character.GetXY()), pixel.V(world.MapToWorld(player.Player1.Input.Coords.X, player.Player1.Input.Coords.Y)))
	imd.Line(2)
	// pc to mouse
	imd.Color = colornames.Green
	imd.EndShape = imdraw.NoEndShape
	imd.Push(pixel.V(player.Player1.Character.GetXY()), player.Player1.Input.World)
	imd.Line(2)
	// path line
	imd.Color = colornames.Blue
	imd.EndShape = imdraw.NoEndShape
	imd.Push(floor.CurrentFloor.PathLine.A, floor.CurrentFloor.PathLine.B)
	imd.Line(2)
}

func DrawLines(win *pixelgl.Window) {
	imd.Draw(win)
}