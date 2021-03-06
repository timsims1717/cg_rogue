package img

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

var (
	IM       = pixel.IM
	Flip     = pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., 1.))
	Flop     = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1., -1.))
	FlipFlop = pixel.IM.ScaledXY(pixel.ZV, pixel.V(-1., -1.))
)

type Sprite struct {
	S *pixel.Sprite
	M pixel.Matrix
}

type SpriteSheet struct {
	Img       pixel.Picture
	Sprites   []pixel.Rect
	SpriteMap map[string]pixel.Rect
}

type spriteFile struct {
	ImgFile   string   `json:"img"`
	Sprites   []sprite `json:"sprites"`
	Width     float64  `json:"width"`
	Height    float64  `json:"height"`
	SingleRow bool     `json:"singleRow"`
}

type sprite struct {
	K string  `json:"key"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

func LoadSpriteSheet(path string) (*SpriteSheet, error) {
	errMsg := "load sprite sheet"
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var fileSheet spriteFile
	err = decoder.Decode(&fileSheet)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	img, err := LoadImage(fmt.Sprintf("%s/%s", filepath.Dir(path), fileSheet.ImgFile))
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	sheet := &SpriteSheet{
		Img:       img,
		Sprites:   make([]pixel.Rect, 0),
		SpriteMap: make(map[string]pixel.Rect, 0),
	}
	x := 0.0
	for _, r := range fileSheet.Sprites {
		var rect pixel.Rect
		w := fileSheet.Width
		h := fileSheet.Height
		if r.W > 0.0 {
			w = r.W
		}
		if fileSheet.SingleRow {
			rect = pixel.R(x, 0.0, x+w, h)
			x += w
		} else {
			if r.H > 0.0 {
				h = r.H
			}
			rect = pixel.R(r.X, r.Y, r.X+w, r.Y+h)
		}
		sheet.Sprites = append(sheet.Sprites, rect)
		if r.K != "" {
			sheet.SpriteMap[r.K] = rect
		}
	}
	return sheet, nil
}
