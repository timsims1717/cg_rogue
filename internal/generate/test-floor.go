package generate

import (
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
	"github.com/timsims1717/cg_rogue_go/internal/floor"
	"github.com/timsims1717/cg_rogue_go/internal/player"
	"github.com/timsims1717/cg_rogue_go/pkg/camera"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
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

	if player.Player1 == nil || player.Player1.Character == nil {
		panic(errors.New("player or player character was nil"))
	}

	floor.CurrentFloor = floor.NewFloor(12, 12, spritesheet)

	pX := rand.Intn(4) + 4
	pY := rand.Intn(4) + 4

	c := world.Coords{X: pX, Y: pY}
	floor.CurrentFloor.PutOccupant(player.Player1.Character, c)

	countBase := 4 + level
	random := (rand.Intn(countBase) - countBase/2) / 2

	enemyCount := countBase + random

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
		occCoords = append(occCoords, world.Coords{X: eX, Y: eY})
		enemyLevel := CreateTestCharacter(world.Coords{X: eX, Y: eY}, treesheet)
		i += enemyLevel
	}

	camera.Cam.CenterOn([]pixel.Vec{player.Player1.Character.Transform.Pos})
}
