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
	Character  *characters.Character
	CurrAction *PlayerAction
	Input      *input.Input
	Hand       *Hand
	PlayCard   *PlayCard
}

func NewPlayer(character *characters.Character) *Player {
	return &Player{
		Character:        character,
		Input:            &input.Input{},
	}
}

func (p *Player) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyA) {
		values := actions.ActionValues{
			Source:  p.Character,
			Damage:  5,
			Move:    0,
			Range:   1,
			Targets: 1,
		}
		sel := selectors.NewTargetSelect()
		p.SetPlayerAction(NewPlayerAction(sel, values, []func([]world.Coords, actions.ActionValues) {BasicAttack}))
	}
	if win.JustPressed(pixelgl.KeyM) {
		values := actions.ActionValues{
			Source:  p.Character,
			Damage:  0,
			Move:    5,
			Range:   0,
			Targets: 0,
		}
		sel := selectors.NewPathSelect()
		sel.Unoccupied = true
		sel.Nonempty = true
		p.SetPlayerAction(NewPlayerAction(sel, values, []func([]world.Coords, actions.ActionValues) {BasicMove}))
	}
	if p.CurrAction != nil {
		p.CurrAction.Update()
		if p.CurrAction.Complete {
			p.CurrAction = nil
		}
	}
}

func (p *Player) SetPlayerAction(act *PlayerAction) {
	act.Complete = false
	act.Selector.Init(p.Input)
	act.Selector.SetValues(act.Values)
	p.CurrAction = act
}

func BasicMove(path []world.Coords, values actions.ActionValues) {
	actions.AddToBot(actions.NewMoveSeriesAction(values.Source, path))
}

func BasicAttack(targets []world.Coords, values actions.ActionValues) {
	if len(targets) > 0 {
		occ := floor.CurrentFloor.GetOccupant(targets[0])
		if occ != nil {
			if target, ok := occ.(objects.Targetable); ok {
				actions.AddToBot(actions.NewDamageAction(values.Source, target, values.Damage))
			}
		}
	}
}