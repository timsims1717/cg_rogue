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
	MakeDecision  func(*characters.Character, []int) ([]*AIAction, int)
	Act           func([]*TempAIAction)
	Actions       []*AIAction
	TempActions   []*TempAIAction
	Character     *characters.Character
	TempCoords    world.Coords
	PrevDecicions []int
}

type AIAction struct {
	Path        []world.Coords
	PathCheck   floor.PathChecks
	TargetArea  *selectors.TargetArea
	TargetCheck floor.PathChecks
	Values      actions.ActionValues
}

type TempAIAction struct {
	Area   []world.Coords
	Values actions.ActionValues
}

func NewAI(makeDecision func(*characters.Character, []int) ([]*AIAction, int), act func([]*TempAIAction), character *characters.Character) *AI {
	return &AI{
		MakeDecision: makeDecision,
		Act:          act,
		Character:    character,
	}
}

func (ai *AI) Decide() {
	if ai.Character.Alive {
		var next int
		ai.Actions, next = ai.MakeDecision(ai.Character, ai.PrevDecicions)
		ai.PrevDecicions = append(ai.PrevDecicions, next)
		for _, act := range ai.Actions {
			act.Values.Source = ai.Character
		}
	}
}

func (ai *AI) TakeTurn() {
	if ai.Character.Alive {
		ai.Act(ai.TempActions)
	}
}

func (ai *AI) Update() {
	if ai.Character.Alive {
		ai.TempActions = make([]*TempAIAction, len(ai.Actions))
		ai.TempCoords = ai.Character.GetCoords()
		for i, act := range ai.Actions {
			// check the path
			var tArea []world.Coords
			var tPath []world.Coords
			if act.Path != nil {
				tCheck := act.PathCheck
				tCheck.Orig = ai.TempCoords
				tPath = floor.CurrentFloor.LongestLegalPath(tCheck.Orig.PathFrom(act.Path), tCheck)
				tArea = tPath
			}
			// check the target area
			if act.TargetArea != nil && len(tPath) > 0 {
				tCheck := act.TargetCheck
				target := tPath[len(tPath)-1]
				tCheck.Orig = target
				tArea = act.TargetArea.SetArea(0, 0, target, tCheck)
			}
			// update the temp actions with the results of the check
			ai.TempActions[i] = &TempAIAction{
				Area:   tArea,
				Values: act.Values,
			}
			if act.Values.Move > 0 && len(tPath) > 0 {
				ai.TempCoords = tPath[len(tPath)-1]
			}
		}
		for _, act := range ai.TempActions {
			if act.Area != nil && len(act.Area) > 0 {
				for _, c := range act.Area {
					if act.Values.Move > 0 {
						ui.AddSelectUI(ui.Move, c.X, c.Y)
					} else {
						ui.AddSelectUI(ui.Attack, c.X, c.Y)
					}
				}
			}
		}
	}
}