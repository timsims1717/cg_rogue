package characters

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
)

type Character struct {
	Coords world.Coords
	Mat pixel.Matrix
	Pos pixel.Vec
	Spr *pixel.Sprite

	CurrHP  int
	MaxHP   int
	LastDmg int
	Dead    bool

	id      uuid.UUID
}

func NewCharacter(sprite *pixel.Sprite, coords world.Coords, maxHP int) *Character {
	if floor.CurrentFloor.IsOccupied(coords) {
		return nil
	}
	c := &Character{
		Coords:  coords,
		Mat:     pixel.Matrix{},
		Pos:     pixel.V(world.MapToWorldHex(coords.X,coords.Y)),
		Spr:     sprite,
		CurrHP:  maxHP,
		MaxHP:   maxHP,
		LastDmg: 0,
		Dead:    false,
		id:      uuid.NewV4(),
	}
	floor.CurrentFloor.PutOccupant(c, coords)
	return c
}

func (c *Character) Update() {
	if !c.IsDestroyed() {
		c.Mat = pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(c.Pos).Moved(cfg.OffsetVector)
	}
}

func (c *Character) Draw(win *pixelgl.Window) {
	if !c.IsDestroyed() {
		c.Spr.Draw(win, c.Mat)
	}
}

func (c *Character) Damage(dmg int) {
	thisDmg := dmg
	if thisDmg < 0 {
		thisDmg = 0
	}
	if thisDmg > 0 {
		c.LastDmg = util.Min(thisDmg, c.CurrHP)
		c.CurrHP -= c.LastDmg
		if c.CurrHP < 1 {
			c.Dead = true
			c.CurrHP = 0
			floor.CurrentFloor.RemoveOccupant(c.Coords)
		}
	}
}

func (c *Character) IsTargeted() {

}

func (c *Character) TryToOccupy() {

}

func (c *Character) IsDestroyed() bool {
	return c.Dead
}

func (c *Character) ID() uuid.UUID {
	return c.id
}