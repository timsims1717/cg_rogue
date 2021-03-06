package generate

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

func CreateTestCharacter(coords world.Coords, spriteSheet *img.SpriteSheet) int {
	r := rand.Intn(100)
	if r < 7 {
		CreateRandomWalker(coords, spriteSheet)
		return 1
	} else if r < 35 {
		CreateFlyChaser(coords, spriteSheet)
		return 1
	} else if r < 60 {
		CreateSkirmisher(coords, spriteSheet)
		return 1
	} else if r < 75 {
		CreateGrenadier(coords, spriteSheet)
		return 2
	} else if r < 90 {
		CreateBruiser(coords, spriteSheet)
		return 2
	} else {
		CreateStationary(coords, spriteSheet)
		return 0
	}
}

func CreateRandomWalker(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := floor.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[8]),
		floor.CurrentFloor,
		coords,
		floor.Enemy,
		1,
	)
	charAi := ai.NewRandomWalker(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateFlyChaser(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := floor.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[5]),
		floor.CurrentFloor,
		coords,
		floor.Enemy,
		3,
	)
	charAi := ai.NewFlyChaser(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateSkirmisher(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := floor.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[0]),
		floor.CurrentFloor,
		coords,
		floor.Enemy,
		4,
	)
	charAi := ai.NewSkirmisher(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateGrenadier(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := floor.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[6]),
		floor.CurrentFloor,
		coords,
		floor.Enemy,
		5,
	)
	charAi := ai.NewGrenadier(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateBruiser(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := floor.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[4]),
		floor.CurrentFloor,
		coords,
		floor.Enemy,
		8,
	)
	charAi := ai.NewBruiser(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateStationary(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := floor.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[7]),
		floor.CurrentFloor,
		coords,
		floor.Neutral,
		6,
	)
	charAi := ai.NewStationary(enemy)
	ai.AIManager.AddAI(charAi)
}
