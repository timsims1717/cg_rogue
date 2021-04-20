package debug

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/player"
)

func Initialize() {
	InitializeText()
	//InitializeLines()
}

func Update() {
	if player.Player1.Input.Debug {
		fmt.Println("DEBUG PAUSED")
	}
	UpdateText()
	//UpdateLines()
}

func Draw(win *pixelgl.Window) {
	DrawText(win)
	//DrawLines(win)
}
