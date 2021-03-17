package game

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/state"
)

func init() {
	state.Machine.State = state.NoState
	state.Machine.NextState = state.MainMenu
	state.Machine.Trans = true
	state.Machine.Phase = state.Undefined
}

func Update(win *pixelgl.Window) {
	UpdateSwitchState(win)
	switch state.Machine.State {
	case state.MainMenu:
		UpdateMenu(win)
	case state.InGame:
		UpdateGame(win)
	case state.Exiting:
		win.SetClosed(true)
	}
}

func Draw(win *pixelgl.Window) {
	switch state.Machine.State {
	case state.MainMenu:
		DrawMenu(win)
	case state.InGame:
		DrawGame(win)
	}
}

func SwitchState(newState state.State) {
	state.Machine.NextState = newState
	state.Machine.Trans = true
}

func UpdateSwitchState(win *pixelgl.Window) {
	if state.Machine.NextState != state.Machine.State {
		if state.Machine.Trans {
			switch state.Machine.State {
			case state.MainMenu:
				state.Machine.Trans = TransitionOutMenu()
			case state.InGame:
				state.Machine.Trans = TransitionOutGame()
			default:
				state.Machine.Trans = false
			}
		} else {
			switch state.Machine.State {
			case state.MainMenu:
				UninitializeMenu()
			case state.InGame:
				UninitializeGame()
			}
			switch state.Machine.NextState {
			case state.MainMenu:
				InitializeMenu(win)
			case state.InGame:
				InitializeGame()
			}
			state.Machine.State = state.Machine.NextState
			state.Machine.Trans = true
		}
	} else if state.Machine.Trans {
		switch state.Machine.State {
		case state.MainMenu:
			state.Machine.Trans = TransitionInMenu()
		case state.InGame:
			state.Machine.Trans = TransitionInGame()
		default:
			state.Machine.Trans = false
		}
	}
}