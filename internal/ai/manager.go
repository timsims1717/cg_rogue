package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/internal/ui"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type AI struct {
	MakeDecision func() []AIAction
	Actions      []AIAction
	TempActions  []AIAction
	Character    *characters.Character
	TempCoords   world.Coords
}

type AIAction struct {
	Path        []world.Coords
	PathCheck   floor.PathChecks
	TargetArea  *selectors.TargetArea
	TargetCheck floor.PathChecks
	Values      actions.ActionValues
}

func NewAI(makeDecision func() []AIAction, character *characters.Character) *AI {
	return &AI{
		MakeDecision: makeDecision,
		Character:    character,
	}
}

func (ai *AI) StartTurn() {
	ai.Actions = ai.MakeDecision()
}

func (ai *AI) Update() {
	ai.TempActions = make([]AIAction, len(ai.Actions))
	ai.TempCoords = ai.Character.GetCoords()
	for i, act := range ai.Actions {
		// check the path
		if act.Path != nil {
			tCheck := act.PathCheck
			tCheck.Orig = ai.TempCoords
			tPath := floor.CurrentFloor.LegalPath(tCheck.Orig.PathFrom(act.Path), tCheck)
			// check the target area

			// update the temp actions with the results of the check
			ai.TempActions[i] = AIAction{
				Path:        tPath,
				PathCheck:   act.PathCheck,
				TargetArea:  nil,
				TargetCheck: act.TargetCheck,
				Values:      act.Values,
			}
			if act.Values.Move > 0 && len(tPath) > 0 {
				ai.TempCoords = tPath[len(tPath)-1]
			}
		}
	}
	for _, act := range ai.TempActions {
		if act.Path != nil {
			for _, c := range act.Path {
				ui.AddSelectUI(ui.Move, c.X, c.Y)
			}
		}
	}
}