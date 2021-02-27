package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

func RandomWalkerDecision(character *characters.Character, previous []int) ([]*AIAction, int) {
	orig := character.Coords
	movCheck := floor.PathChecks{
		NotFilled:  true,
		Unoccupied: true,
		NonEmpty:   true,
		Orig:       orig,
	}
	atkCheck := floor.PathChecks{
		NotFilled:  true,
		Unoccupied: false,
		NonEmpty:   false,
		Orig:       orig,
	}
	var choice int
	neighbors := orig.Neighbors(floor.CurrentFloor.Dimensions())
	indexList := make([]int, 0)
	for i, n := range neighbors {
		a := world.NextHex(orig, n)
		atkCheck.Orig = a
		if h := floor.CurrentFloor.IsLegal(a, atkCheck); h != nil {
			indexList = append(indexList, i)
		}
	}
	if len(indexList) > 0 {
		choice = indexList[rand.Intn(len(indexList))]
	} else {
		choice = rand.Intn(len(neighbors))
	}
	mov := neighbors[choice]
	path := []world.Coords{
		orig,
		mov,
	}
	atk := []world.Coords{
		mov,
		world.NextHex(orig, mov),
	}
	atkCheck.Orig = mov

	return []*AIAction{
		{
			Path:        path,
			PathCheck:   movCheck,
			TargetArea:  nil,
			TargetCheck: floor.PathChecks{},
			Values:      actions.ActionValues{
				Move: 1,
			},
		},
		{
			Path:        atk,
			PathCheck:   floor.PathChecks{},
			TargetArea:  &selectors.TargetArea{SetArea: selectors.SingleTile},
			TargetCheck: atkCheck,
			Values:      actions.ActionValues{
				Damage: 1,
			},
		},
	}, 0
}

func RandomWalkerAct(acts []*TempAIAction) {
	for i, act := range acts {
		switch i % 2 {
		case 0:
			actions.AddToBot(actions.NewMoveSeriesAction(act.Values.Source, act.Values.Source, act.Area))
		case 1:
			actions.AddToBot(actions.NewDamageHexAction(act.Values.Source, act.Area, act.Values.Damage))
		}
	}
}