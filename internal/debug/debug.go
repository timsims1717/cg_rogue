package debug

import "github.com/faiface/pixel/pixelgl"

func Initialize() {
	InitializeText()
	//InitializeLines()
}

func Update() {
	UpdateText()
	//UpdateLines()
}

func Draw(win *pixelgl.Window) {
	DrawText(win)
	//DrawLines(win)
}
