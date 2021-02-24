package actions

import (
	"github.com/phf/go-queue/queue"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"reflect"
)

// ActionManager is the singleton access point to the actionManager struct
var ActionManager *actionManager

// actionManager stores the action queue and the current Action
type actionManager struct {
	qu     queue.Queue
	act    Action
}

type ActionValues struct {
	Source  *characters.Character
	Damage  int
	Move    int
	Range   int
	Targets int
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
func IsActing() bool {
	return notNil(ActionManager.act) || ActionManager.qu.Len() > 0
}

// AddToTop adds an Action to the top of the queue
// to be processed after the current Action is complete.
func AddToTop(a Action) {
	if notNil(a) {
		ActionManager.qu.PushFront(a)
	}
}

// AddToTop adds an Action to the bottom of the queue
// to be processed after all other Action are complete.
func AddToBot(a Action) {
	if notNil(a) {
		ActionManager.qu.PushBack(a)
	}
}

// Update switches to the next Action in the queue if one
// is not already being processed, then processes the next
// Action.
func Update() {
	if !notNil(ActionManager.act) {
		if ActionManager.qu.Front() != nil {
			inter := ActionManager.qu.PopFront()
			if act, ok := inter.(Action); ok {
				ActionManager.act = act
			}
		}
	}
	if notNil(ActionManager.act) {
		ActionManager.act.Update()
		if ActionManager.act.IsDone() {
			ActionManager.act = nil
		}
	}
}

// notNil checks both if a is nil, and if the underlying
// Action is nil.
func notNil(a Action) bool {
	if a == nil {
		return false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(a).IsNil() {
			return false
		}
	}
	return true
}