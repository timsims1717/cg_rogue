package player

type RestAction struct{
	source *Player
	isDone bool
}

func NewRestAction(source *Player) *RestAction {
	if source == nil {
		return nil
	}
	return &RestAction{
		source: source,
	}
}

func (a *RestAction) Update() {
	for i := len(a.source.Discard.Group)-1; i >= 0; i-- {
		//card := a.source.Discard.Group[i]
		CardManager.Move(a.source.Discard, a.source.Hand, i)
	}
	a.isDone = true
}

func (a *RestAction) IsDone() bool {
	return a.isDone
}