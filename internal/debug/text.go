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
	"github.com/timsims1717/cg_rogue_go/internal/state"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/typeface"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"time"
)

var (
	// The canvasText where the debug information is drawn
	canvasText *pixelgl.Canvas
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
	misc    *text.Text
)

// InitializeText creates the debug canvasText and all the text containers.
// This is where the location of the text containers is set.
func InitializeText() {
	canvasText = pixelgl.NewCanvas(pixel.R(0, 0, cfg.WindowWidthF, cfg.WindowHeightF))
	fps = text.New(pixel.ZV, typeface.BasicAtlas)
	mouse = text.New(pixel.ZV, typeface.BasicAtlas)
	worlds = text.New(pixel.ZV, typeface.BasicAtlas)
	maps = text.New(pixel.ZV, typeface.BasicAtlas)
	phase = text.New(pixel.ZV, typeface.BasicAtlas)
	health = text.New(pixel.ZV, typeface.BasicAtlas)
	hand = text.New(pixel.ZV, typeface.BasicAtlas)
	handL = 0
	playing = text.New(pixel.ZV, typeface.BasicAtlas)
	discard = text.New(pixel.ZV, typeface.BasicAtlas)
	discL = 0
	misc = text.New(pixel.ZV, typeface.BasicAtlas)
}

// UpdateText clears the text containers and updates them with the correct information.
func UpdateText() {
	frames++
	select {
	case <-second:
		fps.Clear()
		fmt.Fprintf(fps, "FPS: %d", frames)
		frames = 0
	default:
	}
	mousePtr := player.Player1.Input.Cursor
	wrldPtr := player.Player1.Input.World
	mapX, mapY := world.WorldToMap(wrldPtr.X, wrldPtr.Y)
	mouse.Clear()
	worlds.Clear()
	maps.Clear()
	fmt.Fprintf(mouse, "Mouse (X,Y): (%d,%d)", int(mousePtr.X), int(mousePtr.Y))
	fmt.Fprintf(worlds, "World (X,Y): (%d,%d)", int(wrldPtr.X), int(wrldPtr.Y))
	fmt.Fprintf(maps, "Map (X,Y): (%d,%d)", mapX, mapY)
	phase.Clear()
	fmt.Fprintf(phase, "Phase: %s", state.Machine.Phase.String())
	dispHP = false
	health.Clear()
	if state.Machine.State == state.InGame {
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
		for i := len(player.Player1.Discard.Group) - 1; i >= 0; i-- {
			card := player.Player1.Discard.Group[i]
			fmt.Fprintf(discard, "   %s (%d)\n", card.RawTitle, i)
		}
		discL = len(player.Player1.Discard.Group) + 1
		//misc.Clear()
		//fmt.Fprintf(misc, "center text: (%s, %s)", ui2.CenterText.Pos.X, ui2.CenterText.Pos.Y)
	}
}

// DrawText draws each text container to the canvasText.
// This is where scaling happens.
func DrawText(win *pixelgl.Window) {
	canvasText.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
	height := cfg.WindowHeightF - typeface.BasicAtlas.LineHeight() - 20.
	fps.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
	height -= typeface.BasicAtlas.LineHeight() + 20.
	mouse.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
	height -= typeface.BasicAtlas.LineHeight() + 20.
	worlds.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
	height -= typeface.BasicAtlas.LineHeight() + 20.
	maps.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
	height -= typeface.BasicAtlas.LineHeight() + 20.
	phase.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
	height -= typeface.BasicAtlas.LineHeight() + 20.
	if state.Machine.State == state.InGame {
		if dispHP {
			health.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
		}
		height -= typeface.BasicAtlas.LineHeight() + 20.
		hand.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
		height -= typeface.BasicAtlas.LineHeight()*float64(handL*2) + 20.
		playing.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
		height -= typeface.BasicAtlas.LineHeight() + 20.
		discard.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
		height -= typeface.BasicAtlas.LineHeight()*float64(discL*2) + 20.
		//misc.Draw(canvasText, pixel.IM.Scaled(pixel.ZV, 2.).Moved(pixel.V(20., height)))
	}

	canvasText.Draw(win, pixel.IM.Scaled(pixel.ZV, 1/camera.Cam.Zoom).Moved(pixel.V(camera.Cam.Pos.X, camera.Cam.Pos.Y)))
}