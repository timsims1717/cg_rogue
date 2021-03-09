package camera

import (
	"github.com/faiface/pixel"
)

func UITransform(camera *Camera, pos, scalar pixel.Vec, rot, width, height float64) pixel.Matrix {
	zoom := 1/camera.Zoom
	mat := pixel.IM.ScaledXY(pixel.ZV, scalar.Scaled(zoom))
	mat = mat.Rotated(pixel.ZV, rot)
	mat = mat.Moved(pixel.V(camera.Pos.X, camera.Pos.Y))
	mat = mat.Moved(pixel.V(width, height).Scaled(-0.5 * zoom))
	mat = mat.Moved(pos.Scaled(zoom))
	return mat
}