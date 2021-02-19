package camera

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"math"
)

var Cam *Camera

type Camera struct {
	Mat  pixel.Matrix
	Pos  pixel.Vec
	Zoom float64
	Opt  CameraOptions
}

type CameraOptions struct {
	ScrollSpeed float64
	ZoomSpeed   float64
}

func New() *Camera {
	return &Camera{
		Mat:  pixel.IM,
		Pos:  pixel.ZV,
		Zoom: 1.0,
		Opt:  CameraOptions{
			ScrollSpeed: 500.0,
			ZoomSpeed:   1.2,
		},
	}
}

func (c *Camera) Update(win *pixelgl.Window) {
	c.Mat = pixel.IM.Scaled(c.Pos, c.Zoom).Moved(win.Bounds().Center().Sub(c.Pos))
	win.SetMatrix(c.Mat)
}

func (c *Camera) Left() {
	c.Pos.X -= c.Opt.ScrollSpeed * timing.DT
}

func (c *Camera) Right() {
	c.Pos.X += c.Opt.ScrollSpeed * timing.DT
}

func (c *Camera) Down() {
	c.Pos.Y -= c.Opt.ScrollSpeed * timing.DT
}

func (c *Camera) Up() {
	c.Pos.Y += c.Opt.ScrollSpeed * timing.DT
}

func (c *Camera) ZoomIn(zoom float64) {
	c.Zoom *= math.Pow(c.Opt.ZoomSpeed, zoom)
}