package floor

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/util"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

var CurrentFloor *Floor

type Floor struct {
	floor    [][]Hex
	batch    *pixel.Batch
	update   bool
	checks   PathChecks
	PathLine pixel.Line
	id       uuid.UUID
}

type PathChecks struct {
	NotFilled     bool // true: must not be filled, false: can be a filled tile
	Unoccupied    bool // true: must be unoccupied, false: can have an occupant
	NonEmpty      bool // true: must not be empty, false: can be an empty tile (a pit, or something)
	EndUnoccupied bool // true: last tile must be unoccupied, false: last tile can have an occupant
	HonorClaim    bool // true: can't select claimed tiles, false: doesn't care about claimed tiles
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

func NewFloor(w, h int, spriteSheet *img.SpriteSheet) *Floor {
	if w <= 0 || h <= 0 {
		panic(fmt.Errorf("could not create floor with width of %d and height of %d", w, h))
	}
	floor := &Floor{
		batch:  pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet.Img),
		update: true,
		id:     uuid.NewV4(),
	}
	floor.floor = make([][]Hex, 0)
	for x := 0; x < w; x++ {
		floor.floor = append(floor.floor, make([]Hex, 0))
		for y := 0; y < h; y++ {
			floor.floor[x] = append(floor.floor[x], NewHex(floor, x, y, pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[rand.Intn(len(spriteSheet.Sprites))])))
		}
	}
	return floor
}

func DefaultFloor(w, h int, spriteSheet *img.SpriteSheet) {
	if w <= 0 || h <= 0 {
		panic(fmt.Errorf("could not create floor with width of %d and height of %d", w, h))
	}
	CurrentFloor = &Floor{
		batch:  pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet.Img),
		update: true,
	}
	CurrentFloor.floor = make([][]Hex, 0)
	for x := 0; x < w; x++ {
		CurrentFloor.floor = append(CurrentFloor.floor, make([]Hex, 0))
		for y := 0; y < h; y++ {
			CurrentFloor.floor[x] = append(CurrentFloor.floor[x], NewHex(CurrentFloor, x, y, pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[rand.Intn(len(spriteSheet.Sprites))])))
		}
	}
}

func (f *Floor) Draw(win *pixelgl.Window) {
	if f.update {
		f.batch.Clear()
		w, h := f.Dimensions()
		for y := h - 1; y >= 0; y-- {
			for x := 1; x < w; x += 2 {
				hex := f.Get(world.Coords{x, y})
				mat := pixel.IM.Moved(world.MapToWorld(hex.GetCoords()))
				hex.Tile.Draw(f.batch, mat)
			}
			for x := 0; x < w; x += 2 {
				hex := f.Get(world.Coords{x, y})
				mat := pixel.IM.Moved(world.MapToWorld(hex.GetCoords()))
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

func (f *Floor) IsClaimed(a world.Coords) bool {
	hex := f.Get(a)
	return hex == nil || !util.IsNil(hex.Claimant) || !util.IsNil(hex.Occupant)
}

func (f *Floor) GetClaimant(a world.Coords) *Character {
	hex := f.Get(a)
	if hex != nil && !util.IsNil(hex.Claimant) {
		return hex.Claimant
	}
	return nil
}

func (f *Floor) RemoveClaim(a world.Coords) {
	hex := f.Get(a)
	if hex != nil && !util.IsNil(hex.Claimant) {
		hex.Claimant = nil
	}
}

func (f *Floor) Claim(c *Character, a world.Coords) {
	if !f.Exists(a) {
		return
	}
	hex := f.Get(a)
	if hex.Claimant != nil {
		hex.Claimant.RemoveClaim()
	}
	hex.Claimant = c
	c.Claim = a
}

func (f *Floor) IsOccupied(a world.Coords) bool {
	hex := f.Get(a)
	return hex == nil || !util.IsNil(hex.Occupant)
}

func (f *Floor) GetOccupant(a world.Coords) *Character {
	hex := f.Get(a)
	if hex != nil && !util.IsNil(hex.Occupant) {
		return hex.Occupant
	}
	return nil
}

func (f *Floor) RemoveOccupant(a world.Coords) *Character {
	hex := f.Get(a)
	if hex != nil && !util.IsNil(hex.Occupant) {
		former := hex.Occupant
		hex.Occupant = nil
		return former
	}
	return nil
}

func (f *Floor) PutOccupant(c *Character, e world.Coords) {
	if !f.Exists(e) {
		return
	}
	if c.Floor != nil && c.Floor.id == f.id	{
		f.RemoveOccupant(c.GetCoords())
	} else {
		c.Floor = f
	}
	hex := f.Get(e)
	hex.Claimant = nil
	hex.Occupant = c
	c.Coords = e
	c.SetPos(world.MapToWorld(e))
	c.OnMap = true
}
