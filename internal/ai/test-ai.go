package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/actions"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

type RandomWalker struct {
	*AbstractAI
}

func NewRandomWalker(character *floor.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
		decision: -1,
		currValues: selector.ActionValues{
			Source: character,
		},
	}
	walker := &RandomWalker{
		newAI,
	}
	newAI.AI = walker
	return newAI
}

func (ai *RandomWalker) Decide() {
	orig := ai.Character.Coords
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
			a := world.NextHexLine(orig, n)
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
		world.NextHexLine(orig, mov),
	}
	atkCheck.Orig = mov

	ai.currValues.Damage = 1
	ai.currValues.Move = 1
	ai.Actions = []*AIAction{
		{
			Path:        path,
			PathCheck:   movCheck,
			TargetArea:  nil,
			TargetCheck: floor.PathChecks{},
			Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
			IsMove:      true,
		},
		{
			Path:        atk,
			PathCheck:   floor.PathChecks{},
			TargetArea:  []world.Coords{world.Origin},
			TargetCheck: atkCheck,
			Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
		},
	}
}

func (ai *RandomWalker) TakeTurn() {
	for i, act := range ai.TempActions {
		switch i % 2 {
		case 0:
			AddToBot(actions.NewMoveSeriesAction(ai.currValues.Source, ai.currValues.Source, act.Area), act.Effect)
		case 1:
			AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
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

func NewFlyChaser(character *floor.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
		decision: -1,
		currValues: selector.ActionValues{
			Source: character,
		},
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
	targets := floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 6)
	if len(targets) > 0 {
		for i := 0; i < 3; i++ {
			choice := targets[rand.Intn(len(targets))]
			if path, d, legal := floor.CurrentFloor.FindPathWithinOne(orig, choice, movCheck); legal {
				if d > 2 {
					tPath := floor.CurrentFloor.LongestLegalPath(path, 3, movCheck)
					if len(tPath) > 0 {
						ai.Actions = []*AIAction{
							{
								Path:        tPath,
								PathCheck:   movCheck,
								TargetArea:  nil,
								TargetCheck: floor.PathChecks{},
								Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
								IsMove:      true,
							},
						}
						return
					}
				} else {
					ai.currValues.Damage = 1
					ai.Actions = []*AIAction{
						{
							Path:        path,
							PathCheck:   movCheck,
							TargetArea:  nil,
							TargetCheck: floor.PathChecks{},
							Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
							IsMove:      true,
						},
						{
							Path:        []world.Coords{path[len(path)-1], choice},
							PathCheck:   floor.NoCheck,
							TargetArea:  []world.Coords{world.Origin},
							TargetCheck: atkCheck,
							Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
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
		within, _ := world.Remove(orig, floor.CurrentFloor.AllWithin(orig, 2, movCheck))
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
							Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
							IsMove:      true,
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
			AddToBot(actions.NewMoveSeriesAction(ai.Character, ai.Character, act.Area), act.Effect)
		case 1:
			AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
		}
	}
}

// If the player is further than 10 tiles, patrols between two points
// If the player is between 8-10 tiles, gets w/in 8
// If the player is between 4-7 tiles, strafe and attack from range
// If the player is between 2-3 tiles, retreat
// Otherwise, plink the player
type Skirmisher struct {
	*AbstractAI
	patrol     []world.Coords
	patrolling int
}

func NewSkirmisher(character *floor.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
		decision: -1,
		currValues: selector.ActionValues{
			Source: character,
		},
	}
	skirm := &Skirmisher{
		newAI,
		[]world.Coords{},
		0,
	}
	newAI.AI = skirm
	return newAI
}

func (ai *Skirmisher) Decide() {
	orig := ai.Character.Coords
	movCheck := floor.PathChecks{
		NotFilled:     true,
		Unoccupied:    true,
		NonEmpty:      true,
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
	ai.currValues.Damage = 1
	if len(ai.patrol) == 0 {
		patrolCand, _ := world.Remove(orig, floor.CurrentFloor.AllWithin(orig, 6, movCheck))
		ordered := world.ReverseList(world.OrderByDistSimple(orig, patrolCand))
		choice := 0
		if len(ordered) > 8 {
			choice = rand.Intn(len(ordered) / 8)
		}
		chosen := ordered[choice]
		ai.patrol = []world.Coords{
			orig,
			chosen,
		}
		ai.patrolling = 1
	}
	targets := floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 1)
	if len(targets) > 0 {
		choice := 0
		if len(targets) > 1 {
			choice = rand.Intn(len(targets))
		}
		ai.Actions = []*AIAction{
			{
				Path:        []world.Coords{orig, targets[choice]},
				PathCheck:   floor.NoCheck,
				TargetArea:  []world.Coords{world.Origin},
				TargetCheck: atkCheck,
				Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
			},
		}
		ai.decision++
		return
	}
	targets = floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 3)
	if len(targets) > 0 {
		dist := world.OrderByDistSimple(orig, targets)
		path, _, legal := floor.CurrentFloor.FindPathAwayFrom(orig, dist[0], 3, movCheck)
		if legal {
			ai.Actions = []*AIAction{
				{
					Path:        path,
					PathCheck:   movCheck,
					TargetArea:  nil,
					TargetCheck: floor.PathChecks{},
					Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
					IsMove:      true,
				},
			}
			ai.decision = 0
			return
		}
	}
	targets = floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 7)
	if len(targets) > 0 {
		choice := 0
		if len(targets) > 1 {
			choice = rand.Intn(len(targets))
		}
		path, _, legal := floor.CurrentFloor.FindPathPerpendicularTo(orig, targets[choice], 3, 7, movCheck, atkCheck)
		if legal {
			ai.Actions = []*AIAction{
				{
					Path:        path,
					PathCheck:   movCheck,
					TargetArea:  nil,
					TargetCheck: floor.PathChecks{},
					Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
					IsMove:      true,
				},
			}
			end := path[len(path)-1]
			atkPath := floor.CurrentFloor.LongestLegalPath(floor.CurrentFloor.Line(end, targets[choice], 7), 7, atkCheck)
			if world.CoordsIn(targets[choice], atkPath) {
				ai.Actions = append(ai.Actions, &AIAction{
					Path:        []world.Coords{end, targets[choice]},
					PathCheck:   atkCheck,
					TargetArea:  []world.Coords{world.Origin},
					TargetCheck: atkCheck,
					Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
				})
			}
		}
		ai.decision = 2
		return
	}
	targets = floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 10)
	if len(targets) > 0 {
		choice := targets[rand.Intn(len(targets))]
		if path, _, legal := floor.CurrentFloor.FindPathWithinOne(orig, choice, movCheck); legal {
			tPath := floor.CurrentFloor.LongestLegalPath(path, 3, movCheck)
			if len(tPath) > 0 {
				ai.Actions = []*AIAction{
					{
						Path:        tPath,
						PathCheck:   movCheck,
						TargetArea:  nil,
						TargetCheck: floor.PathChecks{},
						Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
						IsMove:      true,
					},
				}
				ai.decision = 0
				return
			}
		}
	}
	if len(ai.patrol) > 1 {
		if path, _, legal := floor.CurrentFloor.FindPath(orig, ai.patrol[ai.patrolling], movCheck); legal {
			tPath := floor.CurrentFloor.LongestLegalPath(path, 3, movCheck)
			if len(tPath) > 0 {
				ai.Actions = []*AIAction{
					{
						Path:        tPath,
						PathCheck:   movCheck,
						TargetArea:  nil,
						TargetCheck: floor.PathChecks{},
						Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
						IsMove:      true,
					},
				}
				ai.decision = 0
				if tPath[len(tPath)-1] == ai.patrol[ai.patrolling] {
					ai.patrolling = (ai.patrolling + 1) % 2
				}
			}
		} else {
			ai.patrol = []world.Coords{}
		}
	}
}

func (ai *Skirmisher) TakeTurn() {
	if len(ai.TempActions) > 0 {
		switch ai.decision {
		case 0:
			act := ai.TempActions[0]
			AddToBot(actions.NewMoveSeriesAction(ai.Character, ai.Character, act.Area), act.Effect)
		case 1:
			act := ai.TempActions[0]
			AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
		case 2:
			for i, act := range ai.TempActions {
				if i == 0 {
					AddToBot(actions.NewMoveSeriesAction(ai.Character, ai.Character, act.Area), act.Effect)
				} else {
					AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
				}
			}
		}
	}
}

// If the player is further than 10 tiles, do nothing, maybe random walk?
// If the player is between 4-10 tiles, do bombard strike or scatter bomb
// If the player is between 1-3 tiles, blast in a cone
// bombard strike and scatter bomb can be done twice, then rest
// always rest after blast in a cone
type Grenadier struct {
	*AbstractAI
	atkCnt   int
}

func NewGrenadier(character *floor.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
		decision: -1,
		currValues: selector.ActionValues{
			Source: character,
		},
	}
	gren := &Grenadier{
		newAI,
		0,
	}
	newAI.AI = gren
	return newAI
}

func (ai *Grenadier) Decide() {
	if ai.atkCnt > 1 {
		ai.atkCnt = 0
		return
	}
	orig := ai.Character.Coords
	movCheck := floor.PathChecks{
		NotFilled:     true,
		Unoccupied:    true,
		NonEmpty:      true,
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
	ai.currValues.Damage = 0
	targets := floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 3)
	if len(targets) > 0 {
		ai.decision = 2
		choice := 0
		if len(targets) > 1 {
			choice = rand.Intn(len(targets))
		}
		ai.Actions = []*AIAction{}
		area := floor.CurrentFloor.AllInSextant(orig, targets[choice], 3, atkCheck)
		if world.CoordsIn(targets[choice], area) {
			ai.currValues.Damage = 3
			ai.Actions = append(ai.Actions, &AIAction{
				Path:        []world.Coords{orig},
				PathCheck:   atkCheck,
				TargetArea:  area,
				TargetCheck: atkCheck,
				Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
			})
			ai.atkCnt += 2
		}
		return
	}
	targets = floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 10)
	if len(targets) > 0 {
		dec := 0
		if ai.decision == 0 {
			dec = 1
		} else if ai.decision != 1 {
			dec = rand.Intn(2)
		}
		ai.decision = dec
		choice := 0
		if len(targets) > 1 {
			choice = rand.Intn(len(targets))
		}
		ai.Actions = []*AIAction{}
		atkPath := floor.CurrentFloor.LongestLegalPath(floor.CurrentFloor.Line(orig, targets[choice], 10), 10, atkCheck)
		if world.CoordsIn(targets[choice], atkPath) {
			switch ai.decision {
			case 0:
				// bombard
				var best []world.Coords
				pts := -3
				n := world.RandomizeList(targets[choice].Neighbors(floor.CurrentFloor.Dimensions()))
				if len(n) > 1 {
					for _, c := range n {
						tpts := 0
						if cha := floor.CurrentFloor.GetOccupant(c); cha != nil {
							if cha.Diplomacy == floor.Ally {
								tpts += 1
							} else {
								tpts -= 1
							}
						}
						next := world.NextHexRot(c, targets[choice], true)
						if cha := floor.CurrentFloor.GetOccupant(next); cha != nil {
							if cha.Diplomacy == floor.Ally {
								tpts += 1
							} else {
								tpts -= 1
							}
						}
						if tpts > pts || len(best) == 0 {
							best = []world.Coords{c, next}
							pts = tpts
						}
					}
					ai.currValues.Damage = 2
					ai.Actions = append(ai.Actions, &AIAction{
						Path:        []world.Coords{orig, targets[choice]},
						PathCheck:   atkCheck,
						TargetArea:  append([]world.Coords{targets[choice]}, best...),
						TargetCheck: atkCheck,
						Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
					})
					ai.atkCnt += 1
				}
			case 1:
				// scatter shot
				atkCheck.Orig = targets[choice]
				s, _ := world.Remove(targets[choice], floor.CurrentFloor.AllWithin(targets[choice], 2, atkCheck))
				n := world.RandomizeList(s)
				count := len(n) / 3
				if count > 3 {
					count = 3
				}
				var hits []world.Coords
				for i := 0; i < count; i++ {
					hits = append(hits, n[rand.Intn(len(n)-1)])
				}
				ai.currValues.Damage = 2
				ai.Actions = append(ai.Actions, &AIAction{
					Path:        []world.Coords{orig, targets[choice]},
					PathCheck:   atkCheck,
					TargetArea:  append([]world.Coords{targets[choice]}, hits...),
					TargetCheck: atkCheck,
					Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
				})
				ai.atkCnt += 1
			}
		}
		path, _, legal := floor.CurrentFloor.FindPathPerpendicularTo(orig, targets[choice], 3, 10, movCheck, atkCheck)
		if legal {
			ai.Actions = append(ai.Actions, &AIAction{
				Path:        path,
				PathCheck:   movCheck,
				TargetArea:  nil,
				TargetCheck: floor.PathChecks{},
				Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
				IsMove:      true,
			})
		}
	}
}

