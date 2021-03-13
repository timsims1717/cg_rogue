package game

import (
	"github.com/faiface/pixel"
	ui2 "github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"golang.org/x/image/colornames"
)

var (
	CenterText *ui2.ActionEl
)

func Initialize() {
	text := ui2.NewActionText("")
	text.Align = ui2.Center
	text.Transform.Anchor.H = animation.Center
	text.Transform.Anchor.V = animation.Center
	text.Transform.Scalar = pixel.V(6., 6.)
	text.TextColor = colornames.Black
	CenterText = ui2.NewActionEl(text, pixel.R(0, 0, camera.WindowWidthF, camera.WindowHeightF), camera.Cam)
}