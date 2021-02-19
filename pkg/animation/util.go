package animation

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/timsims1717/cg_rogue_go/pkg/img"
	"os"
	"path/filepath"
)

type AnimationCfg struct {
	Sprites         *img.SpriteSheet     `json:"-"`
	SpriteSheetFile string               `json:"spriteSheet"`
	FrameData       map[string]FrameData `json:"frameData"`
	FPS             int                  `json:"fps"`
	State           string
	counter         int
}

func (a *AnimationCfg) Update(dt float64) {

}

type FrameData struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func LoadAnimationCfg(path string) (*AnimationCfg, error) {
	errMsg := "load animation cfg"
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var animCfg AnimationCfg
	err = decoder.Decode(&animCfg)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	animCfg.Sprites, err = img.LoadSpriteSheet(fmt.Sprintf("%s/%s", filepath.Dir(path), animCfg.SpriteSheetFile))
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	return &animCfg, nil
}