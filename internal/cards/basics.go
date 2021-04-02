package cards

import (
	"fmt"
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
		Damage:  4 + level * 2,
		Range:   1,
		Targets: 1,
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Deal %d damage.", values.Damage)
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
	c.RawDesc = fmt.Sprintf("Move %d.", values.Move)
}

func (c *Dash) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewPathSelect(true),
	}
}

func (c *Dash) SetCard(card *player.Card) {
	c.Card = card
}

func CreateDash() *player.Card {
	return player.NewCard("Dash", "Move 4.", &Dash{})
}

type QuickStrike struct {
	*player.Card
}

func (c *QuickStrike) DoActions() {
	manager.ActionManager.AddToBot(actions.NewMoveSeriesAction(c.Values.Source, c.Values.Source, c.Results[0]))
	manager.ActionManager.AddToBot(actions.NewDamageAction(c.Results[1], c.Values))
}

func (c *QuickStrike) SetValues(level int) {
	values := selectors.ActionValues{
		Move:    1,
		Damage:  3,
		Range:   1,
		Targets: 1,
		Checks: floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    true,
			NonEmpty:      true,
			EndUnoccupied: true,
			Orig:          world.Coords{},
		},
	}
	if level >= 1 {
		values.Damage += 1
	}
	if level >= 2 {
		values.Move += 1
	}
	if level >= 3 {
		values.Damage += 1
	}
	if level >= 4 {
		values.Damage += 1
		values.Move += 1
	}
	if level >= 5 {
		values.Damage += 1
		values.Move += 1
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Move %d. Deal %d damage.", values.Move, values.Damage)
}

func (c *QuickStrike) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewPathSelect(true),
		selectors.NewTargetSelect(),
	}
}

func (c *QuickStrike) SetCard(card *player.Card) {
	c.Card = card
}

func CreateQuickStrike() *player.Card {
	return player.NewCard("Quick Strike", "Move 1. Deal 3 damage.", &QuickStrike{})
}

type Sweep struct {
	*player.Card
}

func (c *Sweep) DoActions() {
	manager.ActionManager.AddToBot(actions.NewPushMultiAction(c.Results[0], c.Values))
}

func (c *Sweep) SetValues(level int) {
	values := selectors.ActionValues{
		Damage:   2,
		Range:    1,
		Strength: 1,
		Targets:  3,
		Checks: floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    false,
			NonEmpty:      false,
			EndUnoccupied: false,
			Orig:          world.Coords{},
		},
	}
	if level >= 1 {
		values.Damage += 1
	}
	if level >= 2 {
		values.Targets += 1
	}
	if level >= 3 {
		values.Strength += 1
	}
	if level >= 4 {
		values.Damage += 1
	}
	if level >= 5 {
		values.Targets += 1
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Deal %d damage and push %d away.", values.Damage, values.Strength)
}

func (c *Sweep) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewArcSelect(),
	}
}

func (c *Sweep) SetCard(card *player.Card) {
	c.Card = card
}

func CreateSweep() *player.Card {
	return player.NewCard("Sweep", "Deal 2 damage and push 1 away.", &Sweep{})
}

type Vault struct {
	*player.Card
}

func (c *Vault) DoActions() {
	h := c.Values.Source.GetCoords()
	if len(c.Results[0]) > 0 {
		h = c.Results[0][len(c.Results[0])-1]
	}
	manager.ActionManager.AddToBot(actions.NewMoveAction(c.Values.Source, c.Values.Source, h))
	manager.ActionManager.AddToBot(actions.NewDamageHexAction(c.Results[1], c.Values))
}

func (c *Vault) SetValues(level int) {
	values := selectors.ActionValues{
		Move:    2,
		Range:   1,
		Targets: 1,
		Damage:  1,
	}
	if level >= 1 {
		values.Damage += 1
	}
	if level >= 2 {
		values.Move += 1
	}
	if level >= 3 {
		values.Damage += 1
	}
	if level >= 4 {
		values.Move += 1
	}
	if level >= 5 {
		values.Move += 1
		values.Damage += 1
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Jump %d. Deal %d damage.", values.Move, values.Damage)
}

func (c *Vault) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewMoveHexSelect(),
		selectors.NewHexSelect(false),
	}
}

func (c *Vault) SetCard(card *player.Card) {
	c.Card = card
}

func CreateVault() *player.Card {
	return player.NewCard("Vault", "Jump 2. Deal 1 damage.", &Vault{})
}

type DaggerThrow struct {
	*player.Card
}

func (c *DaggerThrow) DoActions() {
	manager.ActionManager.AddToBot(actions.NewDamageHexAction(c.Results[0], c.Values))
}

func (c *DaggerThrow) SetValues(level int) {
	values := selectors.ActionValues{
		Damage:  2,
		Range:   4,
		Targets: 1,
		Checks: floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    false,
			NonEmpty:      false,
			EndUnoccupied: false,
			Orig:          world.Coords{},
		},
	}
	if level >= 1 {
		values.Damage += 1
	}
	if level >= 2 {
		values.Range += 1
	}
	if level >= 3 {
		values.Damage += 1
	}
	if level >= 4 {
		values.Range += 1
	}
	if level >= 5 {
		values.Range += 1
		values.Damage += 1
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Deal %d damage within %d.", values.Damage, values.Range)
}

func (c *DaggerThrow) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewLineSelect(false),
	}
}

func (c *DaggerThrow) SetCard(card *player.Card) {
	c.Card = card
}

func CreateDaggerThrow() *player.Card {
	return player.NewCard("Dagger Throw", "Deal 2 damage within 4.", &DaggerThrow{})
}
