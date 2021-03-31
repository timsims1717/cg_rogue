package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/internal/manager"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type PlayCard struct {
	player       *Player
	Card         *Card
	update       bool
	CurrSelector *selectors.AbstractSelector
}

func NewPlayCard(player *Player) *PlayCard {
	return &PlayCard{
		player: player,
	}
}

func (p *PlayCard) Update(turn bool) {
	if p.Card != nil && p.player != nil {
		if turn {
			if p.update {
				p.Card.setXY(pixel.V(camera.WindowWidthF - PlayRightPad, PlayBottomPad))
				p.Card.setScalar(PlayCardScale)
				p.update = false
			}
			if p.Card.PointInside(p.player.Input.World) && p.player.Input.Cancel.JustPressed() {
				p.player.Input.Cancel.Consume()
				p.CancelCard()
			}
		}
		p.Card.Values.Source = p.player.Character
		p.Card.Update(pixel.Rect{})
		if turn {
			if p.Card.played {
				if !manager.ActionManager.IsActing() {
					p.player.ActionsThisTurn++
					if p.Card.Rests > 0 {
						CardManager.Move(p, p.player.Discard, p.Card)
					} else {
						CardManager.Move(p, p.Card.Previous, p.Card)
					}
				}
			} else {
				if !p.Card.isPlay {
					p.Card.isPlay = true
					p.Card.actPtr = 0
					p.Card.played = false
					p.Card.Results = make([][]world.Coords, len(p.Card.Selectors))
					p.Card.tempOrig = []world.Coords{p.player.Character.GetCoords()}
					p.NextSelector()
				}
				if p.CurrSelector == nil {
					p.Card.played = true
					p.Card.Action.DoActions()
				} else {
					if !p.CurrSelector.IsCancelled() && !p.CurrSelector.IsDone() {
						p.CurrSelector.Selector.Update(p.player.Input)
					} else if p.CurrSelector.IsDone() {
						results := p.CurrSelector.Finish()
						p.Card.Results[p.Card.actPtr] = results
						moved := p.CurrSelector.IsMove
						p.CurrSelector = nil
						p.Card.actPtr++
						newOrig := p.player.Character.GetCoords()
						if moved && len(results) > 0 {
							newOrig = results[len(results)-1]
						}
						p.Card.tempOrig = append(p.Card.tempOrig, newOrig)
						p.NextSelector()
					} else if p.CurrSelector.IsCancelled() {
						p.Card.actPtr--
						p.Card.tempOrig = p.Card.tempOrig[:len(p.Card.tempOrig)-1]
						if p.Card.actPtr < 0 {
							p.CancelCard()
						} else {
							p.NextSelector()
						}
					}
				}
			}
		} else {
			if p.CurrSelector != nil {
				p.CancelCard()
			}
			CardManager.Move(p, p.Card.Previous, p.Card)
		}
	}
}

func (p *PlayCard) NextSelector() {
	if p.Card.actPtr >= len(p.Card.Selectors) {
		p.CurrSelector = nil
	} else {
		p.Card.Selectors[p.Card.actPtr].Reset(p.Card.tempOrig[len(p.Card.tempOrig)-1])
		p.Card.Selectors[p.Card.actPtr].Selector.SetValues(p.Card.Values)
		p.CurrSelector = p.Card.Selectors[p.Card.actPtr]
	}
}

func (p *PlayCard) Draw(win *pixelgl.Window) {
	if p.Card != nil {
		p.Card.Draw(win)
	}
}

func (p *PlayCard) CancelCard() {
	p.CurrSelector.Cancel()
	if p.Card != nil {
		CardManager.Move(p, p.Card.Previous, p.Card)
	}
}

func (p *PlayCard) AddCard(card *Card) {
	if p.Card != nil {
		p.CancelCard()
	}
	if p.player != nil && card != nil {
		p.update = true
		card.trans = true
		card.isPlay = false
		card.played = false
		p.Card = card
	}
}

func (p *PlayCard) RemoveCard(uuid uuid.UUID) *Card {
	p.update = true
	if p.Card == nil {
		return nil
	}
	p.CurrSelector = nil
	p.Card.isPlay = false
	p.Card.actPtr = -1
	p.Card.played = false
	card := p.Card
	p.Card = nil
	return card
}

