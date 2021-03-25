package player

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/state"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
)

var Player1 *Player

type Player struct {
	Character       *characters.Character
	CurrAction      *PlayerAction
	Input           *input.Input
	Hand            *Hand
	PlayCard        *PlayCard
	Discard         *Discard
	Grid            *Grid
	ActionsThisTurn int
	IsTurn          bool
	RestButton      *ui.ActionEl
	MoveButton      *ui.ActionEl
}

func init() {
	Player1 = NewPlayer(nil)
}

func NewPlayer(character *characters.Character) *Player {
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
	p.CurrAction = nil
}

func (p *Player) Update(win *pixelgl.Window) {
	if state.Machine.State == state.InGame {
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
		if p.RestButton != nil {
			p.RestButton.Disabled = !p.IsTurn
			p.RestButton.Update(p.Input)
		}
		if p.MoveButton != nil {
			p.MoveButton.Disabled = !p.IsTurn
			p.MoveButton.Update(p.Input)
		}
		if p.IsTurn {
			if p.CurrAction != nil && p.PlayCard.Card == nil && p.Input.Cancel.JustPressed() {
				p.Input.Cancel.Consume()
				p.CurrAction = nil
			}
			if p.CurrAction != nil {
				p.CurrAction.Update()
				if p.CurrAction.Complete {
					if p.CurrAction.Complete {
						p.ActionsThisTurn += 1
					}
					p.CurrAction = nil
				}
			}
		}
	}
}

func (p *Player) Draw(win *pixelgl.Window) {
	if p.RestButton != nil {
		p.RestButton.Draw(win)
	}
	if p.MoveButton != nil {
		p.MoveButton.Draw(win)
	}
}

func (p *Player) SetPlayerAction(act *PlayerAction) {
	act.Complete = false
	act.Selector.Init(p.Input)
	act.Selector.SetValues(act.Values)
	p.CurrAction = act
}