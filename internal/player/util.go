package player

import "bytes"

const (
	// Card constants
	BaseCardWidth  = 250
	BaseCardHeight = 350

	// Hand constants
	HandWidth        = 0.7
	HandCardScale    = 0.65
	HandHovCardScale = 0.8
	CardStart        = BaseCardWidth *0.5

	// PlayCard constants
	PlayRightPad  = BaseCardWidth * 0.65
	PlayCardScale = 1.0
	PlayBottomPad = BaseCardHeight * 1.45
)

func MakeKey(args... string) string {
	var b bytes.Buffer
	for i, arg := range args {
		if i != 0 {
			b.Write([]byte("-"))
		}
		b.Write([]byte(arg))
	}
	return b.String()
}