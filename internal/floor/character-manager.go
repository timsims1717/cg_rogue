package floor

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

var CharacterManager *characterManager

type characterManager struct {
	set []*Character
}

func init() {
	CharacterManager = new(characterManager)
	CharacterManager.set = make([]*Character, 0)
}

func (m *characterManager) Add(character *Character) {
	if character == nil {
		return
	}
	character.index = len(m.set)
	m.set = append(m.set, character)
}

func (m *characterManager) Remove(character *Character) {
	m.set = append(m.set[:character.index], m.set[character.index+1:]...)
	m.updateSet()
}

func (m *characterManager) Clear() {
	m.set = []*Character{}
}

func (m *characterManager) updateSet() {
	for i, ch := range m.set {
		ch.index = i
	}
}

func (m *characterManager) GetDiplomatic(d Diplomacy, orig world.Coords, r int) []world.Coords {
	var set []world.Coords
	for _, ch := range m.set {
		if ch.Diplomacy == d && world.DistanceSimple(orig, ch.Coords) <= r && !ch.IsDestroyed() {
			set = append(set, ch.Coords)
		}
	}
	return set
}

func Update() {
	for _, ch := range CharacterManager.set {
		ch.Update()
	}
}

func Draw(win *pixelgl.Window) {
	for _, ch := range CharacterManager.set {
		ch.Draw(win)
	}
}
