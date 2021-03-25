package player

const (
	// Card constants
	BaseCardWidth  = 250
	BaseCardHeight = 350

	// Hand constants
	HandCardScale    = 0.65
	HandHovCardScale = 0.8
	HandLeftPad      = BaseCardWidth *0.5
	HandBottomPad    = BaseCardHeight * 0.25 * HandCardScale

	// PlayCard constants
	PlayRightPad  = BaseCardWidth * 0.65
	PlayCardScale = 1.0
	PlayBottomPad = BaseCardHeight * 1.45

	// Discard constants
	DiscardRightPad  = BaseCardWidth * 1.4
	DiscardBottomPad = BaseCardHeight * 0.25 * DiscardScale
	DiscardScale     = 0.5
	DiscardHovScale  = 0.7

	// Grid constants
	GridCardScale    = 0.75
	GridHovCardScale = 0.9

	// Action Button constants
	ButtonRightPad = 220.
	RestBottomPad = 60.
	MoveBottomPad = 100.
)