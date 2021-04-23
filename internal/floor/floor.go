package floor

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	uuid "github.com/satori/go.uuid"
	"github.com/timsims1717/cg_rogue_go/internal/selector/gfx"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"github.com/timsims1717/cg_rogue_go/pkg/world"
	"math/rand"
)

var CurrentFloor *Floor

type Floor struct {
	floor    [][]Hex
	storage  []*Character
	batches  []*pixel.Batch
	update   []bool
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

var (
	DefaultCheck = PathChecks{
		NotFilled:     true,
	}
	NoCheck = PathChecks{}
)

func NewFloor(w, h int, spriteSheet *img.SpriteSheet) *Floor {
	if w <= 0 || h <= 0 {
		panic(fmt.Errorf("could not create floor with width of %d and height of %d", w, h))
	}
	var batches []*pixel.Batch
	var update []bool
	for i := 0; i < h * 2; i++ {
		batches = append(batches, pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet.Img))
		update = append(update, true)
	}
	floor := &Floor{
		batches: batches,
		update:  update,
		id:      uuid.NewV4(),
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

func (f *Floor) SetWorldDimensions() {
	w, h := f.Dimensions()
	world.Width = w
	world.Height = h
}

func (f *Floor) Update() {
	for _, row := range f.floor {
		for _, h := range row {
			if h.IsOccupied() {
				h.Occupant.Update()
			}
			for _, c := range h.storage {
				c.Update()
			}
		}
	}
}

func (f *Floor) Draw(win *pixelgl.Window) {
	w, h := f.Dimensions()
	i := 0
	for y := h - 1; y >= 0; y-- {
		f.drawRow(win, i, w, y)
		i++
		f.drawRow(win, i, w, y)
		i++
	}
}

func (f *Floor) drawRow(win *pixelgl.Window, i, w, y int) {
	if f.update[i] {
		f.batches[i].Clear()
		for x := (i + 1) % 2; x < w; x += 2 {
			hex := f.Get(world.Coords{X: x, Y: y})
			hex.Tile.Draw(f.batches[i], hex.matrix)
		}
		f.update[i] = false
	}
	f.batches[i].Draw(win)
	gfx.Clear()
	eff := false
	for x := (i + 1) % 2; x < w; x += 2 {
		hex := f.Get(world.Coords{X: x, Y: y})
		for j := 11; j >= 6; j-- {
			if sprites, ok := hex.effects[j]; ok {
				for _, spr := range sprites {
					spr.S.Draw(gfx.Batch(), spr.M.Moved(hex.GetWorlds()))
					eff = true
				}
			}
		}
	}
	if eff {
		gfx.Draw(win)
	}
	for x := (i + 1) % 2; x < w; x += 2 {
		hex := f.Get(world.Coords{X: x, Y: y})
		if hex.IsOccupied() {
			hex.Occupant.Draw(win)
		}
		for _, c := range hex.storage {
			c.Draw(win)
		}
	}
	gfx.Clear()
	eff = false
	for x := (i + 1) % 2; x < w; x += 2 {
		hex := f.Get(world.Coords{X: x, Y: y})
		for j := 5; j >= 0; j-- {
			if sprites, ok := hex.effects[j]; ok {
				for _, spr := range sprites {
					spr.S.Draw(gfx.Batch(), spr.M.Moved(hex.GetWorlds()))
					eff = true
				}
			}
		}
		hex.ClearEffects()
	}
	if eff {
		gfx.Draw(win)
	}
	for x := (i + 1) % 2; x < w; x += 2 {
		hex := f.Get(world.Coords{X: x, Y: y})
		if hex.IsOccupied() {
			hex.Occupant.Health.Draw(win)
			hex.Occupant.Defense.Draw(win)
		}
		for _, c := range hex.storage {
			c.Health.Draw(win)
			c.Defense.Draw(win)
		}
	}
}

func (f *Floor) Dimensions() (int, int) {
	width := len(f.floor)
	height := len(f.floor[0])
	return width, height
}

func (f *Floor) SetDefaultChecks() {
	f.checks = DefaultCheck
}

func (f *Floor) GetDiplomatic(d Diplomacy, orig world.Coords, r int) []world.Coords {
	var set []world.Coords
	within := f.AllWithin(orig, r, NoCheck)
	for _, c := range within {
		ch := f.Get(c).Occupant
		if ch != nil && ch.Diplomacy == d && !ch.IsDestroyed() {
			set = append(set, ch.Coords)
		}
	}
	return set
}

func (f *Floor) Get(a world.Coords) *Hex {
	if f != nil && f.Exists(a) {
		return &(f.floor[a.X][a.Y])
	}
	return nil
}

func (f *Floor) Exists(a world.Coords) bool {
	if f == nil {
		return false
	}
	w, h := f.Dimensions()
	return a.X >= 0 && a.Y >= 0 && a.X < w && a.Y < h
}

func (f *Floor) PutOccupant(c *Character, e world.Coords) {
	if !f.Exists(e) {
		return
	}
	stored := false
	if c.Floor != nil {
		hex := c.Floor.Get(c.GetCoords())
		if hex.IsStoring(c) {
			stored = true
			hex.UnStore(c)
		} else {
			hex.RemoveOccupant()
		}
		c.Floor.UnStore(c)
	}
	c.Floor = f
	hex := f.Get(e)
	if stored {
		hex.Store(c)
	} else {
		hex.Claimant = nil
		hex.Occupant = c
	}
	c.Coords = e
	c.OnMap = true
}

func (f *Floor) Store(c *Character) {
	if c.Floor != nil {
		hex := c.Floor.Get(c.GetCoords())
		if hex.IsStoring(c) {
			hex.UnStore(c)
		} else {
			hex.RemoveOccupant()
			c.Floor.UnStore(c)
		}
	}
	c.Floor = f
	f.storage = append(f.storage, c)
	c.OnMap = false
}

func (f *Floor) UnStore(c *Character) {
	for i, s := range f.storage {
		if s.ID() == c.ID() {
			f.storage = append(f.storage[:i], f.storage[i+1:]...)
			break
		}
	}
}

func (f *Floor) Clear() {
	f.storage = []*Character{}
	for _, row := range f.floor {
		for _, h := range row {
			h.Clear()
		}
	}
}