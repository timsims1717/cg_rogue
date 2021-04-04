package state

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cards"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/run"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

type StartRun struct {}

func (s *StartRun) Initialize() {
	run.StartRun()

	charsheet, err := img.LoadSpriteSheet("assets/character/testmananim.json")
	if err != nil {
		panic(err)
	}
	character := characters.NewCharacter(pixel.NewSprite(charsheet.Img, charsheet.Sprites[rand.Intn(len(charsheet.Sprites))]), world.Coords{X: 0, Y: 0}, characters.Ally, 10)
	player.Player1 = player.NewPlayer(character)

	player.InitializeCards()
	player.Player1.Hand = player.NewHand(player.Player1)
	player.BuildGroup([]*player.Card{
		cards.CreateThrust(),
		cards.CreateDash(),
		cards.CreateQuickStrike(),
		cards.CreateSweep(),
		cards.CreateVault(),
		cards.CreateDaggerThrow(),
		cards.CreateDisengage(),
	}, player.Player1.Hand)
	player.Player1.PlayCard = player.NewPlayCard(player.Player1)
	player.Player1.Discard = player.NewDiscard(player.Player1)
	player.Player1.Grid = player.NewGrid(player.Player1)

	//characters.CharacterManager.Add(character)
}

func (s *StartRun) TransitionIn() bool {
	SwitchState(TheEncounter)
	return false
}

func (s *StartRun) TransitionOut() bool {
	return false
}

func (s *StartRun) Uninitialize() {}

func (s *StartRun) Update(win *pixelgl.Window) {}

func (s *StartRun) Draw(win *pixelgl.Window) {}

func (s *StartRun) String() string {
	return "StartRun"
}