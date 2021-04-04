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
		Damage:  5 + level * 2,
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
	return player.NewCard("Thrust", "Deal 5 damage.", &Thrust{})
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
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Move %d.", values.Move)
}

func (c *Dash) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewPathSelect(true, floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    true,
			NonEmpty:      true,
			EndUnoccupied: true,
			Orig:          world.Coords{},
		}),
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
		selectors.NewPathSelect(true, floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    true,
			NonEmpty:      true,
			EndUnoccupied: true,
			Orig:          world.Coords{},
		}),
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
		Damage:   3,
		Range:    1,
		Strength: 1,
		Targets:  3,
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
		selectors.NewArcSelect(floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    false,
			NonEmpty:      false,
			EndUnoccupied: false,
			Orig:          world.Coords{},
		}),
	}
}

func (c *Sweep) SetCard(card *player.Card) {
	c.Card = card
}

func CreateSweep() *player.Card {
	return player.NewCard("Sweep", "Deal 3 damage and push 1 away.", &Sweep{})
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
		Damage:  2,
	}
	if level >= 1 {
		values.Move += 1
	}
	if level >= 2 {
		values.Damage += 1
	}
	if level >= 3 {
		values.Move += 1
	}
	if level >= 4 {
		values.Damage += 1
	}
	if level >= 5 {
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
	return player.NewCard("Vault", "Jump 2. Deal 2 damage.", &Vault{})
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
		Range:   5,
		Targets: 1,
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
		selectors.NewLineSelect(false, floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    false,
			NonEmpty:      false,
			EndUnoccupied: false,
			Orig:          world.Coords{},
		}),
	}
}

func (c *DaggerThrow) SetCard(card *player.Card) {
	c.Card = card
}

func CreateDaggerThrow() *player.Card {
	return player.NewCard("Dagger Throw", "Deal 2 damage within 5.", &DaggerThrow{})
}

type Disengage struct {
	*player.Card
}

func (c *Disengage) DoActions() {
	manager.ActionManager.AddToBot(actions.NewDamageHexAction(c.Results[0], c.Values))
	manager.ActionManager.AddToBot(actions.NewMoveSeriesAction(c.Values.Source, c.Values.Source, c.Results[1]))
}

func (c *Disengage) SetValues(level int) {
	values := selectors.ActionValues{
		Move:    3,
		Damage:  1,
		Targets: 1,
		Range:   1,
	}
	if level >= 1 {
		values.Move += 1
	}
	if level >= 2 {
		values.Damage += 1
	}
	if level >= 3 {
		values.Targets += 1
	}
	if level >= 4 {
		values.Damage += 1
	}
	if level >= 5 {
		values.Targets += 1
		values.Move += 1
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Deal %d damage. Move %d.", values.Damage, values.Move)
}

func (c *Disengage) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewArcSelect(floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    false,
			NonEmpty:      false,
			EndUnoccupied: false,
			Orig:          world.Coords{},
		}),
		selectors.NewPathSelect(true, floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    true,
			NonEmpty:      true,
			EndUnoccupied: true,
			Orig:          world.Coords{},
		}),
	}
}

func (c *Disengage) SetCard(card *player.Card) {
	c.Card = card
}

func CreateDisengage() *player.Card {
	return player.NewCard("Disengage", "Deal 1 damage. Move 3.", &Disengage{})
}

type Slam struct {
	*player.Card
}

func (c *Slam) DoActions() {
	h := c.Values.Source.GetCoords()
	if len(c.Results[0]) > 0 {
		h = c.Results[0][len(c.Results[0])-1]
	}
	manager.ActionManager.AddToBot(actions.NewMoveAction(c.Values.Source, c.Values.Source, h))
	//manager.ActionManager.AddToBot(actions.NewMoveSeriesAction(c.Values.Source, c.Values.Source, c.Results[1]))
}

func (c *Slam) SetValues(level int) {
	values := selectors.ActionValues{
		Move:    2,
		Damage:  2,
		Range:   1,
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
		values.Range += 1
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Jump %d. Deal %d dmg w/in %d.", values.Move, values.Damage, values.Range)
}

func (c *Slam) InitSelectors() {
	c.Selectors = []*selectors.AbstractSelector{
		selectors.NewMoveHexSelect(),
	}
}

func (c *Slam) SetCard(card *player.Card) {
	c.Card = card
}

func CreateSlam() *player.Card {
	return player.NewCard("Slam", "Jump 2. Deal 2 dmg w/in 1.", &Slam{})
}
