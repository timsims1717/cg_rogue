package ai

var AIManager aiManager

type aiManager struct {
	set      []*AI
	takeTurn bool
	decide   bool
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
	for _, ai := range AIManager.set {
		if AIManager.decide {
			ai.Decide()
		}
		ai.Update()
		if AIManager.takeTurn {
			ai.TakeTurn()
		}
	}
	AIManager.takeTurn = false
	AIManager.decide = false
}