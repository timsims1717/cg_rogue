package generate

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/ai"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

func CreateTestCharacter(coords world.Coords, spriteSheet *img.SpriteSheet) int {
	r := rand.Intn(100)
	if r < 10 {
		CreateRandomWalker(coords, spriteSheet)
		return 1
	} else if r < 60 {
		CreateFlyChaser(coords, spriteSheet)
		return 1
	} else if r < 80 {
		CreateSkirmisher(coords, spriteSheet)
		return 1
	} else {
		CreateGrenadier(coords, spriteSheet)
		return 2
	}
}

func CreateRandomWalker(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := characters.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[8]),
		coords,
		characters.Enemy,
		2,
	)
	characters.CharacterManager.Add(enemy)
	charAi := ai.NewRandomWalker(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateFlyChaser(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := characters.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[5]),
		coords,
		characters.Enemy,
		4,
		)
	characters.CharacterManager.Add(enemy)
	charAi := ai.NewFlyChaser(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateSkirmisher(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := characters.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[0]),
		coords,
		characters.Enemy,
		6,
	)
	characters.CharacterManager.Add(enemy)
	charAi := ai.NewSkirmisher(enemy)
	ai.AIManager.AddAI(charAi)
}

func CreateGrenadier(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := characters.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[6]),
		coords,
		characters.Enemy,
		8,
	)
	characters.CharacterManager.Add(enemy)
	charAi := ai.NewGrenadier(enemy)
	ai.AIManager.AddAI(charAi)
}