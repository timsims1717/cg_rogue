package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/cards"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/generate"
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
)

func InitializeGame() {
	uisheet, err := img.LoadSpriteSheet("assets/img/ui/selectors.json")
	if err != nil {
		panic(err)
	}
	selectors.SelectionSet.SetSpriteSheet(uisheet)

	InitializeCenterText()
	player.Initialize()

	generate.LoadTestFloor(1)

	//floor.DefaultFloor(10, 10, spritesheet)

	//tree := characters.NewCharacter(pixel.NewSprite(treesheet.Img, treesheet.Sprites[rand.Intn(len(treesheet.Sprites))]), world.Coords{8,4}, characters.Enemy, 10)
	//treeAI := ai.NewRandomWalker(tree)
	//flyer := characters.NewCharacter(pixel.NewSprite(treesheet.Img, treesheet.Sprites[rand.Intn(len(treesheet.Sprites))]), world.Coords{2,9}, characters.Enemy, 10)
	//flyerAI := ai.NewFlyChaser(flyer)
	//ai.AIManager.AddAI(treeAI)
	//ai.AIManager.AddAI(flyerAI)

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
	rest.Transform.Scalar = pixel.V(2.5, 2.5)
	rest.TextColor = colornames.Purple
	restButton := ui.NewActionEl(rest, pixel.R(0., 0., rest.Text.BoundsOf(restS).W() * 2.5, rest.Text.BoundsOf(restS).H() * 2.5), true)
	restButton.Show = true
	restButton.Transform.Pos = pixel.V(camera.WindowWidthF - player.ButtonRightPad, player.RestBottomPad)
	restButton.SetOnHoverFn(func() {
		restButton.Text.TextColor = colornames.Forestgreen
	})
	restButton.SetUnHoverFn(func() {
		restButton.Text.TextColor = colornames.Purple
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
	move.Transform.Scalar = pixel.V(2.5, 2.5)
	move.TextColor = colornames.Purple
	moveButton := ui.NewActionEl(move, pixel.R(0., 0., move.Text.BoundsOf(moveS).W() * 2.5, move.Text.BoundsOf(moveS).H() * 2.5), true)
	moveButton.Show = true
	moveButton.Transform.Pos = pixel.V(camera.WindowWidthF - player.ButtonRightPad, player.MoveBottomPad)
	moveButton.SetOnHoverFn(func() {
		moveButton.Text.TextColor = colornames.Forestgreen
	})
	moveButton.SetUnHoverFn(func() {
		moveButton.Text.TextColor = colornames.Purple
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
	state.Machine.Phase = state.EnemyStartTurn
	camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.White, 1.0)
}

func TransitionInGame() bool {
	return camera.Cam.Effect != nil
}

func TransitionOutGame() bool {
	return camera.Cam.Effect != nil
}

func UninitializeGame() {
	floor.CurrentFloor = nil
	InitializeCenterText()
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

	ai.AIManager.Update()
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
	if state.Machine.Phase == state.Undefined {
		return
	}
	if state.Machine.Phase != state.EncounterComplete && state.Machine.Phase != state.GameOver && player.Player1.Character.IsDestroyed() {
		player.Player1.EndTurn()
		CenterText.Text.Raw = "Game Over"
		CenterText.Show = true
		CenterText.Text.TextColor = colornames.Black
		transform := animation.TransformBuilder{
			Transform: CenterText.Text.Transform,
			InterX:    nil,
			InterY:    nil,
			InterR:    nil,
			InterSX:   gween.New(CenterText.Text.Transform.Scalar.X, 7.0, 2.0, ease.Linear),
			InterSY:   gween.New(CenterText.Text.Transform.Scalar.Y, 7.0, 2.0, ease.Linear),
		}
		CenterText.Text.TransformEffect = transform.Build()
		CenterText.Text.ColorEffect = animation.FadeIn(CenterText.Text, 2.0)
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
			CenterText.Text.Raw = "Success!"
			CenterText.Show = true
			CenterText.Text.TextColor = colornames.Black
			transform := animation.TransformBuilder{
				Transform: CenterText.Text.Transform,
				InterX:    nil,
				InterY:    nil,
				InterR:    nil,
				InterSX:   gween.New(CenterText.Text.Transform.Scalar.X, 7.0, 2.0, ease.Linear),
				InterSY:   gween.New(CenterText.Text.Transform.Scalar.Y, 7.0, 2.0, ease.Linear),
			}
			CenterText.Text.TransformEffect = transform.Build()
			CenterText.Text.ColorEffect = animation.FadeIn(CenterText, 2.0)
			state.Machine.Phase = state.EncounterComplete
		}
	}
	switch state.Machine.Phase {
	case state.PlayerStartTurn:
		// todo: effects?
		player.Player1.StartTurn()
		camera.Cam.MoveTo(player.Player1.Character.Transform.Pos, 0.2, true)
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
		camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.Black,4.)
		SwitchState(state.MainMenu)
		state.Machine.Phase = state.Undefined
	}
}