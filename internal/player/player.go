package player

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/internal/state"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

var Player1 *Player

type Player struct {
	Character       *characters.Character
	CurrAction      *PlayerAction
	Input           *input.Input
	Hand            *Hand
	PlayCard        *PlayCard
	Discard         *Discard
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
		Input:     &input.Input{},
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
			if win.JustPressed(pixelgl.KeyA) {
				values := selectors.ActionValues{
					Source:  p.Character,
					Damage:  1,
					Move:    0,
					Range:   1,
					Targets: 1,
				}
				sel := selectors.NewTargetSelect()
				p.PlayCard.CancelCard()
				p.SetPlayerAction(NewPlayerAction(sel, values, BasicAttack))
			}
			if win.JustPressed(pixelgl.KeyM) {
				values := selectors.ActionValues{
					Source:  p.Character,
					Damage:  0,
					Move:    1,
					Range:   0,
					Targets: 0,
					Checks: floor.PathChecks{
						NotFilled:     true,
						Unoccupied:    true,
						NonEmpty:      true,
						EndUnoccupied: true,
						Orig:          world.Coords{},
					},
				}
				sel := selectors.NewPathSelect()
				p.PlayCard.CancelCard()
				p.SetPlayerAction(NewPlayerAction(sel, values, BasicMove))
			}
			if win.JustPressed(pixelgl.KeyR) {
				values := selectors.ActionValues{}
				sel := selectors.NewNullSelect()
				p.PlayCard.CancelCard()
				p.SetPlayerAction(NewPlayerAction(sel, values, p.Rest))
			}
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

func BasicMove(path []world.Coords, values selectors.ActionValues) {
	actions.AddToBot(actions.NewMoveSeriesAction(values.Source, values.Source, path))
}

func BasicAttack(targets []world.Coords, values selectors.ActionValues) {
	if len(targets) > 0 {
		actions.AddToBot(actions.NewDamageAction(targets, values))
	}
}

func (p *Player) Rest(_ []world.Coords, _ selectors.ActionValues) {
	actions.AddToBot(NewRestAction(p))
}