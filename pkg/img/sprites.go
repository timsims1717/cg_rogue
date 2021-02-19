package img

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type SpriteSheet struct {
	Img     pixel.Picture
	Sprites []pixel.Rect
}

type spriteFile struct {
	ImgFile   string `json:"img"`
	Sprites   []rect `json:"sprites"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	SingleRow bool    `json:"singleRow"`
}

type rect struct {
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
		Img: img,
		Sprites: make([]pixel.Rect, 0),
	}
	x := 0.0
	for _, r := range fileSheet.Sprites {
		if fileSheet.SingleRow {
			h := fileSheet.Height
			w := fileSheet.Width
			if r.W > 0.0 {
				w = r.W
			}
			sheet.Sprites = append(sheet.Sprites, pixel.R(x, 0.0, x+w, h))
			x += w
		} else {
			w := fileSheet.Width
			h := fileSheet.Height
			if r.W > 0.0 {
				w = r.W
			}
			if r.H > 0.0 {
				h = r.H
			}
			sheet.Sprites = append(sheet.Sprites, pixel.R(r.X, r.Y, r.X+w, r.Y+h))
		}
	}
	return sheet, nil
}