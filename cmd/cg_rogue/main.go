package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/cards"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/debug"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/game"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
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

	debug.Initialize()
	game.Initialize()
	ui.Initialize()

	treesheet, err := img.LoadSpriteSheet("assets/character/trees.json")
	if err != nil {
		panic(err)
	}
	spritesheet, err := img.LoadSpriteSheet("assets/img/testfloor.json")
	if err != nil {
		panic(err)
	}
	charsheet, err := img.LoadSpriteSheet("assets/character/testmananim.json")
	if err != nil {
		panic(err)
	}
	uisheet, err := img.LoadSpriteSheet("assets/img/ui/selectors.json")
	if err != nil {
		panic(err)
	}
	selectors.SelectionSet.SetSpriteSheet(uisheet)

	camera.Cam = camera.New()

	floor.DefaultFloor(10, 10, spritesheet)

	tree := characters.NewCharacter(pixel.NewSprite(treesheet.Img, treesheet.Sprites[rand.Intn(len(treesheet.Sprites))]), world.Coords{8,4}, characters.Enemy, 10)
	treeAI := ai.NewAI(ai.RandomWalkerDecision, ai.RandomWalkerAct, tree)
	flyer := characters.NewCharacter(pixel.NewSprite(treesheet.Img, treesheet.Sprites[rand.Intn(len(treesheet.Sprites))]), world.Coords{2,9}, characters.Enemy, 10)
	flyerAI := ai.NewAI(ai.FlyChaserDecision, ai.FlyChaserAct, flyer)
	ai.AIManager.AddAI(treeAI)
	ai.AIManager.AddAI(flyerAI)


	character := characters.NewCharacter(pixel.NewSprite(charsheet.Img, charsheet.Sprites[rand.Intn(len(charsheet.Sprites))]), world.Coords{0,0}, characters.Ally, 10)
	player.Player1 = player.NewPlayer(character)
	player.Initialize()
	player.Player1.Hand = player.NewHand(player.Player1)
	player.Player1.Hand.AddCard(cards.CreateThrust())
	player.Player1.Hand.AddCard(cards.CreateDash())
	player.Player1.Hand.AddCard(cards.CreateQuickStrike())
	player.Player1.Hand.AddCard(cards.CreateVault())
	player.Player1.Hand.AddCard(cards.CreateSweep())
	player.Player1.Hand.AddCard(cards.CreateDaggerThrow())
	player.Player1.PlayCard = player.NewPlayCard(player.Player1)
	player.Player1.Discard = player.NewDiscard(player.Player1)

	characters.CharacterManager.Add(tree)
	characters.CharacterManager.Add(flyer)
	characters.CharacterManager.Add(character)

	camera.Cam.CenterOn([]pixel.Vec{character.Pos})
	timing.Reset()
	for !win.Closed() {
		timing.Update()

		debug.Update()
		game.Update()
		player.Player1.Input.Update(win)
		camera.Cam.Update(win)
		ui.Update()

		player.CardManager.Update()
		actions.Update()

		ai.Update()
		characters.Update()
		player.Player1.Update(win)

		win.Clear(colornames.Forestgreen)

		floor.CurrentFloor.Draw(win)
		characters.Draw(win)
		selectors.SelectionSet.Draw(win)
		win.SetSmooth(true)
		player.Player1.Hand.Draw(win)
		player.Player1.PlayCard.Draw(win)
		player.Player1.Discard.Draw(win)
		win.SetSmooth(false)
		ui.Draw(win)
		debug.Draw(win)
		win.Update()
	}
}

// fire the run function (the real main function)
func main() {
	pixelgl.Run(run)
}