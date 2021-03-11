package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/debug"
	"github.com/timsims1717/cg_rogue_go/internal/game"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

func run() {
	rand.Seed(time.Now().UnixNano())
	world.ScaledTileSize = cfg.ScaledTileSize
	cfg.SetWindowSize(1600, 900)
	config := pixelgl.WindowConfig{
		Title:  cfg.Title,
		Bounds: pixel.R(0, 0, cfg.WindowWidthF, cfg.WindowHeightF),
		//VSync: true,
	}
	win, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}

	camera.Cam = camera.New()

	debug.Initialize()

	timing.Reset()
	for !win.Closed() {
		timing.Update()
		debug.Update()
		game.Update(win)

		win.Clear(colornames.Forestgreen)

		game.Draw(win)
		debug.Draw(win)
		win.Update()
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}