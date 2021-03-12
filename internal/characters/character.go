package characters

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
	"image/color"
)

type Character struct {
	Coords world.Coords
	Mat    pixel.Matrix
	Pos    pixel.Vec
	Spr    *pixel.Sprite

	Health Health

	Diplomacy Diplomacy

	id    uuid.UUID
	index int

	mask   color.RGBA
	effect *animation.ColorEffect
}

func NewCharacter(sprite *pixel.Sprite, coords world.Coords, diplomacy Diplomacy, maxHP int) *Character {
	if floor.CurrentFloor.IsOccupied(coords) {
		return nil
	}
	c := &Character{
		Coords: coords,
		Mat:    pixel.Matrix{},
		Pos:    world.MapToWorld(coords),
		Spr:    sprite,
		Health: Health{
			CurrHP:  maxHP,
			MaxHP:   maxHP,
			LastDmg: 0,
			Alive:   true,
			imd:     imdraw.New(nil),
			pos:     pixel.ZV,
		},
		Diplomacy: diplomacy,
		id:        uuid.NewV4(),
		mask:      colornames.White,
	}
	floor.CurrentFloor.PutOccupant(c, coords)
	return c
}

func (c *Character) Update() {
	c.Mat = pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(c.Pos).Moved(cfg.OffsetVector)
	if c.effect != nil {
		c.effect.Update()
		if c.effect.IsDone() {
			c.effect = nil
		}
	}
	c.Health.pos = pixel.V(c.Pos.X, c.Pos.Y + 50.)
	c.Health.Update()
}

func (c *Character) Draw(win *pixelgl.Window) {
	c.Spr.DrawColorMask(win, c.Mat, c.mask)
	c.Health.Draw(win)
}

func (c *Character) Damage(dmg int) {
	thisDmg := dmg
	if thisDmg < 0 {
		thisDmg = 0
	}
	if thisDmg > 0 {
		c.Health.LastDmg = util.Min(thisDmg, c.Health.CurrHP)
		c.Health.CurrHP -= c.Health.LastDmg
		if c.Health.CurrHP < 1 {
			col := colornames.Black
			col.A = 0
			c.effect = animation.FadeOut(c, 0.5)
			c.Health.Alive = false
			c.Health.CurrHP = 0
			floor.CurrentFloor.RemoveOccupant(c.Coords)
		} else {
			c.effect = animation.FadeFrom(c, colornames.Red, 0.5)
		}
	}
}

func (c *Character) IsDestroyed() bool {
	return !c.Health.Alive
}

func (c *Character) ID() uuid.UUID {
	return c.id
}

func (c *Character) GetXY() (float64, float64) {
	return c.Pos.X, c.Pos.Y
}

func (c *Character) SetXY(x, y float64) {
	c.Pos.X = x
	c.Pos.Y = y
}

func (c *Character) GetCoords() world.Coords {
	return c.Coords
}

func (c *Character) SetCoords(a world.Coords) {
	c.Coords = a
}

func (c *Character) GetColor() color.RGBA {
	return c.mask
}

func (c *Character) SetColor(mask color.RGBA) {
	c.mask = mask
}