package manager

import (
	"github.com/phf/go-queue/queue"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
)

// ActionManager is the singleton access point to the actionManager struct
var ActionManager *actionManager

// actionManager stores the action queue and the current Action
type actionManager struct {
	qu  queue.Queue
	act Action
}

// An Action can be updated and checked for completion.
// Everything else is determined by the Action itself.
type Action interface {
	Update()
	IsDone() bool
}

// init creates the singleton ActionManager
func init() {
	ActionManager = new(actionManager)
	ActionManager.qu.Init()
}

// IsActing returns true if there is no current action and no actions in the queue
func (m *actionManager) IsActing() bool {
	return !util.IsNil(m.act) || m.qu.Len() > 0
}

// AddToTop adds an Action to the top of the queue
// to be processed after the current Action is complete.
func (m *actionManager) AddToTop(a Action) {
	if !util.IsNil(a) {
		m.qu.PushFront(a)
	}
}

// AddToTop adds an Action to the bottom of the queue
// to be processed after all other Action are complete.
func (m *actionManager) AddToBot(a Action) {
	if !util.IsNil(a) {
		m.qu.PushBack(a)
	}
}

// Update switches to the next Action in the queue if one
// is not already being processed, then processes the next
// Action.
func (m *actionManager) Update() {
	if util.IsNil(m.act) {
		if m.qu.Front() != nil {
			inter := m.qu.PopFront()
			if act, ok := inter.(Action); ok {
				m.act = act
			}
		}
	}
	if !util.IsNil(m.act) {
		m.act.Update()
		if m.act.IsDone() {
			m.act = nil
		}
	}
}