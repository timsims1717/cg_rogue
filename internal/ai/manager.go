package ai

import (
	"github.com/timsims1717/cg_rogue_go/internal/manager"
)

var AIManager aiManager

type aiManager struct {
	set       []*AbstractAI
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

func (m *aiManager) AddAI(ai *AbstractAI) {
	m.set = append(m.set, ai)
}

func (m *aiManager) Clear() {
	m.set = []*AbstractAI{}
}

func (m *aiManager) AIActing() bool {
	return m.takeTurn
}

func (m *aiManager) Update() {
	if m.takeTurn && !manager.ActionManager.IsActing() {
		if m.turnIndex > -1 && m.turnIndex < len(m.set) {
			m.set[m.turnIndex].Actions = []*AIAction{}
		}
		m.turnIndex += 1
	}
	if m.takeTurn && !manager.ActionManager.IsActing() && m.turnIndex == len(m.set) {
		m.takeTurn = false
		m.turnIndex = -1
	}
	for i, ai := range m.set {
		if ai.IsAlive() {
			if m.decide {
				ai.AI.Decide()
			}
			ai.Update()
			if m.takeTurn {
				if m.turnIndex == i && !manager.ActionManager.IsActing() {
					ai.AI.TakeTurn()
				}
			}
		}
	}
	m.decide = false
}