package state

import "github.com/faiface/pixel/pixelgl"

type State interface{
	//Load()
	//IsLoaded() bool
	Initialize()
	TransitionIn() bool
	TransitionOut() bool
	Uninitialize()
	Update(*pixelgl.Window)
	Draw(*pixelgl.Window)
	String() string
}

type phase int

const (
	Undefined = iota
	PlayerStartTurn
	PlayerTurn
	EnemyStartTurn
	EnemyEndTurn
	GameOver
	EncounterComplete
)

func (p phase) String() string {
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
	case EncounterComplete:
		return "Encounter Complete"
	}
	return "Undefined"
}

var Machine stateMachine

type stateMachine struct {
	Trans     bool
	NextState State
	State     State
	Phase     phase
}