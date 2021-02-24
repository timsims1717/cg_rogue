package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

func SpinnerDecision() []*AIAction {
	orig := world.Coords{
		X: 8,
		Y: 4,
	}
	checks := floor.PathChecks{
		NotFilled:  true,
		Unoccupied: true,
		NonEmpty:   true,
		Orig:       orig,
	}
	path, _, _ := floor.CurrentFloor.FindPath(orig, world.Coords{
		X: 4,
		Y: 3,
	}, checks)
	return []*AIAction{
		{
			Path:        path,
			PathCheck:   checks,
			TargetArea:  nil,
			TargetCheck: floor.PathChecks{},
			Values:      actions.ActionValues{
				Move: 1,
			},
		},
	}
}

func SpinnerAct(acts []*AIAction) {
	for i, act := range acts {
		switch i {
		case 0:
			actions.AddToBot(actions.NewMoveSeriesAction(act.Values.Source, act.Values.Source, act.Path))
		}
	}
}