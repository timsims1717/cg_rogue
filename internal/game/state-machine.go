package game

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"golang.org/x/image/colornames"
)

type Phase int

func (p Phase) String() string {
	switch p {
	case PlayerStartTurn:
		return "Player Start Turn"
	case PlayerTurn:
		return "Player Turn"
	case EnemyStartTurn:
		return "Enemy Start Turn"
	case EnemyEndTurn:
		return "Enemy End Turn"
	case GameOver:
		return "Game Over"
	}
	return "Undefined"
}

const (
	PlayerStartTurn = iota
	PlayerTurn
	EnemyStartTurn
	EnemyEndTurn
	GameOver
)

var StateMachine stateMachine

type stateMachine struct {
	State Phase
}

func Initialize() {
	StateMachine.State = EnemyStartTurn
}

func Update() {
	if StateMachine.State != GameOver && player.Player1.Character.IsDestroyed() {
		player.Player1.EndTurn()
		ui.CenterText.Raw = "Game Over"
		ui.CenterText.Show = true
		ui.CenterText.TextColor = colornames.Black
		transform := animation.TransformBuilder{
			Target:  ui.CenterText,
			InterX:  nil,
			InterY:  nil,
			InterR:  nil,
			InterSX: gween.New(ui.CenterText.Scalar.X, 7.0, 2.0, ease.Linear),
			InterSY: gween.New(ui.CenterText.Scalar.Y, 7.0, 2.0, ease.Linear),
		}
		ui.CenterText.TransformEffect = transform.Build()
		ui.CenterText.ColorEffect = animation.FadeIn(ui.CenterText, 2.0)
		StateMachine.State = GameOver
	}
	switch StateMachine.State {
	case PlayerStartTurn:
		// todo: effects?
		player.Player1.StartTurn()
		camera.Cam.MoveTo(player.Player1.Character.Pos, 0.2, true)
		StateMachine.State = PlayerTurn
	case PlayerTurn:
		if player.Player1.ActionsThisTurn > 0 && player.Player1.PlayCard.Card == nil && !actions.IsActing() {
			player.Player1.EndTurn()
			StateMachine.State = EnemyStartTurn
		}
	case EnemyStartTurn:
		// todo: effects?
		ai.AIManager.StartAITurn()
		StateMachine.State = EnemyEndTurn
	case EnemyEndTurn:
		if !actions.IsActing() {
			ai.AIManager.EndAITurn()
			StateMachine.State = PlayerStartTurn
		}
	case GameOver:

	}
}