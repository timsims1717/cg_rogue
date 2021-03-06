package cards

import (
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type Constructor struct {
	count    int
	title    string
	sections []*player.CardSection
}

func NewConstructor() *Constructor {
	return &Constructor{
		count: 0,
		title: "[MISSING TITLE]",
		sections: []*player.CardSection{},
	}
}

func (c *Constructor) Build() *player.Card {
	if c.count > 0 {
		return player.NewCard(c.title, c.sections)
	} else {
		sec := player.NewCardSection("[NO DESCRIPTION]", nil)
		return player.NewCard(c.title, []*player.CardSection{sec})
	}
}

func (c *Constructor) AddSection(text string, sel selectors.Selector, val selectors.ActionValues, fn func([]world.Coords, selectors.ActionValues)) {
	c.sections = append(c.sections, player.NewCardSection(text, player.NewPlayerAction(sel, val, fn)))
}