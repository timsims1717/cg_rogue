package game

import "github.com/timsims1717/cg_rogue_go/pkg/text"

var (
	centerText       text.TextDisplay
)

func SetCenterText(t string) {
	centerText.Raw = t
}

func UpdateCenterText() {

}