func (ai *Grenadier) TakeTurn() {
	for i, act := range ai.TempActions {
		if i == 0 {
			AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
		} else {
			AddToBot(actions.NewMoveSeriesAction(ai.Character, ai.Character, act.Area), act.Effect)
		}
	}
}

// If the player is further than 10 tiles, do nothing, maybe random walk?
// If the player is between 4-10 tiles, chase player (slowly)
// If the player is between 2-3 tiles, approach and whack is small area
// If the player is 1 tile away, whack in 5 wide ring w/small range boost in front
// Rest after 3-4 attacks
type Bruiser struct {
	*AbstractAI
	atkCnt   int
}

func NewBruiser(character *floor.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
		decision: -1,
		currValues: selector.ActionValues{
			Source: character,
		},
	}
	bruiser := &Bruiser{
		newAI,
		0,
	}
	newAI.AI = bruiser
	return newAI
}

func (ai *Bruiser) Decide() {
	b := rand.Intn(2)
	if ai.atkCnt > 3+b {
		ai.atkCnt = 0
		return
	}
	orig := ai.Character.Coords
	movCheck := floor.PathChecks{
		NotFilled:     true,
		Unoccupied:    true,
		NonEmpty:      true,
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
	targets := floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 1)
	if len(targets) > 0 {
		ai.decision = 1
		choice := 0
		if len(targets) > 1 {
			choice = rand.Intn(len(targets))
		}
		ai.Actions = []*AIAction{}
		area := []world.Coords{orig}
		n1, _ := world.Remove(world.NextHexLine(targets[choice], orig), world.OrderByDist(targets[choice], orig.Neighbors(floor.CurrentFloor.Dimensions())))
		n2 := targets[choice].Neighbors(floor.CurrentFloor.Dimensions())
		area = world.Combine(area, world.Combine(n1, n2))
		if world.CoordsIn(targets[choice], area) {
			ai.currValues.Damage = 4
			ai.Actions = append(ai.Actions, &AIAction{
				Path:        []world.Coords{orig},
				PathCheck:   atkCheck,
				TargetArea:  area,
				TargetCheck: atkCheck,
				Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
			})
			ai.atkCnt++
		}
		return
	}
	targets = floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 3)
	if len(targets) > 0 {
		ai.decision = 2
		choice := 0
		if len(targets) > 1 {
			choice = rand.Intn(len(targets))
		}
		if path, d, legal := floor.CurrentFloor.FindPathWithinOne(orig, targets[choice], movCheck); legal {
			if d <= 3 {
				ai.Actions = []*AIAction{
					{
						Path:        path,
						PathCheck:   movCheck,
						TargetArea:  nil,
						TargetCheck: floor.PathChecks{},
						Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
						IsMove:      true,
					},
				}
				tOrig := path[len(path)-1]
				next := world.NextHexLine(tOrig, targets[choice])
				ai.currValues.Damage = 2
				ai.Actions = append(ai.Actions, &AIAction{
					Path:        []world.Coords{tOrig},
					PathCheck:   floor.NoCheck,
					TargetArea:  []world.Coords{tOrig, targets[choice], next},
					TargetCheck: atkCheck,
					Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
				})
				ai.atkCnt++
				return
			}
		}
	}
	targets = floor.CharacterManager.GetDiplomatic(floor.Ally, orig, 10)
	if len(targets) > 0 {
		ai.decision = 0
		choice := 0
		if len(targets) > 1 {
			choice = rand.Intn(len(targets))
		}
		if path, _, legal := floor.CurrentFloor.FindPathWithinOne(orig, targets[choice], movCheck); legal {
			tPath := floor.CurrentFloor.LongestLegalPath(path, 2, movCheck)
			ai.Actions = append(ai.Actions, &AIAction{
				Path:        tPath,
				PathCheck:   movCheck,
				TargetArea:  nil,
				TargetCheck: floor.PathChecks{},
				Effect:      selector.NewSelectionEffect(&selector.MoveEffect{}, ai.currValues),
				IsMove:      true,
			})
		}
	}
}

