package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type AI struct {
	MakeDecision func()
	Actions      []AIAction
	TempActions  []AIAction
	Character    *characters.Character
}

type AIAction struct {
	Path        []world.Coords
	PathCheck   *floor.PathChecks
	TargetArea  *selectors.TargetArea
	TargetCheck *floor.PathChecks
	Values      actions.ActionValues
}

func NewAI(makeDecision, update, takeTurn func()) *AI {
	return &AI{
		MakeDecision: makeDecision,
	}
}

func (ai *AI) Update() {
	//for _, act := range ai.Actions {
		// check the path

		// update the tempPath with the result of the check

		// update the targets at the result of the check
	//}
}