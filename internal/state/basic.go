package state

import "github.com/faiface/pixel/pixelgl"

type Exiting struct{}

func (s *Exiting) Initialize() {}

func (s *Exiting) TransitionIn() bool {
	return false
}

func (s *Exiting) TransitionOut() bool {
	return false
}

func (s *Exiting) Uninitialize() {}

func (s *Exiting) Update(win *pixelgl.Window) {
	win.SetClosed(true)
}

func (s *Exiting) Draw(win *pixelgl.Window) {}

func (s *Exiting) String() string {
	return "Exiting"
}

type NoState struct{}

func (s *NoState) Initialize() {}

func (s *NoState) TransitionIn() bool {
	return false
}

func (s *NoState) TransitionOut() bool {
	return false
}

func (s *NoState) Uninitialize() {}

func (s *NoState) Update(win *pixelgl.Window) {}

func (s *NoState) Draw(win *pixelgl.Window) {}

func (s *NoState) String() string {
	return "NoState"
}
