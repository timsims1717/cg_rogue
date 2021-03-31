package cards

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/manager"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type Thrust struct {
	*player.Card
}

func (c *Thrust) DoActions() {
	manager.ActionManager.AddToBot(actions.NewDamageAction(c.Results[0], c.Values))
}

func (c *Thrust) SetValues(level int) {
	values := selectors.ActionValues{
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
	c.Values = values
}

func (c *Thrust) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewTargetSelect(),
	}
}

func (c *Thrust) SetCard(card *player.Card) {
	c.Card = card
}

func CreateThrust() *player.Card {
	return player.NewCard("Thrust", "Deal 4 damage.", &Thrust{})
}

type Dash struct {
	*player.Card
}

func (c *Dash) DoActions() {
	manager.ActionManager.AddToBot(actions.NewMoveSeriesAction(c.Values.Source, c.Values.Source, c.Results[0]))
}

func (c *Dash) SetValues(level int) {
	values := selectors.ActionValues{
		Move: 4 + level,
		Checks: floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    true,
			NonEmpty:      true,
			EndUnoccupied: true,
			Orig:          world.Coords{},
		},
	}
	c.Values = values
}

func (c *Dash) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewPathSelect(),
	}
}

func (c *Dash) SetCard(card *player.Card) {
	c.Card = card
}

func CreateDash() *player.Card {
	return player.NewCard("Dash", "Move 4.", &Dash{})
}

//func CreateQuickStrike() *player.Card {
//	valMov, valAtk := QuickStrikeLevel(0)
//	moveAction := player.NewCardAction(selectors.NewPathSelect(), valMov, "Move 1.")
//	attackAction := player.NewCardAction(selectors.NewTargetSelect(), valAtk, "Deal 3 damage.")
//	qs1 := &
//	return player.NewCard("Quick Strike", []*player.AbstractPlayerAction{moveAction, attackAction})
//}
//
//func QuickStrikeLevel(level int) (selectors.ActionValues, selectors.ActionValues) {
//	valMov := selectors.ActionValues{
//		Move:    1,
//		Checks: floor.PathChecks{
//			NotFilled:     true,
//			Unoccupied:    true,
//			NonEmpty:      true,
//			EndUnoccupied: true,
//			Orig:          world.Coords{},
//		},
//	}
//	valAtk := selectors.ActionValues{
//		Damage:  3,
//		Range:   1,
//		Targets: 1,
//	}
//	if level >= 1 {
//		valAtk.Damage += 2
//	}
//	if level >= 2 {
//		valMov.Move += 1
//	}
//	if level >= 3 {
//		valAtk.Damage += 2
//	}
//	if level >= 4 {
//		valAtk.Damage += 1
//		valMov.Move += 1
//	}
//	if level >= 5 {
//		valAtk.Damage += 1
//		valMov.Move += 1
//	}
//	return valMov, valAtk
//}
//
//func CreateSweep() *player.Card {
//	values := SweepLevel(0)
//	fn := func (area []world.Coords, values selectors.ActionValues) {
//		manager.ActionManager.AddToBot(actions.NewPushMultiAction(area, values))
//	}
//	act := player.NewPlayerAction(selectors.NewArcSelect(), values, fn)
//	sec := player.NewCardSection("Deal 2 damage and push 1 away.", act)
//	return player.NewCard("Sweep", []*player.AbstractCardSection{sec})
//}
//
//func SweepLevel(level int) selectors.ActionValues {
//	values := selectors.ActionValues{
//		Damage:   2,
//		Range:    1,
//		Strength: 1,
//		Targets:  3,
//		Checks: floor.PathChecks{
//			NotFilled:     true,
//			Unoccupied:    false,
//			NonEmpty:      false,
//			EndUnoccupied: false,
//			Orig:          world.Coords{},
//		},
//	}
//	return values
//}
//
//func CreateVault() *player.Card {
//	valMov, valAtk := VaultLevel(0)
//	fnMov := func (path []world.Coords, values selectors.ActionValues) {
//		h := values.Source.GetCoords()
//		if len(path) > 0 {
//			h = path[len(path)-1]
//		}
//		manager.ActionManager.AddToBot(actions.NewMoveAction(values.Source, values.Source, h))
//	}
//	fnAtk := func (targets []world.Coords, values selectors.ActionValues) {
//		manager.ActionManager.AddToBot(actions.NewDamageHexAction(targets, values))
//	}
//	selMov := selectors.NewEmptyHexSelect()
//	selAtk := selectors.NewHexSelect()
//	actMov := player.NewPlayerAction(selMov, valMov, fnMov)
//	actAtk := player.NewPlayerAction(selAtk, valAtk, fnAtk)
//	secMov := player.NewCardSection("Jump 2.", actMov)
//	secAtk := player.NewCardSection("Deal 1 damage.", actAtk)
//	return player.NewCard("Vault", []*player.AbstractCardSection{secMov, secAtk})
//}
//
//func VaultLevel(level int) (selectors.ActionValues, selectors.ActionValues) {
//	valMov := selectors.ActionValues{
//		Move:    2,
//		Targets: 1,
//	}
//	valAtk := selectors.ActionValues{
//		Damage:  1 + (level + 1) / 2,
//		Range:   1,
//		Targets: 1 + level / 2,
//	}
//	return valMov, valAtk
//}
//
//func CreateDaggerThrow() *player.Card {
//	values := DaggerThrowLevel(0)
//	fn := func (targets []world.Coords, values selectors.ActionValues) {
//		manager.ActionManager.AddToBot(actions.NewDamageHexAction(targets, values))
//	}
//	sel := selectors.NewLineSelect()
//	act := player.NewPlayerAction(sel, values, fn)
//	sec := player.NewCardSection("Deal 2 damage within 4.", act)
//	return player.NewCard("Dagger Throw", []*player.AbstractCardSection{sec})
//}
//
//func DaggerThrowLevel(level int) selectors.ActionValues {
//	values := selectors.ActionValues{
//		Damage:  2 + (level + 1) / 2,
//		Range:   4 + (level) / 2,
//		Targets: 1,
//		Checks: floor.PathChecks{
//			NotFilled:     true,
//			Unoccupied:    false,
//			NonEmpty:      false,
//			EndUnoccupied: false,
//			Orig:          world.Coords{},
//		},
//	}
//	return values
//}