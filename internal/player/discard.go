package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
)

type Discard struct {
	player *Player
	Group  []*Card
	update bool
	Hover  bool
}

func NewDiscard(player *Player) *Discard {
	return &Discard{
		player: player,
		Group:  []*Card{},
	}
}

func (d *Discard) Update(turn bool) {
	// todo: turn will allow clicking the discard to see the cards in a grid view
	if d.player != nil {
		if len(d.Group) > 0 {
			top := d.Group[len(d.Group)-1]
			hover := top.PointInside(d.player.Input.World) && !top.trans
			if d.Hover != hover {
				d.Hover = hover
				d.update = true
			}
			for i, card := range d.Group {
				if d.update {
					if i == len(d.Group)-1 && d.Hover {
						card.setXY(pixel.V(cfg.WindowWidthF-DiscardRightPad, DiscardBottomPad*2.0))
						card.setScalar(DiscardHovScale)
					} else {
						card.setXY(pixel.V(cfg.WindowWidthF-DiscardRightPad, DiscardBottomPad))
						card.setScalar(DiscardScale)
					}
				}
				card.Update()
			}
		} else {
			d.Hover = false
		}
		if d.update {
			d.update = false
		}
	}
}

func (d *Discard) Draw(win *pixelgl.Window) {
	// loop backwards through Group and mark cards to draw until one isn't transitioning
	// if make it all the way through, draw the discard location as well
	drawBottom := true
	index := 0
	for i := len(d.Group)-1; i >= 0; i-- {
		card := d.Group[i]
		if !card.trans {
			drawBottom = false
			index = i
			break
		}
	}
	for i := index; i < len(d.Group); i++ {
		d.Group[i].Draw(win)
	}
	if drawBottom {
		// todo: draw discard card location
	}
}

func (d *Discard) AddCard(card *Card) {
	d.update = true
	if card != nil {
		card.trans = true
		d.Group = append(d.Group, card)
	}
}

func (d *Discard) RemoveTopCard() *Card {
	d.update = true
	if len(d.Group) > 0 {
		card := d.Group[len(d.Group)-1]
		d.Group = d.Group[:len(d.Group)-1]
		return card
	}
	return nil
}

func (d *Discard) RemoveCard(i int) *Card {
	d.update = true
	if i < 0 || i >= len(d.Group) {
		return nil
	}
	card := d.Group[i]
	d.Group = append(d.Group[0:i], d.Group[i+1:len(d.Group)]...)
	return card
}