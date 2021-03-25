package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
)

type Grid struct {
	player  *Player
	Index   int
	Group   []*Card
	Clicked []*Card
	update  bool
	Hovered int
	Show    bool
}

func NewGrid(player *Player) *Grid {
	return &Grid{
		player:  player,
		Group:   []*Card{},
		Hovered: -1,
	}
}

func (g *Grid) Update() {
	if g.player != nil && g.Show {
		if g.isHovered() {
			if !g.Group[g.Hovered].PointInside(g.player.Input.World) {
				g.Hovered = -1
				g.update = true
			}
		}
		if !g.isHovered() {
			for i := len(g.Group) - 1; i >= 0; i-- {
				card := g.Group[i]
				if card.PointInside(g.player.Input.World) && !card.trans {
					g.Hovered = i
					g.update = true
					break
				}
			}
		}
		for i, card := range g.Group {
			if g.update {
				card.setXY(GridLocation(i))
				if g.Hovered == i {
					card.setScalar(GridHovCardScale)
				} else {
					card.setScalar(GridCardScale)
				}
			}
			card.Update(pixel.Rect{})
		}
		if g.player.Input.Cancel.JustPressed() {
			g.player.Input.Cancel.Consume()
			g.ReturnCards()
		}
		if g.update {
			g.update = false
		}
	}
}

func (g *Grid) Draw(win *pixelgl.Window) {
	if g.Show {
		for _, card := range g.Group {
			card.Draw(win)
		}
	}
}

func (g *Grid) AddCard(card *Card) {
	g.update = true
	if card != nil {
		card.trans = true
		g.Group = append(g.Group, card)
	}
}

func (g *Grid) RemoveCard(id uuid.UUID) *Card {
	g.update = true
	index := -1
	for i, card := range g.Group {
		if card.ID == id {
			index = i
			break
		}
	}
	if index < 0 || index >= len(g.Group) {
		return nil
	}
	card := g.Group[index]
	g.Group = append(g.Group[0:index], g.Group[index+1:len(g.Group)]...)
	return card
}

func (g *Grid) ReturnCards() {
	for _, card := range g.Group {
		card.trans = true
	}
	if g.player != nil {
		g.player.Discard.InGrid = false
		g.player.Discard.update = true
	}
	g.Group = []*Card{}
	g.Show = false
}

func (g *Grid) isHovered() bool {
	return g.Hovered >= 0 && g.Hovered < len(g.Group)
}

// GridLocation returns the position of the card relative to the top center of
// the Grid (which should be the top center of the screen).
func GridLocation(i int) pixel.Vec {
	xPos := i % 5
	yPos := i / 5
	x := BaseCardWidth * float64(xPos - 2) + camera.WindowWidthF * 0.5
	y := camera.WindowHeightF - (BaseCardHeight * float64(yPos) + BaseCardHeight * 0.5)
	return pixel.V(x, y)
}