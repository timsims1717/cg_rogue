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

var Characters characters

type characters struct{
	Set []*Character
}

func (c *characters) Update() {
	for _, ch := range c.Set {
		ch.Update()
	}
}

func (c *characters) Draw(win *pixelgl.Window) {
	for _, ch := range c.Set {
		ch.Draw(win)
	}
}

type Character struct {
	Coords world.Coords
	Mat pixel.Matrix
	Pos pixel.Vec
	Spr *pixel.Sprite

	CurrHP  int
	MaxHP   int
	LastDmg int
	Alive   bool

	Diplomacy Diplomacy

	id    uuid.UUID
	index int
}

func NewCharacter(sprite *pixel.Sprite, coords world.Coords, diplomacy Diplomacy, maxHP int) *Character {
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
		Alive:   true,
		Diplomacy: diplomacy,
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
			c.Alive = false
			c.CurrHP = 0
			floor.CurrentFloor.RemoveOccupant(c.Coords)
		}
	}
}

func (c *Character) IsDestroyed() bool {
	return !c.Alive
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