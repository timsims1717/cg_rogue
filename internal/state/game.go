package state

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/generate"
	"github.com/timsims1717/cg_rogue_go/internal/manager"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/run"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/sfx"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
)

type DefaultActions struct {
	player *player.Player
	Group  []*player.Card
}

func (g *DefaultActions) AddCard(card *player.Card) {
	if card != nil {
		g.Group = append(g.Group, card)
	}
}

func (g *DefaultActions) RemoveCard(id uuid.UUID) *player.Card {
	index := -1
	for i, card := range g.Group {
		if card.ID == id {
			index = i
			break
		}
	}
	if index < 0 || index >= len(g.Group) {
		return nil
	}
	card := g.Group[index]
	g.Group = append(g.Group[0:index], g.Group[index+1:len(g.Group)]...)
	return card
}

type Encounter struct {
	DefaultActions *DefaultActions
	RestButton     *ui.ActionEl
}

func (s *Encounter) Initialize() {
	uisheet, err := img.LoadSpriteSheet("assets/img/ui/selectors.json")
	if err != nil {
		panic(err)
	}
	selector.SelectionSet.SetSpriteSheet(uisheet)

	InitializeCenterText()

	generate.LoadTestFloor(run.CurrentRun.Level)

	restCard := CreateRest(player.Player1)
	moveCard := CreateBasicMove()
	atkCard := CreateCheatAttack()
	restCard.Rests = 0
	moveCard.Rests = 0
	atkCard.Rests = 0
	s.DefaultActions = &DefaultActions{
		player: player.Player1,
	}
	player.BuildGroup([]*player.Card{
		restCard,
		moveCard,
		atkCard,
	}, s.DefaultActions)

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
		player.CardManager.Move(s.DefaultActions, player.Player1.PlayCard, restCard)
	})
	restButton.SetOnDisabledFn(func() {
		restButton.Show = false
	})
	restButton.SetEnabledFn(func() {
		restButton.Show = true
	})
	s.RestButton = restButton
	player.Player1.Input.SetHotKey(pixelgl.KeyR, func() {
		restButton.Click()
	})

	player.Player1.Input.SetHotKey(pixelgl.KeyA, func() {
		player.CardManager.Move(s.DefaultActions, player.Player1.PlayCard, atkCard)
	})

	Machine.Phase = EnemyStartTurn
	camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.White, 1.0)
}

func (s *Encounter) TransitionIn() bool {
	if camera.Cam.Effect != nil {
		return true
	}
	sfx.MusicPlayer.SetCurrentTracks([]string{"test_track", "test_ambience"})
	sfx.MusicPlayer.PlayNextTrack(0.0, 1., 0.5)
	return false
}

func (s *Encounter) TransitionOut() bool {
	return camera.Cam.Effect != nil
}

func (s *Encounter) Uninitialize() {
	floor.CurrentFloor = nil
	InitializeCenterText()
	ai.AIManager.Clear()
	floor.CharacterManager.Clear()
}

func (s *Encounter) Update(win *pixelgl.Window) {
	UpdateEncounterPhase()
	player.Player1.Input.Update(win)
	camera.Cam.Update(win)
	CenterText.Update(player.Player1.Input)

	player.CardManager.Update()
	manager.ActionManager.Update()

	ai.AIManager.Update()
	floor.Update()
	player.Player1.Update()
	if s.RestButton != nil {
		s.RestButton.Disabled = !player.Player1.IsTurn
		s.RestButton.Update(player.Player1.Input)
	}
}

func (s *Encounter) Draw(win *pixelgl.Window) {
	floor.CurrentFloor.Draw(win)
	floor.Draw(win)
	selector.SelectionSet.Draw(win)
	player.Player1.Draw(win)
	win.SetSmooth(true)
	player.Player1.Hand.Draw(win)
	player.Player1.PlayCard.Draw(win)
	player.Player1.Discard.Draw(win)
	if player.Player1.Grid != nil && player.Player1.Grid.Show {
		player.Player1.Grid.Draw(win)
	}
	win.SetSmooth(false)
	if s.RestButton != nil {
		s.RestButton.Draw(win)
	}
	CenterText.Draw(win)
}

