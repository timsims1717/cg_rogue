package player

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/input"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
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
}

func NewPlayer(character *characters.Character) *Player {
	return &Player{
		Character:        character,
		Input:            &input.Input{},
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
	p.Hand.Update(p.IsTurn)
	p.PlayCard.Update(p.IsTurn)
	p.Discard.Update(p.IsTurn)
	if p.IsTurn {
		if win.JustPressed(pixelgl.KeyA) {
			values := selectors.ActionValues{
				Source:  p.Character,
				Damage:  5,
				Move:    0,
				Range:   1,
				Targets: 1,
			}
			sel := selectors.NewTargetSelect()
			p.SetPlayerAction(NewPlayerAction(sel, values, BasicAttack))
		}
		if win.JustPressed(pixelgl.KeyM) {
			values := selectors.ActionValues{
				Source:  p.Character,
				Damage:  0,
				Move:    5,
				Range:   0,
				Targets: 0,
			}
			sel := selectors.NewPathSelect()
			sel.Unoccupied = true
			sel.Nonempty = true
			p.SetPlayerAction(NewPlayerAction(sel, values, BasicMove))
		}
		if win.JustPressed(pixelgl.KeyR) {
			values := selectors.ActionValues{}
			sel := selectors.NewNullSelect()
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
		occ := floor.CurrentFloor.GetOccupant(targets[0])
		if occ != nil {
			if target, ok := occ.(objects.Targetable); ok {
				actions.AddToBot(actions.NewDamageAction(values.Source, target, values.Damage))
			}
		}
	}
}

func (p *Player) Rest(_ []world.Coords, _ selectors.ActionValues) {
	actions.AddToBot(NewRestAction(p))
}