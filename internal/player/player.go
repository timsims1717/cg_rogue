package player

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"golang.org/x/image/colornames"
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
	
	ui              *imdraw.IMDraw
}

func init() {
	Player1 = NewPlayer(nil)
}

func NewPlayer(character *floor.Character) *Player {
	return &Player{
		Character: character,
		Input:     input.NewInput(),
		ui:        imdraw.New(nil),
	}
}

func (p *Player) StartTurn() {
	p.Character.StartTurn()
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
	p.ui.Clear()
	// Health Bar
	p.ui.SetMatrix(camera.Cam.UITransform(pixel.V(camera.WindowWidthF - 75., camera.WindowHeightF - 20.), pixel.V(1., 1.), 0.))
	p.ui.Color = colornames.Darkgray
	p.ui.EndShape = imdraw.NoEndShape
	p.ui.Push(pixel.V(-60., -8.))
	p.ui.Push(pixel.V(-60., 8.))
	p.ui.Push(pixel.V(60., 8.))
	p.ui.Push(pixel.V(60., -8.))
	p.ui.Polygon(0.)
	perc := (120. / float64(p.Character.Health.MaxHP)) * float64(p.Character.Health.MaxHP - p.Character.Health.CurrHP)
	p.ui.Color = colornames.Darkred
	p.ui.EndShape = imdraw.NoEndShape
	p.ui.Push(pixel.V(-60. + perc, -8.))
	p.ui.Push(pixel.V(-60. + perc, 8.))
	p.ui.Push(pixel.V(60., 8.))
	p.ui.Push(pixel.V(60., -8.))
	p.ui.Polygon(0.)
	p.ui.Color = colornames.Lightgray
	p.ui.EndShape = imdraw.NoEndShape
	p.ui.Push(pixel.V(-60., -8.))
	p.ui.Push(pixel.V(-60., 8.))
	p.ui.Push(pixel.V(60., 8.))
	p.ui.Push(pixel.V(60., -8.))
	p.ui.Polygon(2.)
	// Defense Bar
	p.ui.SetMatrix(camera.Cam.UITransform(pixel.V(camera.WindowWidthF - 75., camera.WindowHeightF - 50.), pixel.V(1., 1.), 0.))
	p.ui.Color = colornames.Darkgray
	p.ui.EndShape = imdraw.NoEndShape
	p.ui.Push(pixel.V(-60., -8.))
	p.ui.Push(pixel.V(-60., 8.))
	p.ui.Push(pixel.V(60., 8.))
	p.ui.Push(pixel.V(60., -8.))
	p.ui.Polygon(0.)
	perc = (120. / float64(p.Character.Defense.MaxDef)) * float64(p.Character.Defense.MaxDef - p.Character.Defense.CurrDef)
	p.ui.Color = colornames.Mediumblue
	p.ui.EndShape = imdraw.NoEndShape
	p.ui.Push(pixel.V(-60. + perc, -8.))
	p.ui.Push(pixel.V(-60. + perc, 8.))
	p.ui.Push(pixel.V(60., 8.))
	p.ui.Push(pixel.V(60., -8.))
	p.ui.Polygon(0.)
	p.ui.Color = colornames.Lightgray
	p.ui.EndShape = imdraw.NoEndShape
	p.ui.Push(pixel.V(-60., -8.))
	p.ui.Push(pixel.V(-60., 8.))
	p.ui.Push(pixel.V(60., 8.))
	p.ui.Push(pixel.V(60., -8.))
	p.ui.Polygon(2.)
}

func (p *Player) Draw(win *pixelgl.Window) {
	// This is where the player's UI is drawn
	p.ui.Draw(win)
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
