package actions

import "github.com/timsims1717/cg_rogue_go/internal/player"

type RestAction struct{
	source *player.Player
	isDone bool
}

func NewRestAction(source *player.Player) *RestAction {
	if source == nil {
		return nil
	}
	return &RestAction{
		source: source,
	}
}

func (a *RestAction) Update() {
	for i := len(a.source.Discard.Group)-1; i >= 0; i-- {
		card := a.source.Discard.Group[i]
		player.CardManager.Move(a.source.Discard, a.source.Hand, card)
	}
	a.isDone = true
}

func (a *RestAction) IsDone() bool {
	return a.isDone
}