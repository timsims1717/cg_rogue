package floor

import (
	"fmt"
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
	floor    [][]Hex
	batch    *pixel.Batch
	update   bool
	checks   PathChecks
	PathLine pixel.Line
}

type PathChecks struct {
	NotFilled     bool // true: must not be filled, false: can be a filled tile
	Unoccupied    bool // true: must be unoccupied, false: can have an occupant
	NonEmpty      bool // true: must not be empty, false: can be an empty tile (a pit, or something)
	EndUnoccupied bool // true: last tile must be unoccupied, false: last tile can have an occupant
	Orig          world.Coords
}

var DefaultCheck = PathChecks{
	NotFilled:     true,
	Unoccupied:    false,
	NonEmpty:      false,
	EndUnoccupied: false,
	Orig:          world.Coords{},
}

var NoCheck = PathChecks{
	NotFilled:     false,
	Unoccupied:    false,
	NonEmpty:      false,
	EndUnoccupied: false,
	Orig:          world.Coords{},
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
				mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(pixel.V(world.MapToWorld(hex.X, hex.Y))).Moved(pixel.V(-4.0, 0.0))
				hex.Tile.Draw(f.batch, mat)
			}
			for x := 0; x < w; x += 2 {
				hex := f.Get(world.Coords{x,y})
				mat := pixel.IM.Scaled(pixel.ZV, cfg.Scalar).Moved(pixel.V(world.MapToWorld(hex.X, hex.Y))).Moved(pixel.V(-4.0, 0.0))
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

func (f *Floor) SetDefaultChecks() {
	f.checks = DefaultCheck
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
	return hex == nil || objects.NotNilMov(hex.Occupant)
}

func (f *Floor) GetOccupant(a world.Coords) objects.Moveable {
	hex := f.Get(a)
	if hex != nil && objects.NotNilMov(hex.Occupant) {
		return hex.Occupant
	}
	return nil
}

func (f *Floor) HasOccupant(a world.Coords) bool {
	hex := f.Get(a)
	return hex != nil && objects.NotNilMov(hex.Occupant)
}

func (f *Floor) PutOccupant(m objects.Moveable, a world.Coords) bool {
	hex := f.Get(a)
	if hex != nil && !objects.NotNilMov(hex.Occupant) {
		hex.Occupant = m
		return true
	}
	return false
}

func (f *Floor) RemoveOccupant(a world.Coords) bool {
	hex := f.Get(a)
	if hex != nil && objects.NotNilMov(hex.Occupant) {
		hex.Occupant = nil
		return true
	}
	return false
}

func (f *Floor) MoveOccupant(m objects.Moveable, a, b world.Coords) bool {
	if !f.Exists(a) || !f.Exists(b) {
		return false
	}
	success := f.RemoveOccupant(a)
	if success {
		success = f.PutOccupant(m, b)
	}
	return success
}