func UpdateEncounterPhase() {
	if Machine.Phase == Undefined {
		return
	}
	if Machine.Phase != EncounterComplete && Machine.Phase != GameOver && player.Player1.Character.IsDestroyed() {
		player.Player1.EndTurn()
		CenterText.Text.Raw = "Game Over"
		CenterText.Show = true
		CenterText.Text.TextColor = colornames.Black
		transform := animation.TransformBuilder{
			Transform: CenterText.Text.Transform,
			InterSX:   gween.New(CenterText.Text.Transform.Scalar.X, 7.0, 2.0, ease.Linear),
			InterSY:   gween.New(CenterText.Text.Transform.Scalar.Y, 7.0, 2.0, ease.Linear),
		}
		CenterText.Text.TransformEffect = transform.Build()
		CenterText.Text.ColorEffect = animation.FadeIn(CenterText.Text, 2.0)
		Machine.Phase = GameOver
		camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.Black,4.)
		sfx.MusicPlayer.FadeOut(4.)
		SwitchState(TheMainMenu)
	}
	if Machine.Phase != EncounterComplete && Machine.Phase != GameOver {
		allDead := true
		for _, c := range floor.CharacterManager.GetDiplomatic(floor.Enemy, player.Player1.Character.GetCoords(), 50) {
			if cha := floor.CurrentFloor.GetOccupant(c); cha != nil {
				if !cha.IsDestroyed() {
					allDead = false
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
				InterSX:   gween.New(CenterText.Text.Transform.Scalar.X, 7.0, 2.0, ease.Linear),
				InterSY:   gween.New(CenterText.Text.Transform.Scalar.Y, 7.0, 2.0, ease.Linear),
			}
			CenterText.Text.TransformEffect = transform.Build()
			CenterText.Text.ColorEffect = animation.FadeIn(CenterText, 2.0)
			Machine.Phase = EncounterComplete
			manager.ActionManager.AddToBot(actions.NewHealAction([]world.Coords{player.Player1.Character.Coords}, selector.ActionValues{
				Source:   player.Player1.Character,
				Heal:     1,
			}))
			manager.ActionManager.AddToBot(actions.NewRestAction(player.Player1))
			camera.Cam.Effect = animation.FadeTo(camera.Cam, colornames.Black,4.)
			SwitchState(TheUpgrade)
		}
	}
	switch Machine.Phase {
	case PlayerStartTurn:
		// todo: effects?
		player.Player1.StartTurn()
		camera.Cam.MoveTo(player.Player1.Character.Transform.Pos, 0.2, true)
		Machine.Phase = PlayerTurn
	case PlayerTurn:
		if player.Player1.ActionsThisTurn > 0 && player.Player1.PlayCard.Card == nil && !manager.ActionManager.IsActing() {
			player.Player1.EndTurn()
			Machine.Phase = EnemyStartTurn
		}
	case EnemyStartTurn:
		// todo: effects?
		ai.AIManager.StartAITurn()
		Machine.Phase = EnemyEndTurn
	case EnemyEndTurn:
		if !ai.AIManager.AIActing() && !manager.ActionManager.IsActing() {
			ai.AIManager.EndAITurn()
			Machine.Phase = PlayerStartTurn
		}
	}
}

func (s *Encounter) String() string {
	return "Encounter"
}

type Rest struct {
	*player.Card
	player *player.Player
}

func (r *Rest) DoActions() {
	manager.ActionManager.AddToBot(actions.NewRestAction(r.player))
}

func (r *Rest) SetValues(_ int) {}

func (r *Rest) InitSelectors() {}

func (r *Rest) SetCard(card *player.Card) {
	r.Card = card
}

func CreateRest(thePlayer *player.Player) *player.Card {
	card := player.NewCard("", "", &Rest{
		player: thePlayer,
	})
	card.SetDraw(false)
	return card
}

type BasicMove struct {
	*player.Card
}

func (m *BasicMove) DoActions() {
	manager.ActionManager.AddToBot(actions.NewMoveSeriesAction(m.Values.Source, m.Values.Source, m.Results[0][0].Area))
}

func (m *BasicMove) SetValues(_ int) {
	values := selector.ActionValues{
		Move: 1,
	}
	m.Values = values
}

func (m *BasicMove) InitSelectors() {
	m.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.PathSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    true,
				NonEmpty:      true,
				EndUnoccupied: true,
				Orig:          world.Coords{},
			},
			Effect: selector.NewSelectionEffect(&selector.MoveEffect{}),
		}, true),
	}
}

func (m *BasicMove) SetCard(card *player.Card) {
	m.Card = card
}

func CreateBasicMove() *player.Card {
	card := player.NewCard("", "", &BasicMove{})
	card.SetDraw(false)
	return card
}

func (m *BasicMove) DoAction(path []world.Coords) {
	manager.ActionManager.AddToBot(actions.NewMoveSeriesAction(m.Values.Source, m.Values.Source, path))
}

type CheatAttack struct {
	*player.Card
}

func (a *CheatAttack) DoActions() {
	manager.ActionManager.AddToBot(actions.NewDamageHexAction(a.Results[0][0].Area, a.Values))
}

func (a *CheatAttack) SetValues(_ int) {
	values := selector.ActionValues{
		Damage:  10,
		Range:   10,
		Targets: 5,
	}
	a.Values = values
}

func (a *CheatAttack) InitSelectors() {
	a.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.HexSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    false,
				NonEmpty:      false,
				EndUnoccupied: false,
				Orig:          world.Coords{},
			},
			Effect: selector.NewSelectionEffect(&selector.AttackEffect{}),
		}, false),
	}
}

func (a *CheatAttack) SetCard(card *player.Card) {
	a.Card = card
}

func CreateCheatAttack() *player.Card {
	card := player.NewCard("", "", &CheatAttack{})
	card.SetDraw(false)
	return card
}