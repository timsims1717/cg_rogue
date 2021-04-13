package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type AI interface {
	Decide()
	TakeTurn()
}

type AbstractAI struct {
	AI            AI
	Actions       []*AIAction
	TempActions   []*TempAIAction
	Character     *floor.Character
	TempCoords    world.Coords
	PrevDecicions []int
}

type AIAction struct {
	Path        []world.Coords
	PathCheck   floor.PathChecks
	TargetArea  []world.Coords
	TargetCheck floor.PathChecks
	Values      selector.ActionValues
}

type TempAIAction struct {
	Area   []world.Coords
	Values selector.ActionValues
}

func (ai *AbstractAI) IsAlive() bool {
	return ai.Character.Health.Alive
}

func (ai *AbstractAI) Update() {
	if ai.Character.Health.Alive {
		for _, act := range ai.Actions {
			act.Values.Source = ai.Character
		}
		ai.TempActions = make([]*TempAIAction, len(ai.Actions))
		ai.TempCoords = ai.Character.GetCoords()
		for i, act := range ai.Actions {
			// check the path
			var tArea []world.Coords
			var tPath []world.Coords
			if act.Path != nil {
				tCheck := act.PathCheck
				tCheck.Orig = ai.TempCoords
				tPath = floor.CurrentFloor.LongestLegalPath(tCheck.Orig.PathFrom(act.Path), 0, tCheck)
				tArea = tPath
			}
			// check the target area
			if act.TargetArea != nil && len(tPath) > 0 {
				tCheck := act.TargetCheck
				target := tPath[len(tPath)-1]
				tCheck.Orig = target
				tArea, _ = world.Remove(ai.TempCoords, floor.CurrentFloor.IsSetLegal(tCheck.Orig.PathFrom(act.TargetArea), tCheck))
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
						if len(act.Area) > 1 {
							selector.AddSelectUI(selector.Move, c.X, c.Y)
						}
					} else {
						selector.AddSelectUI(selector.Attack, c.X, c.Y)
					}
				}
			}
		}
	}
}
