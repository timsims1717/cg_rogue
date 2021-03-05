package debug

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/game"
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
	// The text containers for the debug elements
	fps     *text.Text
	mouse   *text.Text
	worlds  *text.Text
	maps    *text.Text
	phase   *text.Text
	health  *text.Text
	dispHP  bool
	hand    *text.Text
	handL   int
	playing *text.Text
	discard *text.Text
	discL   int
)

// Initialize creates the debug canvas and all the text containers.
// This is where the location of the text containers is set.
func Initialize() {
	canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(cfg.WindowWidth), float64(cfg.WindowHeight)))
	fps = text.New(pixel.ZV, text2.BasicAtlas)
	mouse = text.New(pixel.ZV, text2.BasicAtlas)
	worlds = text.New(pixel.ZV, text2.BasicAtlas)
	maps = text.New(pixel.ZV, text2.BasicAtlas)
	phase = text.New(pixel.ZV, text2.BasicAtlas)
	health = text.New(pixel.ZV, text2.BasicAtlas)
	hand = text.New(pixel.ZV, text2.BasicAtlas)
	handL = 0
	playing = text.New(pixel.ZV, text2.BasicAtlas)
	discard = text.New(pixel.ZV, text2.BasicAtlas)
	discL = 0
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
	phase.Clear()
	fmt.Fprintf(phase, "Phase: %s", game.StateMachine.State.Phase.String())
	dispHP = false
	health.Clear()
	occ := floor.CurrentFloor.GetOccupant(world.Coords{mapX, mapY})
	if objects.NotNil(occ) {
		if cha, ok := occ.(*characters.Character); ok {
			fmt.Fprintf(health, "Health: %d/%d HP", cha.CurrHP, cha.MaxHP)
			dispHP = true
		}
	}
	hand.Clear()
	fmt.Fprintf(hand, "Player 1 Hand: (Hovered: %d)\n", player.Player1.Hand.Hovered)
	for i, card := range player.Player1.Hand.Group {
		hovered := card.PointInside(wrldPtr)
		fmt.Fprintf(hand, "   %s (%d): Hovered: %t\n", card.RawTitle, i, hovered)
	}
	handL = len(player.Player1.Hand.Group) + 1
	playing.Clear()
	playcard := "none"
	if player.Player1.PlayCard.Card != nil {
		playcard = player.Player1.PlayCard.Card.RawTitle
	}
	fmt.Fprintf(playing, "Player 1 Playing: %s", playcard)
	discard.Clear()
	fmt.Fprintf(discard, "Player 1 Discard: (Hovered: %t)\n", player.Player1.Discard.Hover)
	for i := len(player.Player1.Discard.Group)-1; i >= 0; i-- {
		card := player.Player1.Discard.Group[i]
		fmt.Fprintf(discard, "   %s (%d)\n", card.RawTitle, i)
	}
	discL = len(player.Player1.Discard.Group)
}

// Draw draws each text container to the canvas.
// This is where scaling happens.
func Draw(win *pixelgl.Window) {
	canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
	height := float64(cfg.WindowHeight) - text2.BasicAtlas.LineHeight() - 20.
	fps.Draw(canvas, pixel.IM.Scaled(fps.Orig, 2.).Moved(pixel.V(20., height)))
	height -= text2.BasicAtlas.LineHeight() + 20.
	mouse.Draw(canvas, pixel.IM.Scaled(mouse.Orig, 2.).Moved(pixel.V(20., height)))
	height -= text2.BasicAtlas.LineHeight() + 20.
	worlds.Draw(canvas, pixel.IM.Scaled(worlds.Orig, 2.).Moved(pixel.V(20., height)))
	height -= text2.BasicAtlas.LineHeight() + 20.
	maps.Draw(canvas, pixel.IM.Scaled(maps.Orig, 2.).Moved(pixel.V(20., height)))
	height -= text2.BasicAtlas.LineHeight() + 20.
	phase.Draw(canvas, pixel.IM.Scaled(phase.Orig, 2.).Moved(pixel.V(20., height)))
	height -= text2.BasicAtlas.LineHeight() + 20.
	if dispHP {
		health.Draw(canvas, pixel.IM.Scaled(health.Orig, 2.).Moved(pixel.V(20., height)))
	}
	height -= text2.BasicAtlas.LineHeight() + 20.
	hand.Draw(canvas, pixel.IM.Scaled(hand.Orig, 2.).Moved(pixel.V(20., height)))
	height -= text2.BasicAtlas.LineHeight() * float64(handL * 2) + 20.
	playing.Draw(canvas, pixel.IM.Scaled(playing.Orig, 2.).Moved(pixel.V(20., height)))
	height -= text2.BasicAtlas.LineHeight() + 20.
	discard.Draw(canvas, pixel.IM.Scaled(hand.Orig, 2.).Moved(pixel.V(20., height)))

	canvas.Draw(win, pixel.IM.Scaled(pixel.ZV, 1/camera.Cam.Zoom).Moved(pixel.V(camera.Cam.Pos.X, camera.Cam.Pos.Y)))
}