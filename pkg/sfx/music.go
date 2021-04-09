package sfx

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/pkg/errors"
	gween "github.com/timsims1717/cg_rogue_go/pkg/gween64"
	"github.com/timsims1717/cg_rogue_go/pkg/gween64/ease"
	"github.com/timsims1717/cg_rogue_go/pkg/timing"
	"math/rand"
	"os"
	"time"
)

var MusicPlayer *musicPlayer

type musicPlayer struct {
	curr     []string
	next     string
	tracks   map[string]string
	ctrl     *beep.Ctrl
	volume   *effects.Volume
	interV   *gween.Tween
	format   beep.Format
	silent   bool
	volNum   float64
	wait     float64
	variance float64
}

func init() {
	MusicPlayer = &musicPlayer{
		tracks: make(map[string]string),
	}
}

func (p *musicPlayer) Update() {
	if p.next != "" {
		if p.volume == nil || p.volume.Silent {
			if err := p.loadTrack(p.next); err != nil {
				fmt.Printf("music player error %s: %s\n", p.next, err)
			}
			p.next = ""
		}
	}
	if p.volume != nil {
		if p.interV != nil {
			v, fin := p.interV.Update(timing.DT)
			speaker.Lock()
			if fin {
				p.volume.Silent = true
				p.silent = true
				p.interV = nil
			} else {
				p.volume.Volume = v
			}
			speaker.Unlock()
		}
		if p.silent != p.volume.Silent {
			speaker.Lock()
			p.volume.Silent = p.silent
			speaker.Unlock()
		}
		if p.volNum != p.volume.Volume {
			speaker.Lock()
			p.volume.Volume = p.volNum
			speaker.Unlock()
		}
	}
}

func (p *musicPlayer) RegisterMusicTrack(path, key string) {
	p.tracks[key] = path
}

func (p *musicPlayer) SetCurrentTracks(keys []string) {
	p.curr = keys
}

func (p *musicPlayer) PlayTrack(key string, fadeOut, wait, variance float64) {
	p.next = key
	p.wait = wait
	p.variance = variance
	if p.volume != nil {
		p.interV = gween.New(p.volume.Volume, -10., fadeOut, ease.Linear)
	}
}

func (p *musicPlayer) PlayNextTrack(fadeOut, wait, variance float64) {
	if len(p.curr) > 0 {
		p.PlayTrack(p.curr[rand.Intn(len(p.curr))], fadeOut, wait, variance)
	}
}

func (p *musicPlayer) FadeOut(fade float64) {
	if p.volume != nil {
		p.interV = gween.New(p.volume.Volume, -10., fade, ease.Linear)
	}
}

func (p *musicPlayer) Silence(s bool) {
	p.silent = s
}

func (p *musicPlayer) SetVolume(v float64) {
	p.volNum = v
}

func (p *musicPlayer) loadTrack(key string) error {
	errMsg := "load track"
	if path, ok := p.tracks[key]; ok {
		file, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, errMsg)
		}
		streamer, format, err := mp3.Decode(file)
		if err != nil {
			return errors.Wrap(err, errMsg)
		}
		speaker.Clear()
		p.ctrl = &beep.Ctrl{
			Streamer: streamer,
			Paused:   false,
		}
		p.volume = &effects.Volume{
			Streamer: p.ctrl,
			Base:     2,
			Volume:   0,
			Silent:   false,
		}
		p.volNum = 0
		p.silent = false
		p.format = format
		speaker.Play(beep.Seq(
			beep.Callback(func() {
				v := (rand.Float64() * 2. - 1.) * p.variance
				time.Sleep(time.Duration(p.wait + v) * time.Second)
			}),
			p.volume,
			beep.Callback(func() {
			if len(p.curr) > 0 {
				p.PlayTrack(p.curr[rand.Intn(len(p.curr) - 1)], 0., 8., 5.)
			}
		})))
		return nil
	}
	return errors.Wrap(fmt.Errorf("key %s is not a registered track", key), errMsg)
}

func (p *musicPlayer) stopMusic() {
	speaker.Clear()
	p.ctrl = nil
	p.volume = nil
}