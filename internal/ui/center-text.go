package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/ui"
	"golang.org/x/image/colornames"
)

var (
	CenterText *ui.TextDisplay
)

func Initialize() {
	CenterText = &ui.TextDisplay{
		Text:      text.New(pixel.ZV, ui.BasicAtlas),
		Raw:       "",
		Show:      false,
		UI:        true,
		Align:     ui.Center,
		Pos:       pixel.V(cfg.WindowWidthF * 0.5, cfg.WindowHeightF * 0.5),
		Scalar:    pixel.V(6., 6.),
		TextColor: colornames.Black,
	}
}

func Update() {
	CenterText.Update()
}

func Draw(win *pixelgl.Window) {
	CenterText.Draw(win)
}