package camera

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/pkg/animation"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

var Cam *Camera

type Camera struct {
	Mat    pixel.Matrix
	Pos    pixel.Vec
	zoom   float64
	zStep  float64
	Opt    Options
	Mask   color.RGBA
	Effect *animation.ColorEffect

	interX *gween.Tween
	interY *gween.Tween
	interZ *gween.Tween
	lock   bool
}

type Options struct {
	ScrollSpeed float64
	ZoomStep    float64
	ZoomSpeed   float64
}

func New() *Camera {
	return &Camera{
		Mat:   pixel.IM,
		Pos:   pixel.ZV,
		zoom:  1.0,
		zStep: 1.0,
		Opt: Options{
			ScrollSpeed: 500.0,
			ZoomStep:    1.2,
			ZoomSpeed:   0.2,
		},
		Mask: colornames.Black,
	}
}

func (c *Camera) Moving() bool {
	return c.lock
}

func (c *Camera) Update(win *pixelgl.Window) {
	fin := true
	if c.interX != nil {
		x, finX := c.interX.Update(timing.DT)
		c.Pos.X = x
		if finX {
			c.interX = nil
		} else {
			fin = false
		}
	}
	if c.interY != nil {
		y, finY := c.interY.Update(timing.DT)
		c.Pos.Y = y
		if finY {
			c.interY = nil
		} else {
			fin = false
		}
	}
	if c.interZ != nil {
		z, finZ := c.interZ.Update(timing.DT)
		c.zoom = z
		if finZ {
			c.interZ = nil
		} else {
			fin = false
		}
	}
	if fin && c.lock {
		c.lock = false
	}
	if c.Effect != nil {
		c.Effect.Update()
		if c.Effect.IsDone() {
			c.Effect = nil
		}
	}
	c.Mat = pixel.IM.Scaled(c.Pos, c.zoom).Moved(win.Bounds().Center().Sub(c.Pos))
	win.SetMatrix(c.Mat)
	win.SetColorMask(c.Mask)
}

func (c *Camera) Stop() {
	c.lock = false
	c.interX = nil
	c.interY = nil
}

func (c *Camera) MoveTo(v pixel.Vec, dur float64, lock bool) {
	if !c.lock {
		c.interX = gween.New(c.Pos.X, v.X, dur, ease.InOutQuad)
		c.interY = gween.New(c.Pos.Y, v.Y, dur, ease.InOutQuad)
		c.lock = lock
	}
}

func (c *Camera) CenterOn(points []pixel.Vec) {
	if !c.lock {
		if points == nil || len(points) == 0 {
			return
		} else if len(points) == 1 {
			c.Pos = points[0]
		} else {
			// todo: center on multiple points + change zoom
		}
	}
}

func (c *Camera) Left() {
	if !c.lock {
		c.Pos.X -= c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) Right() {
	if !c.lock {
		c.Pos.X += c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) Down() {
	if !c.lock {
		c.Pos.Y -= c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) Up() {
	if !c.lock {
		c.Pos.Y += c.Opt.ScrollSpeed * timing.DT
	}
}

func (c *Camera) SetZoom(zoom float64) {
	c.zoom = zoom
	c.zStep = zoom
}

func (c *Camera) ZoomIn(zoom float64) {
	if !c.lock {
		c.zStep *= math.Pow(c.Opt.ZoomStep, zoom)
		c.interZ = gween.New(c.zoom, c.zStep, c.Opt.ZoomSpeed, ease.OutQuad)
	}
}

// UITransform returns a pixel.Matrix that can move the center of a pixel.Rect
// to the bottom left of the screen.
func (c *Camera) UITransform(pos, scalar pixel.Vec, rot float64) pixel.Matrix {
	zoom := 1 / c.zoom
	mat := pixel.IM.ScaledXY(pixel.ZV, scalar.Scaled(zoom))
	mat = mat.Rotated(pixel.ZV, rot)
	mat = mat.Moved(pixel.V(c.Pos.X, c.Pos.Y))
	mat = mat.Moved(pixel.V(WindowWidthF, WindowHeightF).Scaled(-0.5 * zoom))
	mat = mat.Moved(pos.Scaled(zoom))
	return mat
}

func (c *Camera) GetColor() color.RGBA {
	return c.Mask
}

func (c *Camera) SetColor(col color.RGBA) {
	c.Mask = col
}
