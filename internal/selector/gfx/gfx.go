package gfx

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
)

var (
	SelectionSprites map[string]*pixel.Sprite
	batch            *pixel.Batch
)

func SetSpriteSheet(sheet *img.SpriteSheet) {
	batch = pixel.NewBatch(&pixel.TrianglesData{}, sheet.Img)
	SelectionSprites = make(map[string]*pixel.Sprite)
	for k, r := range sheet.SpriteMap {
		SelectionSprites[k] = pixel.NewSprite(sheet.Img, r)
	}
}

func Clear() {
	batch.Clear()
}

func Batch() *pixel.Batch {
	return batch
}

func Draw(target pixel.Target) {
	batch.Draw(target)
}