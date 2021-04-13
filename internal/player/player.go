package player

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
)

var Player1 *Player

type Player struct {
	Character       *floor.Character
	Input           *input.Input
	Hand            *Hand
	PlayCard        *PlayCard
	Discard         *Discard
	Grid            *Grid
	ActionsThisTurn int
	IsTurn          bool
}

func init() {
	Player1 = NewPlayer(nil)
}

func NewPlayer(character *floor.Character) *Player {
	return &Player{
		Character: character,
		Input:     input.NewInput(),
	}
}

func (p *Player) StartTurn() {
	p.ActionsThisTurn = 0
	p.IsTurn = true
}

func (p *Player) EndTurn() {
	p.IsTurn = false
}

func (p *Player) Update() {
	if p.Grid != nil {
		p.Grid.Update()
	}
	if p.Hand != nil {
		p.Hand.Update(p.IsTurn)
	}
	if p.PlayCard != nil {
		p.PlayCard.Update(p.IsTurn)
	}
	if p.Discard != nil {
		p.Discard.Update(p.IsTurn)
	}
}

func (p *Player) Draw(win *pixelgl.Window) {
}

func (p *Player) GetDeck() []*Card {
	var deck []*Card
	if p.Hand != nil {
		for _, c := range p.Hand.Group {
			deck = append(deck, c)
		}
	}
	if p.PlayCard != nil && p.PlayCard.Card != nil {
		deck = append(deck, p.PlayCard.Card)
	}
	if p.Discard != nil {
		for _, c := range p.Discard.Group {
			deck = append(deck, c)
		}
	}
	return deck
}
