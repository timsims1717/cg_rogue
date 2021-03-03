package player

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"reflect"
)

type PlayerAction struct {
	Selector selectors.Selector
	Action   func([]world.Coords, actions.ActionValues)
	Values   actions.ActionValues
	Complete bool
	Cancel   bool
}

func NewPlayerAction(sel selectors.Selector, values actions.ActionValues, act func([]world.Coords, actions.ActionValues)) *PlayerAction {
	if sel == nil {
		return nil
	}
	switch reflect.TypeOf(sel).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(sel).IsNil() {
			return nil
		}
	}
	return &PlayerAction{
		Selector: sel,
		Action:   act,
		Values:   values,
	}
}

func (p *PlayerAction) Update() {
	if p.Complete {
		return
	}
	p.Selector.Update()
	if p.Selector.IsDone() {
		result := p.Selector.Finish()
		p.Action(result, p.Values)
		p.Complete = true
	} else if p.Selector.IsCancelled() {
		p.Cancel = true
	}
}