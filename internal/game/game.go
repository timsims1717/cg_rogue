package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/cards"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/internal/state"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

func InitializeGame() {
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

	floor.DefaultFloor(10, 10, spritesheet)
	Initialize()

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

	restS := "Rest (R)"
	rest := ui.NewActionText(restS)
	rest.VAlign = ui.Center
	rest.Scalar = pixel.V(2.5, 2.5)
	rest.TextColor = colornames.Purple
	restButton := ui.NewActionEl(rest, pixel.R(0., 0., rest.Text.BoundsOf(restS).W() * 2.5, rest.Text.BoundsOf(restS).H() * 2.5))
	restButton.Show = true
	restButton.UI = true
	restButton.Pos = pixel.V(cfg.WindowWidthF - player.ButtonRightPad, player.RestBottomPad)
	restButton.SetOnHoverFn(func() {
		restButton.T.TextColor = colornames.Forestgreen
	})
	restButton.SetUnHoverFn(func() {
		restButton.T.TextColor = colornames.Purple
	})
	restButton.SetClickFn(func() {
		values := selectors.ActionValues{}
		sel := selectors.NewNullSelect()
		player.Player1.PlayCard.CancelCard()
		player.Player1.SetPlayerAction(player.NewPlayerAction(sel, values, player.Player1.Rest))
	})
	restButton.SetOnDisabledFn(func() {
		restButton.Show = false
	})
	restButton.SetEnabledFn(func() {
		restButton.Show = true
	})
	player.Player1.RestButton = restButton

	moveS := "Move 1 (M)"
	move := ui.NewActionText(moveS)
	move.VAlign = ui.Center
	move.Scalar = pixel.V(2.5, 2.5)
	move.TextColor = colornames.Purple
	moveButton := ui.NewActionEl(move, pixel.R(0., 0., move.Text.BoundsOf(moveS).W() * 2.5, move.Text.BoundsOf(moveS).H() * 2.5))
	moveButton.Show = true
	moveButton.UI = true
	moveButton.Pos = pixel.V(cfg.WindowWidthF - player.ButtonRightPad, player.MoveBottomPad)
	moveButton.SetOnHoverFn(func() {
		moveButton.T.TextColor = colornames.Forestgreen
	})
	moveButton.SetUnHoverFn(func() {
		moveButton.T.TextColor = colornames.Purple
	})
	moveButton.SetClickFn(func() {
		values := selectors.ActionValues{
		Source:  player.Player1.Character,
		Damage:  0,
		Move:    1,
		Range:   0,
		Targets: 0,
		Checks: floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    true,
			NonEmpty:      true,
			EndUnoccupied: true,
			Orig:          world.Coords{},
		},
	}
		sel := selectors.NewPathSelect()
		player.Player1.PlayCard.CancelCard()
		player.Player1.SetPlayerAction(player.NewPlayerAction(sel, values, player.BasicMove))
	})
	moveButton.SetOnDisabledFn(func() {
		moveButton.Show = false
	})
	moveButton.SetEnabledFn(func() {
		moveButton.Show = true
	})
	player.Player1.MoveButton = moveButton

	characters.CharacterManager.Add(tree)
	characters.CharacterManager.Add(flyer)
	characters.CharacterManager.Add(character)

	camera.Cam.CenterOn([]pixel.Vec{character.Pos})
	state.Machine.Phase = state.EnemyStartTurn
}

func TransitionInGame() bool {
	return false
}

func TransitionOutGame() bool {
	return false
}

func UninitializeGame() {
	floor.CurrentFloor = nil
	Initialize()
	ai.AIManager.Clear()
	characters.CharacterManager.Clear()
	player.Player1.Hand = nil
	player.Player1.PlayCard = nil
	player.Player1.Discard = nil
}

