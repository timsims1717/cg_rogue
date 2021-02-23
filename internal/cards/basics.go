package cards

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

func CreateStrike() *player.Card {
	values := StrikeLevel(0)
	fn := func (targets []world.Coords, values actions.ActionValues) {
		if len(targets) > 0 {
			occ := floor.CurrentFloor.GetOccupant(targets[0])
			if occ != nil {
				if target, ok := occ.(objects.Targetable); ok {
					actions.AddToBot(actions.NewDamageAction(values.Source, target, values.Damage))
				}
			}
		}
	}
	act := player.NewPlayerAction(selectors.NewTargetSelect(), values, fn)
	sec := player.NewCardSection("Deal 4 damage.", act)
	return player.NewCard("Strike", []*player.CardSection{sec})
}

func StrikeLevel(level int) actions.ActionValues {
	values := actions.ActionValues{
		Damage:  4,
		Range:   1,
		Targets: 1,
	}
	if level >= 1 {
		values.Damage += 2
	}
	if level >= 2 {
		values.Targets += 1
	}
	if level >= 3 {
		values.Damage += 3
	}
	if level >= 4 {
		values.Targets += 1
	}
	if level >= 5 {
		values.Damage += 4
	}
	return values
}

func CreateManeuver() *player.Card {
	values := ManeuverLevel(0)
	fn := func (path []world.Coords, values actions.ActionValues) {
		actions.AddToBot(actions.NewMoveSeriesAction(values.Source, path))
	}
	sel := selectors.NewPathSelect()
	sel.Unoccupied = true
	sel.Nonempty = true
	act := player.NewPlayerAction(sel, values, fn)
	sec := player.NewCardSection("Move 3.", act)
	return player.NewCard("Maneuver", []*player.CardSection{sec})
}

func ManeuverLevel(level int) actions.ActionValues {
	values := actions.ActionValues{
		Move:    3 + level,
	}
	return values
}

func CreateCharge() *player.Card {
	valMov, valAtk := ChargeLevel(0)
	fnMov := func (path []world.Coords, values actions.ActionValues) {
		actions.AddToBot(actions.NewMoveSeriesAction(values.Source, path))
	}
	fnAtk := func (targets []world.Coords, values actions.ActionValues) {
		if len(targets) > 0 {
			occ := floor.CurrentFloor.GetOccupant(targets[0])
			if occ != nil {
				if target, ok := occ.(objects.Targetable); ok {
					actions.AddToBot(actions.NewDamageAction(values.Source, target, values.Damage))
				}
			}
		}
	}
	selMov := selectors.NewPathSelect()
	selMov.Unoccupied = true
	selMov.Nonempty = true
	actMov := player.NewPlayerAction(selMov, valMov, fnMov)
	actAtk := player.NewPlayerAction(selectors.NewTargetSelect(), valAtk, fnAtk)
	secMov := player.NewCardSection("Move 1.", actMov)
	secAtk := player.NewCardSection("Deal 2 damage.", actAtk)
	return player.NewCard("Charge", []*player.CardSection{secMov, secAtk})
}

func ChargeLevel(level int) (actions.ActionValues, actions.ActionValues) {
	valMov := actions.ActionValues{
		Move:    1,
	}
	valAtk := actions.ActionValues{
		Damage:  2,
		Range:   1,
		Targets: 1,
	}
	if level >= 1 {
		valAtk.Damage += 2
	}
	if level >= 2 {
		valMov.Move += 1
	}
	if level >= 3 {
		valAtk.Damage += 2
	}
	if level >= 4 {
		valAtk.Damage += 1
		valMov.Move += 1
	}
	if level >= 5 {
		valAtk.Damage += 1
		valMov.Move += 1
	}
	return valMov, valAtk
}