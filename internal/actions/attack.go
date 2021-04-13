package actions

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

func SetAttackTransformSingle(source *floor.Character, target world.Coords) {
	p := source.GetPos()
	b := world.MapToWorld(target).Sub(p)
	r := util.Normalize(b).Scaled(0.4 * cfg.ScaledTileSize)
	e := p.Add(r)
	transform := animation.TransformBuilder{
		Transform: source.GetTransform(),
		InterX:    gween.New(p.X, e.X, 0.12, ease.InBack),
		InterY:    gween.New(p.Y, e.Y, 0.12, ease.InBack),
	}
	source.SetTransformEffect(transform.Build())
}

func SetAttackTransform(source *floor.Character, targets []world.Coords) {
	p := source.GetPos()
	xs := 0.
	ys := 0.
	for _, c := range targets {
		t := world.MapToWorld(c)
		xs += t.X
		ys += t.Y
	}
	l := float64(len(targets))
	a := pixel.V(xs/l, ys/l)
	if a.Eq(p) {
		// todo: add jump/slam animation?
		a.Y += 1.0
	}
	b := a.Sub(p)
	r := util.Normalize(b).Scaled(0.4 * cfg.ScaledTileSize)
	e := p.Add(r)
	transform := animation.TransformBuilder{
		Transform: source.GetTransform(),
		InterX:    gween.New(p.X, e.X, 0.12, ease.InBack),
		InterY:    gween.New(p.Y, e.Y, 0.12, ease.InBack),
	}
	source.SetTransformEffect(transform.Build())
}

func SetResetTransform(source *floor.Character) {
	o := world.MapToWorld(source.GetCoords())
	v := source.GetPos()
	transform := animation.TransformBuilder{
		Transform: source.GetTransform(),
		InterX:    gween.New(v.X, o.X, 0.2, ease.InQuint),
		InterY:    gween.New(v.Y, o.Y, 0.2, ease.InQuint),
	}
	source.SetTransformEffect(transform.Build())
}
