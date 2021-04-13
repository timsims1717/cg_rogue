package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
)

type Discard struct {
	player *Player
	Group  []*Card
	update bool
	Hover  bool
	InGrid bool
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
			hover := top.PointInside(d.player.Input.World) && !top.trans && !d.InGrid
			if d.Hover != hover {
				d.Hover = hover
				d.update = true
			}
			if turn && d.Hover && d.player.Input.Select.JustPressed() && d.player.Grid != nil {
				d.player.Input.Select.Consume()
				d.player.Grid.ReturnCards()
				d.InGrid = true
				d.player.Grid.Show = true
				for _, card := range d.Group {
					CardManager.Move(nil, d.player.Grid, card)
				}
			}
			for i, card := range d.Group {
				if !d.InGrid {
					if d.update {
						if i == len(d.Group)-1 && d.Hover {
							card.setXY(pixel.V(camera.WindowWidthF-DiscardRightPad, DiscardBottomPad*2.0))
							card.setScalar(DiscardHovScale)
						} else {
							card.setXY(pixel.V(camera.WindowWidthF-DiscardRightPad, DiscardBottomPad))
							card.setScalar(DiscardScale)
						}
					}
				}
				card.Update(pixel.Rect{})
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
	for i := len(d.Group) - 1; i >= 0; i-- {
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
	if d.InGrid {
		CardManager.Move(nil, d.player.Grid, card)
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

func (d *Discard) RemoveCard(id uuid.UUID) *Card {
	d.update = true
	index := -1
	for i, card := range d.Group {
		if card.ID == id {
			index = i
			break
		}
	}
	if index < 0 || index >= len(d.Group) {
		return nil
	}
	card := d.Group[index]
	d.Group = append(d.Group[0:index], d.Group[index+1:len(d.Group)]...)
	return card
}
