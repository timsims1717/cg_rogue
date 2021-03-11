package game

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	ui2 "github.com/timsims1717/cg_rogue_go/internal/ui"
	"golang.org/x/image/colornames"
)

var (
	CenterText *ui2.ActionEl
)

func Initialize() {
	text := ui2.NewActionText("")
	text.Align = ui2.Center
	text.VAlign = ui2.Center
	text.Scalar = pixel.V(6., 6.)
	text.TextColor = colornames.Black
	CenterText = ui2.NewActionEl(text, pixel.R(0, 0, cfg.WindowWidthF, cfg.WindowHeightF))
	CenterText.UI = true
	CenterText.Pos = pixel.V(cfg.WindowWidthF * 0.5, cfg.WindowHeightF * 0.5)
}