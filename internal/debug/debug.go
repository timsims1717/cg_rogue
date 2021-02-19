package debug

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	text2 "github.com/timsims1717/cg_rogue_go/pkg/text"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"time"
)

var (
	// The canvas where the debug information is drawn
	canvas *pixelgl.Canvas
	// The number of frames since a second has passed
	frames = 0
	second = time.Tick(time.Second)
	// The text containers for fps, worlds coordinates, and maps coordinates
	fps    *text.Text
	mouse  *text.Text
	worlds *text.Text
	maps   *text.Text
	health *text.Text
	dispHP bool
	cards  *text.Text
)

// Initialize creates the debug canvas and all the text containers.
// This is where the location of the text containers is set.
func Initialize() {
	canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(cfg.WindowWidth), float64(cfg.WindowHeight)))
	fps = text.New(pixel.V(20., float64(cfg.WindowHeight) - text2.BasicAtlas.LineHeight() - 20.), text2.BasicAtlas)
	mouse = text.New(pixel.V(20., float64(cfg.WindowHeight) - (text2.BasicAtlas.LineHeight() + 20.) * 2.0), text2.BasicAtlas)
	worlds = text.New(pixel.V(20., float64(cfg.WindowHeight) - (text2.BasicAtlas.LineHeight() + 20.) * 3.0), text2.BasicAtlas)
	maps = text.New(pixel.V(20., float64(cfg.WindowHeight) - (text2.BasicAtlas.LineHeight() + 20.) * 4.0), text2.BasicAtlas)
	health = text.New(pixel.V(20., float64(cfg.WindowHeight) - (text2.BasicAtlas.LineHeight() + 20.) * 5.0), text2.BasicAtlas)
	cards = text.New(pixel.V(20., float64(cfg.WindowHeight) - (text2.BasicAtlas.LineHeight() + 20.) * 6.0), text2.BasicAtlas)
}

// Update clears the text containers and updates them with the correct information.
func Update(win *pixelgl.Window) {
	frames++
	select {
	case <-second:
		fps.Clear()
		fmt.Fprintf(fps, "FPS: %d", frames)
		frames = 0
	default:
	}
	mousePtr := win.MousePosition()
	wrldPtr := camera.Cam.Mat.Unproject(mousePtr)
	mapX, mapY := world.WorldToMapHex(wrldPtr.X, wrldPtr.Y)
	mouse.Clear()
	worlds.Clear()
	maps.Clear()
	fmt.Fprintf(mouse, "Mouse (X,Y): (%d,%d)", int(mousePtr.X), int(mousePtr.Y))
	fmt.Fprintf(worlds, "World (X,Y): (%d,%d)", int(wrldPtr.X), int(wrldPtr.Y))
	fmt.Fprintf(maps, "Map (X,Y): (%d,%d)", mapX, mapY)
	dispHP = false
	health.Clear()
	occ := floor.CurrentFloor.GetOccupant(world.Coords{mapX, mapY})
	if objects.NotNil(occ) {
		if cha, ok := occ.(*characters.Character); ok {
			fmt.Fprintf(health, "Health: %d/%d HP", cha.CurrHP, cha.MaxHP)
			dispHP = true
		}
	}
	cards.Clear()
	fmt.Fprintf(cards, "Player 1 Cards: (Hovered: %d)\n", player.Player1.Hand.Hovered)
	for i, card := range player.Player1.Hand.Group {
		hovered := card.PointInside(wrldPtr)
		fmt.Fprintf(cards, "   %s (%d): Hovered: %t\n", card.RawTitle, i, hovered)
	}
}

// Draw draws each text container to the canvas.
// This is where scaling happens.
func Draw(win *pixelgl.Window) {
	canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
	fps.Draw(canvas, pixel.IM.Scaled(fps.Orig, 2.))
	mouse.Draw(canvas, pixel.IM.Scaled(mouse.Orig, 2.))
	worlds.Draw(canvas, pixel.IM.Scaled(worlds.Orig, 2.))
	maps.Draw(canvas, pixel.IM.Scaled(maps.Orig, 2.))
	if dispHP {
		health.Draw(canvas, pixel.IM.Scaled(health.Orig, 2.))
	}
	cards.Draw(canvas, pixel.IM.Scaled(cards.Orig, 2.))

	canvas.Draw(win, pixel.IM.Scaled(pixel.ZV, 1/camera.Cam.Zoom).Moved(pixel.V(camera.Cam.Pos.X, camera.Cam.Pos.Y)))
}