package cards

import (
	"fmt"
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type Thrust struct {
	*player.Card
}

func (c *Thrust) DoActions() {
	AddToBot(actions.NewDamageAction(c.Results[0][0].Area, c.Values), c.Results[0][0].Effect)
}

func (c *Thrust) SetValues(level int) {
	values := selector.ActionValues{
		Damage:  5 + level*2,
		Range:   1,
		Targets: 1,
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Deal %d damage.", values.Damage)
}

func (c *Thrust) InitSelectors() {
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.TargetSelect{
			Effect: selector.NewSelectionEffect(&selector.AttackTriangleEffect{}, c.Values),
		}, false),
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
	AddToBot(actions.NewMoveSeriesAction(c.Values.Source, c.Values.Source, c.Results[0][0].Area), c.Results[0][0].Effect)
}

func (c *Dash) SetValues(level int) {
	values := selector.ActionValues{
		Move: 5 + level,
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Move %d.", values.Move)
}

func (c *Dash) InitSelectors() {
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.PathSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    true,
				NonEmpty:      true,
				EndUnoccupied: true,
			},
			Effect: selector.NewSelectionEffect(&selector.MoveSeriesEffect{}, c.Values),
		}, true),
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
	AddToBot(actions.NewMoveSeriesAction(c.Values.Source, c.Values.Source, c.Results[0][0].Area), c.Results[0][0].Effect)
	AddToBot(actions.NewDamageAction(c.Results[1][0].Area, c.Values), c.Results[1][0].Effect)
}

func (c *QuickStrike) SetValues(level int) {
	values := selector.ActionValues{
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
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.PathSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    true,
				NonEmpty:      true,
				EndUnoccupied: true,
			},
			Effect: selector.NewSelectionEffect(&selector.MoveSeriesEffect{}, c.Values),
		}, true),
		selector.NewSelector(&selector.TargetSelect{
			Effect: selector.NewSelectionEffect(&selector.AttackTriangleEffect{}, c.Values),
		}, false),
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
	AddToBot(actions.NewPushMultiAction(c.Results[0][0].Area, c.Values), c.Results[0][0].Effect)
}

func (c *Sweep) SetValues(level int) {
	values := selector.ActionValues{
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
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.ArcSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    false,
				NonEmpty:      false,
				EndUnoccupied: false,
			},
			Effect: selector.NewSelectionEffect(&selector.AttackTriangleEffect{}, c.Values),
		}, false),
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
		h = c.Results[0][0].Area[len(c.Results[0])-1]
	}
	AddToBot(actions.NewMoveAction(c.Values.Source, c.Values.Source, h), c.Results[0][0].Effect)
	AddToBot(actions.NewDamageHexAction(c.Results[1][0].Area, c.Values), c.Results[1][0].Effect)
}

func (c *Vault) SetValues(level int) {
	values := selector.ActionValues{
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
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.HexSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    true,
				NonEmpty:      true,
				EndUnoccupied: true,
			},
			Effect: selector.NewSelectionEffect(&selector.MoveEffect{}, c.Values),
		}, true),
		selector.NewSelector(&selector.HexSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    false,
				NonEmpty:      false,
				EndUnoccupied: false,
			},
			Effect: selector.NewSelectionEffect(&selector.AttackTriangleEffect{}, c.Values),
		}, false),
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
	AddToBot(actions.NewDamageHexAction(c.Results[0][0].Area, c.Values), c.Results[0][0].Effect, c.Results[0][1].Effect)
}

func (c *DaggerThrow) SetValues(level int) {
	values := selector.ActionValues{
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
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.LineTargetSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    false,
				NonEmpty:      false,
				EndUnoccupied: false,
			},
			Effect:    selector.NewSelectionEffect(&selector.AttackTargetEffect{}, c.Values),
			SecEffect: selector.NewSelectionEffect(&selector.HighlightEffect{}, c.Values),
		}, false),
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
	AddToBot(actions.NewDamageHexAction(c.Results[0][0].Area, c.Values), c.Results[0][0].Effect)
	AddToBot(actions.NewMoveSeriesAction(c.Values.Source, c.Values.Source, c.Results[1][0].Area), c.Results[1][0].Effect)
}

func (c *Disengage) SetValues(level int) {
	values := selector.ActionValues{
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
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.ArcSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    false,
				NonEmpty:      false,
				EndUnoccupied: false,
			},
			Effect: selector.NewSelectionEffect(&selector.AttackTriangleEffect{}, c.Values),
		}, false),
		selector.NewSelector(&selector.PathSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    true,
				NonEmpty:      true,
				EndUnoccupied: true,
			},
			Effect: selector.NewSelectionEffect(&selector.MoveSeriesEffect{}, c.Values),
		}, true),
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
		h = c.Results[0][0].Area[len(c.Results[0][0].Area)-1]
	}
	AddToBot(actions.NewSlamAction(h, c.Results[0][1].Area, c.Values), c.Results[0][0].Effect, c.Results[0][1].Effect)
}

func (c *Slam) SetValues(level int) {
	values := selector.ActionValues{
		Move:    1,
		Damage:  1,
		Targets: 1,
		Area:    world.Spiral(6),
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
		values.Damage += 1
	}
	c.Values = values
	c.RawDesc = fmt.Sprintf("Jump %d. Deal %d dmg w/in 1.", values.Move, values.Damage)
}

func (c *Slam) InitSelectors() {
	c.Selectors = []*selector.AbstractSelector{
		selector.NewSelector(&selector.HexAreaSplitSelect{
			PathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    true,
				NonEmpty:      true,
				EndUnoccupied: true,
			},
			SecPathChecks: floor.PathChecks{
				NotFilled:     true,
				Unoccupied:    false,
				NonEmpty:      false,
				EndUnoccupied: false,
			},
			Effect:    selector.NewSelectionEffect(&selector.MoveEffect{}, c.Values),
			SecEffect: selector.NewSelectionEffect(&selector.AttackTriangleEffect{}, c.Values),
		}, true),
	}
}

func (c *Slam) SetCard(card *player.Card) {
	c.Card = card
}

func CreateSlam() *player.Card {
	return player.NewCard("Slam", "Jump 1. Deal 1 dmg w/in 1.", &Slam{})
}
