package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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

func (h *Hand) Update() {
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
		card.Update()
	}
	if h.isHovered() && h.player.Input.Select.JustPressed() {
		h.player.Input.Select.Consume()
		CardManager.Move(h.player.PlayCard, h, h.Hovered)
		CardManager.Move(h, h.player.PlayCard, h.Hovered)
	}
	if h.update {
		h.update = false
	}
}

func (h *Hand) Draw(win *pixelgl.Window) {
	for i, card := range h.Group {
		if i != h.Hovered {
			card.Draw(win)
		}
	}
	if h.isHovered() {
		h.Group[h.Hovered].Draw(win)
	}
}

func (h *Hand) AddCard(card *Card) {
	h.update = true
	if card != nil {
		card.trans = true
		h.Group = append(h.Group, card)
	}
}

func (h *Hand) RemoveCard(i int) *Card {
	h.update = true
	if i < 0 || i >= len(h.Group) {
		return nil
	}
	card := h.Group[i]
	h.Group = append(h.Group[0:i], h.Group[i+1:len(h.Group)]...)
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
		offset = -20.0
	} else {
		offset = 20.0
	}
	x := offset + CardStart + (float64(i) * BaseCardWidth * HandCardScale * 0.85)
	y := BaseCardHeight * 0.25 * HandCardScale
	if hovered == i {
		y *= 2.0
	}
	return pixel.V(x, y)
}