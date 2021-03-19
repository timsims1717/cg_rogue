package generate

import (
	"github.com/faiface/pixel"
	"github.com/timsims1717/cg_rogue_go/internal/characters"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

func LoadTestFloor(level int) {
	treesheet, err := img.LoadSpriteSheet("assets/character/trees.json")
	if err != nil {
		panic(err)
	}
	spritesheet, err := img.LoadSpriteSheet("assets/img/testfloor.json")
	if err != nil {
		panic(err)
	}
	charsheet, err := img.LoadSpriteSheet("assets/character/testmananim.json")
	if err != nil {
		panic(err)
	}
	floor.CurrentFloor = floor.NewFloor(12, 12, spritesheet)

	pX := rand.Intn(4) + 4
	pY := rand.Intn(4) + 4
	character := characters.NewCharacter(pixel.NewSprite(charsheet.Img, charsheet.Sprites[rand.Intn(len(charsheet.Sprites))]), world.Coords{pX, pY}, characters.Ally, 10)
	player.Player1 = player.NewPlayer(character)

	characters.CharacterManager.Add(character)

	enemyCount := 4 + util.Max((level * 2) + rand.Intn(4) - 2, 0)

	var occCoords []world.Coords
	for i := -2; i < 3; i++ {
		for j := -2; j < 3; j++ {
			occCoords = append(occCoords, world.Coords{
				X: pX + i,
				Y: pY + j,
			})
		}
	}
	for i := 0; i < enemyCount; {
		eX := pX
		eY := pY
		found := true
		for found {
			found = false
			for _, c := range occCoords {
				if c.X == eX && c.Y == eY {
					eX = rand.Intn(12)
					eY = rand.Intn(12)
					found = true
					break
				}
			}
		}
		occCoords = append(occCoords, world.Coords{ X: eX, Y: eY })
		enemyLevel := CreateTestCharacter(world.Coords{X: eX, Y: eY}, treesheet)
		i += enemyLevel
	}

	camera.Cam.CenterOn([]pixel.Vec{character.Transform.Pos})
}