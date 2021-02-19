package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
)

type PlayCard struct {
	player *Player
	Card   *Card
	update bool
}

func NewPlayCard(player *Player) *PlayCard {
	return &PlayCard{
		player: player,
	}
}

func (p *PlayCard) Update() {
	if p.Card != nil {
		if p.update {
			p.Card.setXY(pixel.V(float64(cfg.WindowWidth) -PlayRightPad, PlayBottomPad))
			p.Card.setScalar(PlayCardScale)
			p.update = false
		}
		p.Card.Update()
		if p.Card.PointInside(p.player.Input.World) && p.player.Input.Cancel.JustPressed() && p.Card.canCancel {
			p.player.Input.Cancel.Consume()
			p.Card.stop()
		}
		if p.Card.action.Complete {
			CardManager.Move(p, p.player.Hand, 0)
		}
	}
}

func (p *PlayCard) Draw(win *pixelgl.Window) {
	if p.Card != nil {
		p.Card.Draw(win)
	}
}

func (p *PlayCard) AddCard(card *Card) {
	p.update = true
	if card != nil {
		p.Card = card
		p.Card.play(p.player)
	}
}

func (p *PlayCard) RemoveCard(i int) *Card {
	p.update = true
	if p.Card == nil {
		return nil
	}
	p.Card.stop()
	card := p.Card
	p.Card = nil
	return card
}

