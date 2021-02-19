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
	act := player.NewPlayerAction(selectors.NewTargetSelect(), values, []func([]world.Coords, actions.ActionValues) {fn})
	return player.NewCard("Strike", "Deal 5 damage.", act)
}

func StrikeLevel(level int) actions.ActionValues {
	values := actions.ActionValues{
		Source:  nil,
		Damage:  4,
		Move:    0,
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
	act := player.NewPlayerAction(sel, values, []func([]world.Coords, actions.ActionValues) {fn})
	return player.NewCard("Maneuver", "Move 3.", act)
}

func ManeuverLevel(level int) actions.ActionValues {
	values := actions.ActionValues{
		Source:  nil,
		Damage:  0,
		Move:    3 + level,
		Range:   0,
		Targets: 0,
	}
	return values
}

func Charge(path []world.Coords, values actions.ActionValues) {

}