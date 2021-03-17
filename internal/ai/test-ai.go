package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

type RandomWalker struct {
	*AbstractAI
}

func NewRandomWalker(character *characters.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
	}
	walker := &RandomWalker{
		newAI,
	}
	newAI.AI = walker
	return newAI
}

func (r *RandomWalker) Decide() {
	orig := r.Character.Coords
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

	r.Actions = []*AIAction{
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
	}
}

func (r *RandomWalker) TakeTurn() {
	for i, act := range r.TempActions {
		switch i % 2 {
		case 0:
			actions.AddToBot(actions.NewMoveSeriesAction(act.Values.Source, act.Values.Source, act.Area))
		case 1:
			actions.AddToBot(actions.NewDamageHexAction(act.Area, act.Values))
		}
	}
}

// If the player is further than 6 tiles, 50% chance to move 1-2 tiles
// Otherwise, chases the player:
//   If w/in 3, move 2 and attack
//   Otherwise move 2
// After 3 attacks, must rest
type FlyChaser struct {
	*AbstractAI
	atkCnt int
}

func NewFlyChaser(character *characters.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
	}
	flyChaser := &FlyChaser{
		newAI,
		0,
	}
	newAI.AI = flyChaser
	return newAI
}

func (ai *FlyChaser) Decide() {
	if ai.atkCnt >= 3 {
		ai.Actions = []*AIAction{}
		ai.atkCnt = 0
		return
	}
	orig := ai.Character.Coords
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
						ai.Actions = []*AIAction{
							{
								Path:        tPath,
								PathCheck:   movCheck,
								TargetArea:  nil,
								TargetCheck: floor.PathChecks{},
								Values: selectors.ActionValues{
									Move: len(tPath),
								},
							},
						}
						return
					}
				} else {
					ai.Actions = []*AIAction{
						{
							Path:        path,
							PathCheck:   movCheck,
							TargetArea:  nil,
							TargetCheck: floor.PathChecks{},
							Values: selectors.ActionValues{
								Move: 1,
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
					}
					ai.atkCnt++
					return
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
					ai.Actions = []*AIAction{
						{
							Path:        path,
							PathCheck:   movCheck,
							TargetArea:  nil,
							TargetCheck: floor.PathChecks{},
							Values: selectors.ActionValues{
								Move: d,
							},
						},
					}
				}
			}
		}
	}
}

func (ai *FlyChaser) TakeTurn() {
	for i, act := range ai.TempActions {
		switch i % 2 {
		case 0:
			actions.AddToBot(actions.NewMoveSeriesAction(act.Values.Source, act.Values.Source, act.Area))
		case 1:
			actions.AddToBot(actions.NewDamageHexAction(act.Area, act.Values))
		}
	}
}

// If the player is further than 10 tiles, patrols between two points
// If the player is between 8-10 tiles, gets w/in 8
// If the player is between 4-7 tiles, strafe and attack from range
// If the player is between 2-3 tiles, retreat
// Otherwise, plink the player