func (ai *Bruiser) TakeTurn() {
	if len(ai.TempActions) > 0 {
		if ai.decision == 0 {
			act := ai.TempActions[0]
			AddToBot(actions.NewMoveSeriesAction(ai.Character, ai.Character, act.Area), act.Effect)
		} else if ai.decision == 1 {
			act := ai.TempActions[0]
			AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
		} else {
			for i, act := range ai.TempActions {
				if i == 0 {
					AddToBot(actions.NewMoveSeriesAction(ai.Character, ai.Character, act.Area), act.Effect)
				} else {
					AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
				}
			}
		}
	}
}

// Every other turn, attack each tile around the enemy
type Stationary struct {
	*AbstractAI
	decision int
}

func NewStationary(character *floor.Character) *AbstractAI {
	newAI := &AbstractAI{
		Character: character,
		currValues: selector.ActionValues{
			Source: character,
		},
	}
	stat := &Stationary{
		newAI,
		0,
	}
	newAI.AI = stat
	return newAI
}

func (ai *Stationary) Decide() {
	if ai.decision == 0 {
		orig := ai.Character.Coords
		atkCheck := floor.PathChecks{
			NotFilled:     true,
			Unoccupied:    false,
			NonEmpty:      false,
			EndUnoccupied: false,
			Orig:          orig,
		}
		ai.currValues.Damage = 2
		area := append([]world.Coords{orig}, orig.Neighbors(floor.CurrentFloor.Dimensions())...)
		ai.Actions = []*AIAction{
			{
				Path:        []world.Coords{orig},
				PathCheck:   floor.NoCheck,
				TargetArea:  area,
				TargetCheck: atkCheck,
				Effect:      selector.NewSelectionEffect(&selector.AttackEffect{}, ai.currValues),
			},
		}
		ai.decision = 1
	} else {
		ai.decision = 0
	}
}

func (ai *Stationary) TakeTurn() {
	if len(ai.TempActions) > 0 {
		act := ai.TempActions[0]
		AddToBot(actions.NewDamageHexAction(act.Area, ai.currValues), act.Effect)
	}
}
