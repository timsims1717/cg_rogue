package actions

import (
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

func SetAttackTransform(source objects.Moveable, target world.Coords) {
	p := source.GetPos()
	b := world.MapToWorld(target)
	t := p.To(b).Scaled(0.3)
	e := p.Add(t)
	transform := animation.TransformBuilder{
		Transform: source.GetTransform(),
		InterX:    gween.New(p.X, e.X, 0.12, ease.InBack),
		InterY:    gween.New(p.Y, e.Y, 0.12, ease.InBack),
	}
	source.SetTransformEffect(transform.Build())
}

func SetResetTransform(source objects.Moveable) {
	o := world.MapToWorld(source.GetCoords())
	v := source.GetPos()
	transform := animation.TransformBuilder{
		Transform: source.GetTransform(),
		InterX:    gween.New(v.X, o.X, 0.2, ease.InQuint),
		InterY:    gween.New(v.Y, o.Y, 0.2, ease.InQuint),
	}
	source.SetTransformEffect(transform.Build())
}