package debug

import (
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
	if player.Player1.Character != nil {
		// pc to coords
		imd.Color = colornames.Red
		imd.EndShape = imdraw.NoEndShape
		imd.Push(player.Player1.Character.GetPos(), world.MapToWorld(player.Player1.Input.Coords))
		imd.Line(2)
		// pc to mouse
		imd.Color = colornames.Green
		imd.EndShape = imdraw.NoEndShape
		imd.Push(player.Player1.Character.GetPos(), player.Player1.Input.World)
		imd.Line(2)
	}
	if floor.CurrentFloor != nil {
		// path line
		imd.Color = colornames.Blue
		imd.EndShape = imdraw.NoEndShape
		imd.Push(floor.CurrentFloor.PathLine.A, floor.CurrentFloor.PathLine.B)
		imd.Line(2)
	}
}

func DrawLines(win *pixelgl.Window) {
	imd.Draw(win)
}