func UpdateGame(win *pixelgl.Window) {
	UpdateGamePhase()
	player.Player1.Input.Update(win)
	camera.Cam.Update(win)
	CenterText.Update(player.Player1.Input)

	player.CardManager.Update()
	actions.Update()

	ai.Update()
	characters.Update()
	player.Player1.Update(win)
}

func DrawGame(win *pixelgl.Window) {
	floor.CurrentFloor.Draw(win)
	characters.Draw(win)
	selectors.SelectionSet.Draw(win)
	player.Player1.Draw(win)
	win.SetSmooth(true)
	player.Player1.Hand.Draw(win)
	player.Player1.PlayCard.Draw(win)
	player.Player1.Discard.Draw(win)
	win.SetSmooth(false)
	CenterText.Draw(win)
}

func UpdateGamePhase() {
	if state.Machine.Phase != state.EncounterComplete && state.Machine.Phase != state.GameOver && player.Player1.Character.IsDestroyed() {
		player.Player1.EndTurn()
		CenterText.T.Raw = "Game Over"
		CenterText.Show = true
		CenterText.T.TextColor = colornames.Black
		transform := animation.TransformBuilder{
			Target:  CenterText.T,
			InterX:  nil,
			InterY:  nil,
			InterR:  nil,
			InterSX: gween.New(CenterText.T.Scalar.X, 7.0, 2.0, ease.Linear),
			InterSY: gween.New(CenterText.T.Scalar.Y, 7.0, 2.0, ease.Linear),
		}
		CenterText.T.TransformEffect = transform.Build()
		CenterText.T.ColorEffect = animation.FadeIn(CenterText.T, 2.0)
		state.Machine.Phase = state.GameOver
	}
	if state.Machine.Phase != state.EncounterComplete && state.Machine.Phase != state.GameOver {
		allDead := true
		for _, c := range characters.CharacterManager.GetDiplomatic(characters.Enemy, player.Player1.Character.GetCoords(), 50) {
			if occ := floor.CurrentFloor.GetOccupant(c); occ != nil {
				if cha, ok := occ.(*characters.Character); ok {
					if !cha.IsDestroyed() {
						allDead = false
					}
				}
			}
		}
		if allDead {
			player.Player1.EndTurn()
			CenterText.T.Raw = "Success!"
			CenterText.Show = true
			CenterText.T.TextColor = colornames.Black
			transform := animation.TransformBuilder{
				Target:  CenterText.T,
				InterX:  nil,
				InterY:  nil,
				InterR:  nil,
				InterSX: gween.New(CenterText.T.Scalar.X, 7.0, 2.0, ease.Linear),
				InterSY: gween.New(CenterText.T.Scalar.Y, 7.0, 2.0, ease.Linear),
			}
			CenterText.T.TransformEffect = transform.Build()
			CenterText.T.ColorEffect = animation.FadeIn(CenterText, 2.0)
			state.Machine.Phase = state.EncounterComplete
		}
	}
	switch state.Machine.Phase {
	case state.PlayerStartTurn:
		// todo: effects?
		player.Player1.StartTurn()
		camera.Cam.MoveTo(player.Player1.Character.Pos, 0.2, true)
		state.Machine.Phase = state.PlayerTurn
	case state.PlayerTurn:
		if player.Player1.ActionsThisTurn > 0 && player.Player1.PlayCard.Card == nil && !actions.IsActing() {
			player.Player1.EndTurn()
			state.Machine.Phase = state.EnemyStartTurn
		}
	case state.EnemyStartTurn:
		// todo: effects?
		ai.AIManager.StartAITurn()
		state.Machine.Phase = state.EnemyEndTurn
	case state.EnemyEndTurn:
		if !ai.AIManager.AIActing() && !actions.IsActing() {
			ai.AIManager.EndAITurn()
			state.Machine.Phase = state.PlayerStartTurn
		}
	case state.EncounterComplete:
		fallthrough
	case state.GameOver:
		go func() {
			time.Sleep(3 * time.Second)
			SwitchState(state.MainMenu)
		}()
	}
}