package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math"
)

type AttackTargetEffect struct {
	*AbstractSelectionEffect
	Center         bool
	targetMatrices [3]pixel.Matrix
}

func (e *AttackTargetEffect) Update() {
	e.targetMatrices = [3]pixel.Matrix{}
	if len(e.area) > 0 {
		var v pixel.Vec
		if e.Center {
			centerX := 0.
			centerY := 0.
			for _, c := range e.area {
				v := world.MapToWorld(c)
				centerX += v.X
				centerY += v.Y
			}
			centerX = centerX / float64(len(e.area))
			centerY = centerY / float64(len(e.area))
			v = pixel.V(centerX, centerY)
		} else {
			v = world.MapToWorld(e.area[0])
		}
		dist := world.DistanceWorld(world.MapToWorld(e.orig), v) - 6.0
		s := v.Sub(world.MapToWorld(e.orig))
		angle := s.Angle() - math.Pi * 0.5
		midX := world.MapToWorld(e.orig).X + (v.X - world.MapToWorld(e.orig).X) * 0.5
		midY := world.MapToWorld(e.orig).Y + (v.Y - world.MapToWorld(e.orig).Y) * 0.5
		e.targetMatrices = [3]pixel.Matrix{
			pixel.IM.Rotated(pixel.ZV, angle).Moved(world.MapToWorld(e.orig)),
			pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., dist / cfg.TileSize)).Rotated(pixel.ZV, angle).Moved(pixel.V(midX, midY)),
			pixel.IM.Moved(v),
		}
		sprites := attackAreas(e.area)
		for i, c := range e.area {
			hex := floor.CurrentFloor.Get(c)
			hex.AddEffect(sprites[i], 10)
		}
	}
}

func (e *AttackTargetEffect) Draw(target pixel.Target) {
	if len(e.area) > 0 {
		gfx.SelectionSprites["attack_target_mid"].Draw(target, e.targetMatrices[1])
		gfx.SelectionSprites["attack_target_start"].Draw(target, e.targetMatrices[0])
		gfx.SelectionSprites["attack_target"].Draw(target, e.targetMatrices[2])
	}
}

func (e *AttackTargetEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
