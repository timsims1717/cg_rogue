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
	switch rand.Intn(2) {
	case 0:
		CreateRandomWalker(coords, spriteSheet)
		return 1
	case 1:
		CreateFlyChaser(coords, spriteSheet)
		return 1
	}
	return 0
}

func CreateRandomWalker(coords world.Coords, spriteSheet *img.SpriteSheet) {
	enemy := characters.NewCharacter(
		pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[8]),
		coords,
		characters.Enemy,
		6,
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