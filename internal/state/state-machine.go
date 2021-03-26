package state

import (
	"github.com/faiface/pixel/pixelgl"
)

var (
	TheMainMenu = &MainMenu{}
	TheStartRun = &StartRun{}
	TheEncounter = &Encounter{}
	TheUpgrade   = &Upgrade{}
)

func init() {
	Machine.State = &NoState{}
	Machine.NextState = TheMainMenu
	Machine.Trans = true
	Machine.Phase = Undefined
}

func Update(win *pixelgl.Window) {
	UpdateSwitchState()
	Machine.State.Update(win)
}

func Draw(win *pixelgl.Window) {
	Machine.State.Draw(win)
}

func SwitchState(newState State) {
	Machine.NextState = newState
	Machine.Trans = true
}

func UpdateSwitchState() {
	if Machine.NextState.String() != Machine.State.String() {
		if Machine.Trans {
			Machine.Trans = Machine.State.TransitionOut()
		} else {
			Machine.State.Uninitialize()
			Machine.NextState.Initialize()
			Machine.State = Machine.NextState
			Machine.Trans = true
		}
	} else if Machine.Trans {
		Machine.Trans = Machine.State.TransitionIn()
	}
}