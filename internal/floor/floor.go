package floor

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/timsims1717/cg_rogue_go/internal/cfg"
	"github.com/timsims1717/cg_rogue_go/internal/objects"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

var CurrentFloor Floor

type Floor struct{
	floor  [][]Hex
	batch  *pixel.Batch
	update bool
	checks pathChecks
}

type pathChecks struct {
	unoccupied bool
	nonEmpty   bool
	orig       world.Coords
}

func DefaultFloor(w, h int, spriteSheet *img.SpriteSheet) {
	if w <= 0 || h <= 0 {
		panic(fmt.Errorf("could not create floor with width of %d and height of %d", w, h))
	}
	CurrentFloor = Floor{
		batch: pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet.Img),
		update: true,
	}
	CurrentFloor.floor = make([][]Hex, 0)
	for x := 0; x < w; x++ {
		CurrentFloor.floor = append(CurrentFloor.floor, make([]Hex, 0))
		for y := 0; y < h; y++ {
			CurrentFloor.floor[x] = append(CurrentFloor.floor[x], NewHex(&CurrentFloor, x, y, pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[rand.Intn(len(spriteSheet.Sprites))])))
		}
	}
}

func (f *Floor) Draw(win *pixelgl.Window) {
	if f.update {
		f.batch.Clear()
		w, h := f.Dimensions()
		for y := h - 1; y >= 0; y-- {
			for x := 1; x < w; x += 2 {
				hex := f.Get(world.Coords{x,y})
				mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(pixel.V(world.MapToWorldHex(hex.X, hex.Y))).Moved(pixel.V(-4.0, 0.0))
				hex.Tile.Draw(f.batch, mat)
			}
			for x := 0; x < w; x += 2 {
				hex := f.Get(world.Coords{x,y})
				mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(pixel.V(world.MapToWorldHex(hex.X, hex.Y))).Moved(pixel.V(-4.0, 0.0))
				hex.Tile.Draw(f.batch, mat)
			}
		}
		f.update = false
	}
	f.batch.Draw(win)
}

func (f *Floor) Dimensions() (int, int) {
	width := len(f.floor)
	height := len(f.floor[0])
	return width, height
}

func (f *Floor) Get(a world.Coords) *Hex {
	if f.Exists(a) {
		return &(f.floor[a.X][a.Y])
	}
	return nil
}

func (f *Floor) Exists(a world.Coords) bool {
	w, h := f.Dimensions()
	return a.X >= 0 && a.Y >= 0 && a.X < w && a.Y < h
}

func (f *Floor) IsOccupied(a world.Coords) bool {
	hex := f.Get(a)
	return hex == nil || objects.NotNilOcc(hex.Occupant)
}

func (f *Floor) GetOccupant(a world.Coords) objects.Occupant {
	hex := f.Get(a)
	if hex != nil && objects.NotNilOcc(hex.Occupant) {
		return hex.Occupant
	}
	return nil
}

func (f *Floor) PutOccupant(o objects.Occupant, a world.Coords) bool {
	hex := f.Get(a)
	if hex != nil && !objects.NotNilOcc(hex.Occupant) {
		hex.Occupant = o
		return true
	}
	return false
}

func (f *Floor) RemoveOccupant(a world.Coords) bool {
	hex := f.Get(a)
	if hex != nil && objects.NotNilOcc(hex.Occupant) {
		hex.Occupant = nil
		return true
	}
	return false
}

func (f *Floor) MoveOccupant(o objects.Occupant, a, b world.Coords) bool {
	if !f.Exists(a) || !f.Exists(b) {
		return false
	}
	success := f.RemoveOccupant(a)
	if success {
		success = f.PutOccupant(o, b)
	}
	return success
}

func (f *Floor) IfLegal(a world.Coords) *Hex {
	hex := f.Get(a)
	if hex != nil {
		if (a.X == f.checks.orig.X && a.Y == f.checks.orig.Y) || ((!f.checks.unoccupied || hex.Occupant == nil) && (!f.checks.nonEmpty || !hex.Empty)) {
			return hex
		}
	}
	return nil
}

func (f *Floor) FindPath(a, b world.Coords, unoccupied, nonEmpty bool) ([]*Hex, int, bool) {
	f.checks = pathChecks{
		unoccupied: unoccupied,
		nonEmpty:   nonEmpty,
		orig:       a,
	}
	pathA, distance, found := astar.Path(f.Get(b), f.Get(a))
	var path []*Hex
	for _, h := range pathA {
		path = append(path, h.(*Hex))
	}
	f.checks = pathChecks{
		unoccupied: false,
		nonEmpty:   false,
		orig:       world.Coords{},
	}
	return path, int(distance), found
}

// Neighbors returns each legal hex
func (f *Floor) Neighbors(hex *Hex) []*Hex {
	width, height := f.Dimensions()
	neighbors := make([]*Hex, 0)
	if hex.Y > 0 {
		if n := f.IfLegal(world.Coords{hex.X, hex.Y-1}); n != nil {
			neighbors = append(neighbors, n)
		}
	}
	if hex.X < width - 1 {
		if n := f.IfLegal(world.Coords{hex.X+1, hex.Y}); n != nil {
			neighbors = append(neighbors, n)
		}
	}
	if hex.Y < height - 1 {
		if n := f.IfLegal(world.Coords{hex.X, hex.Y+1}); n != nil {
			neighbors = append(neighbors, n)
		}
	}
	if hex.X > 0 {
		if n := f.IfLegal(world.Coords{hex.X-1, hex.Y}); n != nil {
			neighbors = append(neighbors, n)
		}
	}
	if hex.X % 2 == 0 {
		if hex.X < width - 1 && hex.Y > 0 {
			if n := f.IfLegal(world.Coords{hex.X+1, hex.Y-1}); n != nil {
				neighbors = append(neighbors, n)
			}
		}
		if hex.X > 0 && hex.Y > 0 {
			if n := f.IfLegal(world.Coords{hex.X-1, hex.Y-1}); n != nil {
				neighbors = append(neighbors, n)
			}
		}
	} else {
		if hex.X < width - 1 && hex.Y < height - 1 {
			if n := f.IfLegal(world.Coords{hex.X+1, hex.Y+1}); n != nil {
				neighbors = append(neighbors, n)
			}
		}
		if hex.X > 0 && hex.Y < height - 1 {
			if n := f.IfLegal(world.Coords{hex.X-1, hex.Y+1}); n != nil {
				neighbors = append(neighbors, n)
			}
		}
	}
	return neighbors
}
