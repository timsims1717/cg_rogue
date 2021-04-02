package player

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/internal/selectors"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	text2 "github.com/timsims1717/cg_rogue_go/pkg/typeface"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"golang.org/x/image/colornames"
)

var CardBG *pixel.Sprite

func InitializeCards() {
	cardsheet, err := img.LoadSpriteSheet("assets/cards/testcard.json")
	if err != nil {
		panic(err)
	}
	CardBG = pixel.NewSprite(cardsheet.Img, cardsheet.Sprites[0])
}

type CardAction interface {
	DoActions()
	//PreSelect(int)
	//PostSelect(int)
	SetValues(int)
	InitSelectors()
	SetCard(*Card)
}

type Card struct {
	Action CardAction

	RawTitle string
	RawDesc  string
	title    *text.Text
	desc     *text.Text
	canvas   *pixelgl.Canvas

	draw      bool
	Transform *animation.Transform
	PosEffect *animation.TransformEffect
	ScaEffect *animation.TransformEffect
	trans     bool

	Selectors []*selectors.AbstractSelector
	Results   [][]world.Coords
	Values    selectors.ActionValues
	actPtr    int
	isPlay    bool
	played    bool
	tempOrig  []world.Coords
	Rests     int
	Level     int

	Current  CardGroup
	Previous CardGroup
	ID       uuid.UUID
}

func NewCard(title, desc string, action CardAction) *Card {
	transform := animation.NewTransform(true)
	transform.Anchor = animation.Anchor{
		H: animation.Center,
		V: animation.Center,
	}
	transform.Rect = pixel.R(0, 0, BaseCardWidth, BaseCardHeight)
	newCard := &Card{
		Action:    action,
		RawTitle:  title,
		RawDesc:   desc,
		title:     text.New(pixel.ZV, text2.BasicAtlas),
		desc:      text.New(pixel.ZV, text2.BasicAtlas),
		canvas:    pixelgl.NewCanvas(transform.Rect),
		draw:      true,
		Transform: transform,
		actPtr:    -1,
		Rests:     1,
		ID:        uuid.NewV4(),
	}
	action.SetCard(newCard)
	action.InitSelectors()
	action.SetValues(0)
	return newCard
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

	// card desc
	c.desc.Clear()
	c.desc.Color = colornames.Black
	c.desc.Dot.X -= c.desc.BoundsOf(c.RawDesc).W() / 2.
	fmt.Fprintln(c.desc, c.RawDesc)

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
	c.trans = c.trans && moved
}

func (c *Card) Draw(target pixel.Target) {
	if c.draw {
		c.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
		CardBG.Draw(c.canvas, pixel.IM.Moved(pixel.V(BaseCardWidth*0.5, BaseCardHeight*0.5)))
		c.desc.Draw(c.canvas, pixel.IM.Scaled(c.desc.Orig, 1.2).Moved(pixel.V(BaseCardWidth*0.5, BaseCardHeight*0.5)))
		//c.desc.Draw(c.canvas, pixel.IM.Scaled(pa.text.Orig, 1.2).Moved(pixel.V(BaseCardWidth*0.5, BaseCardHeight*0.5-float64(offset)*(text2.BasicAtlas.LineHeight()+20.))))
		c.title.Draw(c.canvas, pixel.IM.Scaled(c.title.Orig, 2.0).Moved(pixel.V(BaseCardWidth*0.5, BaseCardHeight-32.0)))
		c.Transform.Mat = camera.Cam.UITransform(c.Transform.RPos, c.Transform.Scalar, c.Transform.Rot)
		c.canvas.Draw(target, c.Transform.Mat)
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

func (c *Card) SetDraw(draw bool) {
	c.draw = draw
}

func (c *Card) Upgrade() {
	if c.Level < 6 {
		c.Level++
		c.Action.SetValues(c.Level)
	}
}