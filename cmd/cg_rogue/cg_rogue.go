package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/debug"
	"github.com/timsims1717/cg_rogue_go/internal/state"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/sfx"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

func run() {
	//seed := int64(1617372779626977911)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	fmt.Println("Seed:", seed)
	world.ScaledTileSize = cfg.ScaledTileSize
	camera.SetWindowSize(1600, 900)
	config := pixelgl.WindowConfig{
		Title:  cfg.Title,
		Bounds: pixel.R(0, 0, camera.WindowWidthF, camera.WindowHeightF),
		//VSync: true,
	}
	win, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}

	camera.Cam = camera.New()

	sfx.MusicPlayer.RegisterMusicTrack("assets/music/test_track.MP3", "test_track")
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/test_ambience.MP3", "test_ambience")

	debug.Initialize()

	timing.Reset()
	for !win.Closed() {
		timing.Update()
		debug.Update()
		state.Update(win)
		sfx.MusicPlayer.Update()

		win.Clear(colornames.Black)

		state.Draw(win)
		debug.Draw(win)
		win.Update()
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}