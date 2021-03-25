package player

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	text2 "github.com/timsims1717/cg_rogue_go/pkg/typeface"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"golang.org/x/image/colornames"
)

var CardBG *pixel.Sprite

func Initialize() {
	cardsheet, err := img.LoadSpriteSheet("assets/cards/testcard.json")
	if err != nil {
		panic(err)
	}
	CardBG = pixel.NewSprite(cardsheet.Img, cardsheet.Sprites[0])
}

type Card struct {
	RawTitle string
	title    *text.Text
	canvas   *pixelgl.Canvas

	draw      bool
	//Pos    pixel.Vec
	//Scalar float64
	//Rot    float64
	//Mat    pixel.Matrix
	//interX *gween.Tween
	//interY *gween.Tween
	//interR *gween.Tween
	//interS *gween.Tween
	Transform *animation.Transform
	PosEffect *animation.TransformEffect
	ScaEffect *animation.TransformEffect
	trans     bool

	isPlay    bool
	played    bool
	canCancel bool
	sections  []*CardSection
	actPtr    int
	player    *Player

	Current  CardGroup
	Previous CardGroup
	ID       uuid.UUID
}

func NewCard(title string, sections []*CardSection) *Card {
	transform := animation.NewTransform(true)
	transform.Anchor = animation.Anchor{
		H: animation.Center,
		V: animation.Center,
	}
	transform.Rect = pixel.R(0, 0, BaseCardWidth, BaseCardHeight)
	return &Card{
		RawTitle: title,
		canvas:   pixelgl.NewCanvas(transform.Rect),
		title:    text.New(pixel.ZV, text2.BasicAtlas),

		draw:      true,
		Transform: transform,

		sections:  sections,
		actPtr:    -1,
		canCancel: true,

		ID: uuid.NewV4(),
	}
}

func (c *Card) PointInside(point pixel.Vec) bool {
	return util.PointInside(point, c.canvas.Bounds(), c.Transform.Mat)
}

func (c *Card) Update(r pixel.Rect) {
	// card title
	c.title.Clear()
	c.title.Color = colornames.Black
	c.title.Dot.X -= c.title.BoundsOf(c.RawTitle).W() / 2.
	fmt.Fprintln(c.title, c.RawTitle)

	for _, cs := range c.sections {
		cs.Update()
	}

	moved := false
	// card position, scaling, and rotation
	c.Transform.Update(r)
	if c.PosEffect != nil {
		c.PosEffect.Update()
		moved = true
		if c.PosEffect.IsDone() {
			c.PosEffect = nil
		}
	}
	if c.ScaEffect != nil {
		c.ScaEffect.Update()
		moved = true
		if c.ScaEffect.IsDone() {
			c.ScaEffect = nil
		}
	}
	//if c.interX != nil {
	//	x, finX := c.interX.Update(timing.DT)
	//	c.Pos.X = x
	//	if finX {
	//		c.interX = nil
	//	} else {
	//		moved = true
	//	}
	//}
	//if c.interY != nil {
	//	y, finY := c.interY.Update(timing.DT)
	//	c.Pos.Y = y
	//	if finY {
	//		c.interY = nil
	//	} else {
	//		moved = true
	//	}
	//}
	//if c.interS != nil {
	//	s, finS := c.interS.Update(timing.DT)
	//	c.Scalar = s
	//	if finS {
	//		c.interS = nil
	//	} else {
	//		moved = true
	//	}
	//}
	c.trans = c.trans && moved

	if c.isPlay {
		if c.actPtr >= len(c.sections) {
			c.played = true
			c.stop()
		} else {
			section := c.sections[c.actPtr]
			if !section.start {
				section.action.Values.Source = c.player.Character
				c.player.SetPlayerAction(section.action)
				section.start = true
			} else if section.action.Complete {
				section.isDone = true
				c.actPtr++
				c.canCancel = false
			} else {
				section.action.Selector.SetValues(section.action.Values)
			}
		}
	}
}

func (c *Card) Draw(target pixel.Target) {
	c.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
	CardBG.Draw(c.canvas, pixel.IM.Moved(pixel.V(BaseCardWidth* 0.5, BaseCardHeight* 0.5)))
	for i, cs := range c.sections {
		cs.Draw(c.canvas, i)
	}
	c.title.Draw(c.canvas, pixel.IM.Scaled(c.title.Orig, 2.0).Moved(pixel.V(BaseCardWidth* 0.5, BaseCardHeight- 32.0)))
	c.Transform.Mat = camera.Cam.UITransform(c.Transform.RPos, c.Transform.Scalar, c.Transform.Rot)
	c.canvas.Draw(target, c.Transform.Mat)
}

func (c *Card) play(player *Player) {
	c.actPtr = 0
	c.isPlay = true
	c.player = player
	c.played = false
	c.canCancel = true
	for _, section := range c.sections {
		section.start = false
		section.isDone = false
		section.action.Complete = false
	}
}

func (c *Card) stop() {
	c.actPtr = -1
	c.isPlay = false
	c.canCancel = true
	c.player.CurrAction = nil
	for _, section := range c.sections {
		section.start = true
		section.isDone = true
		section.action.Complete = false
	}
}

func (c *Card) setXY(v pixel.Vec) {
	transform := animation.TransformBuilder{
		Transform: c.Transform,
		InterX:    gween.New(c.Transform.Pos.X, v.X, 0.2, ease.InOutQuad),
		InterY:    gween.New(c.Transform.Pos.Y, v.Y, 0.2, ease.InOutQuad),
	}
	c.PosEffect = transform.Build()
}

func (c *Card) setScalar(s float64) {
	transform := animation.TransformBuilder{
		Transform: c.Transform,
		InterSX:   gween.New(c.Transform.Scalar.X, s, 0.2, ease.InOutCubic),
		InterSY:   gween.New(c.Transform.Scalar.Y, s, 0.2, ease.InOutCubic),
	}
	c.ScaEffect = transform.Build()
}

type CardSection struct {
	rawText string
	action  *PlayerAction
	text    *text.Text
	isDone  bool
	start   bool
}

func NewCardSection(rawText string, action *PlayerAction) *CardSection {
	return &CardSection{
		rawText: rawText,
		action:  action,
		text:    text.New(pixel.ZV, text2.BasicAtlas),
	}
}

func (cs *CardSection) Update() {
	cs.text.Clear()
	cs.text.Color = colornames.Black
	cs.text.Dot.X -= cs.text.BoundsOf(cs.rawText).W() / 2.
	fmt.Fprintln(cs.text, cs.rawText)
}

func (cs *CardSection) Draw(canvas *pixelgl.Canvas, offset int) {
	cs.text.Draw(canvas, pixel.IM.Scaled(cs.text.Orig, 1.2).Moved(pixel.V(BaseCardWidth*0.5, BaseCardHeight*0.5 - float64(offset) * (text2.BasicAtlas.LineHeight() + 20.))))
}