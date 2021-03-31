package player

import (
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
)

func init() {
	CardManager = cardManager{
		moves:  []cardMove{},
	}
}

type CardGroup interface {
	AddCard(*Card)
	RemoveCard(uuid.UUID) *Card
}

func BuildGroup(cards []*Card, group CardGroup) {
	for _, c := range cards {
		addCard(group, c)
	}
}

func addCard(a CardGroup, card *Card) {
	card.Current = a
	a.AddCard(card)
}

func moveCard(a, b CardGroup, id uuid.UUID) {
	card := a.RemoveCard(id)
	if card != nil {
		card.Previous = card.Current
		card.Current = b
		b.AddCard(card)
	}
}

var CardManager cardManager

type cardManager struct {
	moves  []cardMove
}

type cardMove struct {
	from CardGroup
	to   CardGroup
	card *Card
}

func (m *cardManager) Update() {
	nextMoves := m.moves
	m.moves = []cardMove{}
	for _, move := range nextMoves {
		if util.IsNil(move.from) {
			addCard(move.to, move.card)
		} else {
			moveCard(move.from, move.to, move.card.ID)
		}
	}
}

func (m *cardManager) Move(from, to CardGroup, card *Card) {
	m.moves = append(m.moves, cardMove{
		from: from,
		to:   to,
		card: card,
	})
}