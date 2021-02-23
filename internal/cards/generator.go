package cards

import (
	"github.com/timsims1717/cg_rogue_go/internal/player"
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

//func (c *Constructor) Build() *player.Card {
//	if c.count > 0 {
//
//	} else {
//		sec := player.NewCardSection("[NO DESCRIPTION]", nil)
//		return player.NewCard(c.title, []*player.CardSection{sec})
//	}
//}