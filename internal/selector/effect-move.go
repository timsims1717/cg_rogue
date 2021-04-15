package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math"
)

type MoveEffect struct {
	*AbstractSelectionEffect
	//imd *imdraw.IMDraw
	matrices [][3]pixel.Matrix
}

func (e *MoveEffect) Update() {
	e.matrices = [][3]pixel.Matrix{}
	prev := e.orig
	if len(e.area) > 0 {
		for _, c := range e.area {
			dist := world.DistanceWorld(world.MapToWorld(prev), world.MapToWorld(c)) - cfg.TileSize
			angle := world.AngleBetween(c, prev) - math.Pi * 0.5
			midX := world.MapToWorld(prev).X + (world.MapToWorld(c).X - world.MapToWorld(prev).X) * 0.5
			midY := world.MapToWorld(prev).Y + (world.MapToWorld(c).Y - world.MapToWorld(prev).Y) * 0.5
			mat := [3]pixel.Matrix{
				pixel.IM.Rotated(pixel.ZV, angle).Moved(world.MapToWorld(prev)),
				pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., dist / cfg.TileSize)).Rotated(pixel.ZV, angle).Moved(pixel.V(midX, midY)),
				pixel.IM.Rotated(pixel.ZV, angle).Moved(world.MapToWorld(c)),
			}
			e.matrices = append(e.matrices, mat)
			prev = c
		}
	}
}

func (e *MoveEffect) Draw(target pixel.Target) {
	if len(e.area) > 0 {
		for i, c := range e.area {
			if c != e.orig {
				mat := pixel.IM.Moved(world.MapToWorld(c))
				SelectionSprites["move_single"].Draw(target, mat)

				SelectionSprites["move_u_d"].Draw(target, e.matrices[i][1])
				SelectionSprites["move_start_u"].Draw(target, e.matrices[i][0])
				SelectionSprites["move_end_alone"].Draw(target, e.matrices[i][2])
			}
		}
	}
	//e.imd.Draw(target)
}

func (e *MoveEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
