package game

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
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
	if player.Player1.Character.IsDestroyed() {
		player.Player1.EndTurn()
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