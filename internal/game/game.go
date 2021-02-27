package game

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/player"
)

type state struct {
	Phase    Phase
	Start    bool
}

type Phase int

func (p Phase) String() string {
	switch p {
	case PlayerTurn:
		return "Player Turn"
	case EnemyTurn:
		return "Enemy Turn"
	}
	return "Undefined"
}

const (
	PlayerTurn = iota
	EnemyTurn
)

var StateMachine stateMachine

type stateMachine struct {
	State state
}

func Initialize() {
	StateMachine.State.Start = false
	StateMachine.State.Phase = EnemyTurn
}

func Update() {
	switch StateMachine.State.Phase {
	case PlayerTurn:
		if StateMachine.State.Start {
			// todo: effects?
			player.Player1.StartTurn()
			StateMachine.State.Start = false
		} else if player.Player1.CardsPlayed > 0 && !actions.IsActing() {
			player.Player1.EndTurn()
			StateMachine.State.Phase = EnemyTurn
			StateMachine.State.Start = true
		}
	case EnemyTurn:
		if StateMachine.State.Start {
			// todo: effects?
			StateMachine.State.Start = false
			ai.AIManager.StartAITurn()
		} else if !actions.IsActing() {
			ai.AIManager.EndAITurn()
			StateMachine.State.Phase = PlayerTurn
			StateMachine.State.Start = true
		}
	}
}