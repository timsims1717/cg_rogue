package floor

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
	"image/color"
)

type Character struct {
	Coords      world.Coords
	Claim       world.Coords
	Transform   *animation.Transform
	transEffect *animation.TransformEffect
	Spr         *pixel.Sprite
	OnMap       bool
	Floor       *Floor

	Health    Health
	Defense   Defense
	Diplomacy Diplomacy

	id    uuid.UUID
	index int

	mask   color.RGBA
	effect *animation.ColorEffect
}

// NewCharacter will override any character on the world.Coords supplied in the arguments.
func NewCharacter(sprite *pixel.Sprite, floor *Floor, coords world.Coords, diplomacy Diplomacy, maxHP int) *Character {
	transform := animation.NewTransform(true)
	transform.Pos = world.MapToWorld(coords)
	transform.Offset = cfg.OffsetVector
	c := &Character{
		Coords:    coords,
		Claim:     world.Nil,
		Transform: transform,
		Spr:       sprite,
		Floor:     floor,
		Health: Health{
			CurrHP:  maxHP,
			MaxHP:   maxHP,
			Alive:   true,
			Display: true,
			imd:     imdraw.New(nil),
		},
		Defense: Defense{
			MaxDef:  6,
			Alive:   true,
			Display: true,
			imd:     imdraw.New(nil),
		},
		Diplomacy: diplomacy,
		id:        uuid.NewV4(),
		mask:      colornames.White,
	}
	if c.Floor != nil {
		c.Floor.PutOccupant(c, coords)
	}
	return c
}

func (c *Character) Update() {
	c.Transform.Update(pixel.Rect{})
	if c.transEffect != nil {
		c.transEffect.Update()
		if c.transEffect.IsDone() {
			c.transEffect = nil
		}
	}
	if c.effect != nil {
		c.effect.Update()
		if c.effect.IsDone() {
			c.effect = nil
		}
	}
	c.Health.pos = pixel.V(c.Transform.Pos.X, c.Transform.Pos.Y+29.)
	c.Defense.pos = pixel.V(c.Transform.Pos.X, c.Transform.Pos.Y+32.)
	c.Health.Update()
	c.Defense.Update()
}

func (c *Character) Draw(win *pixelgl.Window) {
	c.Spr.DrawColorMask(win, c.Transform.Mat, c.mask)
	c.Health.Draw(win)
	c.Defense.Draw(win)
}

func (c *Character) StartTurn() {
	c.Defense.RemoveAll()
}

func (c *Character) Heal(amt int) {
	thisAmt := amt
	if thisAmt < 0 {
		thisAmt = 0
	}
	if thisAmt > 0 {
		c.Health.CurrHP += thisAmt
		if c.Health.CurrHP > c.Health.MaxHP {
			c.Health.CurrHP = c.Health.MaxHP
		}
		c.effect = animation.FadeFrom(c, colornames.Lightgoldenrodyellow, 0.5)
	}
}

func (c *Character) Damage(dmg int) {
	defDmg := c.Defense.Damage(dmg)
	c.Health.Damage(defDmg)
	if c.Health.CurrHP < 1 {
		col := colornames.Black
		col.A = 0
		c.effect = animation.FadeOut(c, 0.5)
		c.Health.Alive = false
		c.Health.CurrHP = 0
		c.Floor.Store(c)
	} else {
		c.effect = animation.FadeFrom(c, colornames.Red, 0.5)
	}
}

func (c *Character) IsDestroyed() bool {
	return !c.Health.Alive
}

func (c *Character) ID() uuid.UUID {
	return c.id
}

func (c *Character) RemoveClaim() {
	if c != nil && c.Floor != nil {
		if c.Floor.Get(c.Claim).IsClaimed() {
			c.Floor.Get(c.Claim).RemoveClaim()
			c.Claim = world.Nil
		}
	}
}

func (c *Character) MakeClaim(a world.Coords) {
	if c.Floor != nil {
		if a == c.Claim {
			return
		}
		if c.Floor.Get(c.Claim).IsClaimed() {
			c.Floor.Get(c.Claim).RemoveClaim()
		}
		c.Floor.Get(a).Claim(c)
	}
}

func (c *Character) GetPos() pixel.Vec {
	return c.Transform.Pos
}

func (c *Character) SetPos(v pixel.Vec) {
	c.Transform.Pos = v
}

func (c *Character) GetCoords() world.Coords {
	return c.Coords
}

func (c *Character) GetTransform() *animation.Transform {
	return c.Transform
}

func (c *Character) SetTransformEffect(effect *animation.TransformEffect) {
	c.transEffect = effect
}

func (c *Character) IsMoving() bool {
	return c.transEffect != nil
}

func (c *Character) GetColor() color.RGBA {
	return c.mask
}

func (c *Character) SetColor(mask color.RGBA) {
	c.mask = mask
}
