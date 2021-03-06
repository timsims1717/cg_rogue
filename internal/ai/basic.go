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
		NotFilled:     true,
		Unoccupied:    true,
		NonEmpty:      true,
		EndUnoccupied: true,
		Orig:          orig,
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
		if m := floor.CurrentFloor.IsLegal(n, movCheck); m != nil {
			a := world.NextHex(orig, n)
			atkCheck.Orig = a
			if h := floor.CurrentFloor.IsLegal(a, atkCheck); h != nil {
				indexList = append(indexList, i)
			}
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
			Values:      selectors.ActionValues{
				Move: 1,
			},
		},
		{
			Path:        atk,
			PathCheck:   floor.PathChecks{},
			TargetArea:  []world.Coords{world.Origin},
			TargetCheck: atkCheck,
			Values:      selectors.ActionValues{
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

// If the player is further than 6 tiles, 50% chance to move 1-2 tiles
// Otherwise, chases the player:
//   If w/in 3, move 2 and attack
//   Otherwise move 2
// After 3 attacks, must rest
func FlyChaserDecision(character *characters.Character, previous []int) ([]*AIAction, int) {
	prev := 0
	if len(previous) > 0 {
		prev = previous[len(previous)-1]
	}
	if prev == 3 {
		return []*AIAction{}, 0
	}
	orig := character.Coords
	movCheck := floor.PathChecks{
		NotFilled:     true,
		Unoccupied:    false,
		NonEmpty:      false,
		EndUnoccupied: true,
		Orig:          orig,
	}
	atkCheck := floor.PathChecks{
		NotFilled:     true,
		Unoccupied:    false,
		NonEmpty:      false,
		EndUnoccupied: false,
		Orig:          orig,
	}
	targets := characters.CharacterManager.GetDiplomatic(characters.Ally, orig, 6)
	if len(targets) > 0 {
		for i := 0; i < 3; i++ {
			choice := targets[rand.Intn(len(targets))]
			if path, d, legal := floor.CurrentFloor.FindPathWithinOne(orig, choice, movCheck); legal {
				if d > 2 {
					 tPath := floor.CurrentFloor.LongestLegalPath(path[:3], movCheck)
					 if len(tPath) > 0 {
						 return []*AIAction{
							 {
								 Path:        tPath,
								 PathCheck:   movCheck,
								 TargetArea:  nil,
								 TargetCheck: floor.PathChecks{},
								 Values: selectors.ActionValues{
									 Move: len(tPath),
								 },
							 },
						 }, prev
					 }
				} else {
					return []*AIAction{
						{
							Path:        path,
							PathCheck:   movCheck,
							TargetArea:  nil,
							TargetCheck: floor.PathChecks{},
							Values: selectors.ActionValues{
								Move: d,
							},
						},
						{
							Path:        []world.Coords{path[len(path)-1], choice},
							PathCheck:   floor.NoCheck,
							TargetArea:  []world.Coords{world.Origin},
							TargetCheck: atkCheck,
							Values:      selectors.ActionValues{
								Damage: 1,
							},
						},
					}, prev + 1
				}
			}
		}
	}
	if rand.Intn(2) > 0 {
		// 50% chance random walk
		within := floor.CurrentFloor.AllWithinNoPath(orig, 2, movCheck)
		if len(within) > 0 {
			for i := 0; i < 3; i++ {
				choice := within[rand.Intn(len(within))]
				if path, d, legal := floor.CurrentFloor.FindPath(orig, choice, movCheck); legal && d <= 2 {
					return []*AIAction{
						{
							Path:        path,
							PathCheck:   movCheck,
							TargetArea:  nil,
							TargetCheck: floor.PathChecks{},
							Values: selectors.ActionValues{
								Move: d,
							},
						},
					}, prev
				}
			}
		}
	}
	return []*AIAction{}, prev
}

func FlyChaserAct(acts []*TempAIAction) {
	for i, act := range acts {
		switch i % 2 {
		case 0:
			actions.AddToBot(actions.NewMoveSeriesAction(act.Values.Source, act.Values.Source, act.Area))
		case 1:
			actions.AddToBot(actions.NewDamageHexAction(act.Values.Source, act.Area, act.Values.Damage))
		}
	}
}