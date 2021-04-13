package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/action"
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
	Stamina       int
	decision      int
	currValues    selector.ActionValues
}

type AIAction struct {
	Effect      *selector.AbstractSelectionEffect
	Path        []world.Coords
	PathCheck   floor.PathChecks
	TargetArea  []world.Coords
	TargetCheck floor.PathChecks
	IsMove      bool
}

type TempAIAction struct {
	Area   []world.Coords
	Effect *selector.AbstractSelectionEffect
}

func (ai *AbstractAI) IsAlive() bool {
	return ai.Character.Health.Alive
}

func (ai *AbstractAI) Update() {
	if ai.Character.Health.Alive {
		ai.currValues.Source = ai.Character
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
				Effect: act.Effect,
			}
			if act.IsMove && len(tPath) > 0 {
				ai.TempCoords = tPath[len(tPath)-1]
			}
		}
		for _, act := range ai.TempActions {
			if act.Area != nil && len(act.Area) > 0 && act.Effect != nil {
				act.Effect.SetArea(act.Area)
				selector.AddSelectionEffect(act.Effect)
			}
		}
	}
}

func AddToTop(a action.Action, effects... *selector.AbstractSelectionEffect) {
	action.ActionManager.AddToTop(a, effects)
}

func AddToBot(a action.Action, effects... *selector.AbstractSelectionEffect) {
	action.ActionManager.AddToBot(a, effects)
}

