package objects

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type Targetable interface {
	Damage(dmg int)
}

type Moveable interface {
	GetCoords() world.Coords
	SetCoords(world.Coords)
	GetTransform() *animation.Transform
	SetTransformEffect(*animation.TransformEffect)
	GetPos() pixel.Vec
	SetPos(pixel.Vec)
	IsMoving() bool
}