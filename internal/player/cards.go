package player

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	text2 "github.com/timsims1717/cg_rogue_go/pkg/text"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
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
	rawText  string
	RawTitle string
	desc     *text.Text
	title    *text.Text
	canvas   *pixelgl.Canvas
	updated  bool

	draw   bool
	Pos    pixel.Vec
	Scalar float64
	Rot    float64
	Mat    pixel.Matrix
	interX *gween.Tween
	interY *gween.Tween
	interR *gween.Tween
	interS *gween.Tween

	isPlay    bool
	canCancel bool
	isDone    bool
	action    *PlayerAction
	actPtr    int
}

func NewCard(title string, rawText string, action *PlayerAction) *Card {
	return &Card{
		rawText:  rawText,
		RawTitle: title,
		canvas:   pixelgl.NewCanvas(pixel.R(0, 0, BaseCardWidth, BaseCardHeight)),
		desc:     text.New(pixel.ZV, text2.BasicAtlas),
		title:    text.New(pixel.ZV, text2.BasicAtlas),
		updated:  true,

		draw:   true,
		Pos:    pixel.ZV,
		Scalar: 1.0,
		Rot:    0.0,
		Mat:    pixel.IM,

		action: action,
		actPtr: -1,
		canCancel: true,
	}
}

func (c *Card) PointInside(point pixel.Vec) bool {
	if c.canvas.Bounds().Moved(pixel.V(-(BaseCardWidth / 2.0), -(BaseCardHeight / 2.0))).Contains(c.Mat.Unproject(point)) {
		return true
	}
	return false
}

func (c *Card) Update() {
	// card text
	c.desc.Clear()
	c.desc.Color = colornames.Black
	c.desc.Dot.X -= c.desc.BoundsOf(c.rawText).W() / 2.
	fmt.Fprintln(c.desc, c.rawText)

	c.title.Clear()
	c.title.Color = colornames.Black
	c.title.Dot.X -= c.title.BoundsOf(c.RawTitle).W() / 2.
	fmt.Fprintln(c.title, c.RawTitle)

	// card position, scaling, and rotation
	if c.interX != nil {
		x, finX := c.interX.Update(timing.DT)
		c.Pos.X = x
		if finX {
			c.interX = nil
		}
	}
	if c.interY != nil {
		y, finY := c.interY.Update(timing.DT)
		c.Pos.Y = y
		if finY {
			c.interY = nil
		}
	}
	if c.interS != nil {
		s, finS := c.interS.Update(timing.DT)
		c.Scalar = s
		if finS {
			c.interS = nil
		}
	}
}

func (c *Card) Draw(win *pixelgl.Window) {
	c.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
	CardBG.Draw(c.canvas, pixel.IM.Moved(pixel.V(BaseCardWidth* 0.5, BaseCardHeight* 0.5)))
	c.desc.Draw(c.canvas, pixel.IM.Scaled(c.desc.Orig, 1.2).Moved(pixel.V(BaseCardWidth* 0.5, BaseCardHeight* 0.5)))
	c.title.Draw(c.canvas, pixel.IM.Scaled(c.title.Orig, 2.0).Moved(pixel.V(BaseCardWidth* 0.5, BaseCardHeight- 32.0)))
	zoom := 1/camera.Cam.Zoom
	c.Mat = pixel.IM.Scaled(pixel.ZV, zoom * c.Scalar)
	c.Mat = c.Mat.Moved(pixel.V(camera.Cam.Pos.X, camera.Cam.Pos.Y))
	c.Mat = c.Mat.Moved(pixel.V(float64(cfg.WindowWidth), float64(cfg.WindowHeight)).Scaled(-0.5 * zoom))
	c.Mat = c.Mat.Moved(c.Pos.Scaled(zoom))
	c.canvas.Draw(win, c.Mat)
}

func (c *Card) play(player *Player) {
	c.actPtr = 0
	c.isPlay = true
	c.action.Values.Source = player.Character
	player.SetPlayerAction(c.action)
}

func (c *Card) stop() {
	c.actPtr = -1
	c.isPlay = false
	c.action.Complete = true
}

func (c *Card) setXY(v pixel.Vec) {
	c.interX = gween.New(c.Pos.X, v.X, 0.2, ease.InOutQuad)
	c.interY = gween.New(c.Pos.Y, v.Y, 0.2, ease.InOutQuad)
}

func (c *Card) setScalar(s float64) {
	c.interS = gween.New(c.Scalar, s, 0.2, ease.InOutCubic)
}

type CardSection struct {
	rawText string
	action  *PlayerAction
	text    *text.Text
	isDone  bool
}

func NewCardSection(rawText string, action *PlayerAction) *CardSection {
	return &CardSection{
		rawText: rawText,
		action:  action,
		text:    text.New(pixel.ZV, text2.BasicAtlas),
	}
}