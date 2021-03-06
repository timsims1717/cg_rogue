package ai

import "github.com/timsims1717/cg_rogue_go/internal/actions"

var AIManager aiManager

type aiManager struct {
	set       []*AI
	takeTurn  bool
	turnIndex int
	decide    bool
}

func (m *aiManager) StartAITurn() {
	m.takeTurn = true
}

func (m *aiManager) EndAITurn() {
	m.decide = true
}

func (m *aiManager) AddAI(ai *AI) {
	m.set = append(m.set, ai)
}

func Update() {
	if AIManager.takeTurn && !actions.IsActing() {
		AIManager.turnIndex += 1
	}
	if AIManager.takeTurn && !actions.IsActing() && AIManager.turnIndex == len(AIManager.set) {
		AIManager.takeTurn = false
		AIManager.turnIndex = -1
	}
	for i, ai := range AIManager.set {
		if AIManager.decide {
			ai.Decide()
		}
		ai.Update()
		if AIManager.takeTurn {
			if AIManager.turnIndex == i && !actions.IsActing() {
				ai.TakeTurn()
			}
		}
	}
	AIManager.decide = false
}