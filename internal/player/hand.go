package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
)

type Hand struct {
	player  *Player
	Group   []*Card
	update  bool
	Hovered int
}

func NewHand(player *Player) *Hand {
	return &Hand{
		player:  player,
		Group:   []*Card{},
		Hovered: -1,
	}
}

func (h *Hand) Update(turn bool) {
	if h.player != nil {
		if h.isHovered() {
			if !h.Group[h.Hovered].PointInside(h.player.Input.World) {
				h.Hovered = -1
				h.update = true
			}
		}
		if !h.isHovered() {
			for i := len(h.Group) - 1; i >= 0; i-- {
				card := h.Group[i]
				if card.PointInside(h.player.Input.World) && !card.trans {
					h.Hovered = i
					h.update = true
					break
				}
			}
		}
		for i, card := range h.Group {
			if h.update {
				card.setXY(HandLocation(i, h.Hovered))
				if h.Hovered == i {
					card.setScalar(HandHovCardScale)
				} else {
					card.setScalar(HandCardScale)
				}
			}
			card.Update(pixel.Rect{})
		}
		if h.isHovered() && h.player.Input.Select.JustPressed() && turn {
			h.player.Input.Select.Consume()
			//CardManager.Move(h.player.PlayCard, h, h.Group[h.Hovered])
			CardManager.Move(h, h.player.PlayCard, h.Group[h.Hovered])
		}
		if h.update {
			h.update = false
		}
	}
}

func (h *Hand) Draw(win *pixelgl.Window) {
	for _, card := range h.Group {
		card.Draw(win)
	}
}

func (h *Hand) AddCard(card *Card) {
	h.update = true
	if card != nil {
		card.trans = true
		h.Group = append(h.Group, card)
	}
}

func (h *Hand) RemoveCard(id uuid.UUID) *Card {
	h.update = true
	index := -1
	for i, card := range h.Group {
		if card.ID == id {
			index = i
			break
		}
	}
	if index < 0 || index >= len(h.Group) {
		return nil
	}
	card := h.Group[index]
	h.Group = append(h.Group[0:index], h.Group[index+1:len(h.Group)]...)
	return card
}

func (h *Hand) isHovered() bool {
	return h.Hovered >= 0 && h.Hovered < len(h.Group)
}

func HandLocation(i, hovered int) pixel.Vec {
	var offset float64
	if hovered == -1 || hovered == i {
		offset = 0.0
	} else if hovered > i {
		offset = -40.0
	} else {
		offset = 40.0
	}
	x := offset + HandLeftPad + (float64(i) * BaseCardWidth * HandCardScale * 0.85)
	y := HandBottomPad
	if hovered == i {
		y *= 2.0
	}
	return pixel.V(x, y)
}
