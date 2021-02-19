package player

func init() {
	CardManager = cardManager{
		moves:  []cardMove{},
	}
}

type CardGroup interface {
	AddCard(*Card)
	RemoveCard(int) *Card
}

func moveCard(a, b CardGroup, i int) {
	b.AddCard(a.RemoveCard(i))
}

var CardManager cardManager

type cardManager struct {
	moves  []cardMove
}

type cardMove struct {
	from CardGroup
	to   CardGroup
	i    int
}

func (m *cardManager) Update() {
	for _, move := range m.moves {
		moveCard(move.from, move.to, move.i)
	}
	m.moves = []cardMove{}
}

func (m *cardManager) Move(from, to CardGroup, i int) {
	m.moves = append(m.moves, cardMove{
		from: from,
		to:   to,
		i:    i,
	})
}