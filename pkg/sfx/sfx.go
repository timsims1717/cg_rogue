package sfx

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"time"
)

const sampleRate = beep.SampleRate(44100)

func init() {
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}
}