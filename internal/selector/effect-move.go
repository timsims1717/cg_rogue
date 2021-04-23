package selector

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math"
)

type MoveEffect struct {
	*AbstractSelectionEffect
	//imd *imdraw.IMDraw
	matrices [][3]pixel.Matrix
	End      bool
}

func (e *MoveEffect) Update() {
	e.matrices = [][3]pixel.Matrix{}
	prev := e.orig
	if len(e.area) > 0 && e.area[len(e.area)-1] != e.orig {
		for _, c := range e.area {
			if e.End || c != e.orig {
				hex := floor.CurrentFloor.Get(c)
				hex.AddEffect([]img.Sprite{{
					S: gfx.SelectionSprites["move_single"],
					M: img.IM,
				}}, 11)
			}
		}
		if e.End {
			c := e.area[len(e.area)-1]
			dist := world.DistanceWorld(world.MapToWorld(prev), world.MapToWorld(c)) - cfg.TileSize
			angle := world.AngleBetween(c, prev) - math.Pi*0.5
			midX := world.MapToWorld(prev).X + (world.MapToWorld(c).X-world.MapToWorld(prev).X)*0.5
			midY := world.MapToWorld(prev).Y + (world.MapToWorld(c).Y-world.MapToWorld(prev).Y)*0.5
			mat := [3]pixel.Matrix{
				pixel.IM.Rotated(pixel.ZV, angle).Moved(world.MapToWorld(prev)),
				pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., dist/cfg.TileSize)).Rotated(pixel.ZV, angle).Moved(pixel.V(midX, midY)),
				pixel.IM.Rotated(pixel.ZV, angle).Moved(world.MapToWorld(c)),
			}
			e.matrices = append(e.matrices, mat)
		} else {
			for _, c := range e.area {
				dist := world.DistanceWorld(world.MapToWorld(prev), world.MapToWorld(c)) - cfg.TileSize
				angle := world.AngleBetween(c, prev) - math.Pi*0.5
				midX := world.MapToWorld(prev).X + (world.MapToWorld(c).X-world.MapToWorld(prev).X)*0.5
				midY := world.MapToWorld(prev).Y + (world.MapToWorld(c).Y-world.MapToWorld(prev).Y)*0.5
				mat := [3]pixel.Matrix{
					pixel.IM.Rotated(pixel.ZV, angle).Moved(world.MapToWorld(prev)),
					pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., dist/cfg.TileSize)).Rotated(pixel.ZV, angle).Moved(pixel.V(midX, midY)),
					pixel.IM.Rotated(pixel.ZV, angle).Moved(world.MapToWorld(c)),
				}
				e.matrices = append(e.matrices, mat)
				prev = c
			}
		}
	}
}

func (e *MoveEffect) Draw(target pixel.Target) {
	if len(e.area) > 0 && e.area[len(e.area)-1] != e.orig {
		if e.End {
			gfx.SelectionSprites["move_u_d"].Draw(target, e.matrices[0][1])
			gfx.SelectionSprites["move_start_u"].Draw(target, e.matrices[0][0])
			gfx.SelectionSprites["move_end_alone"].Draw(target, e.matrices[0][2])
		} else {
			for i, c := range e.area {
				if c != e.orig {
					gfx.SelectionSprites["move_u_d"].Draw(target, e.matrices[i][1])
					gfx.SelectionSprites["move_start_u"].Draw(target, e.matrices[i][0])
					gfx.SelectionSprites["move_end_alone"].Draw(target, e.matrices[i][2])
				}
			}
		}
	}
	//e.imd.Draw(target)
}

func (e *MoveEffect) SetAbstract(effect *AbstractSelectionEffect) {
	e.AbstractSelectionEffect = effect
}
