package sfx

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"time"
)

const sampleRate = beep.SampleRate(44100)

// Volumes are stored as integers from 0 to 100.
var (
	masterVolume = 100
	masterMuted  = false
	musicVolume  = 100
	musicMuted   = false
	soundVolume  = 100
	soundMuted   = false
	sfxVolume    = map[string]int{}
	sfxMuted     = map[string]bool{}
)

func init() {
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}
}

func getMasterVolume() float64 {
	if masterMuted {
		return -8.
	} else {
		return float64(masterVolume) / 10. - 8.
	}
}

func getMusicVolume() float64 {
	if musicMuted || masterMuted {
		return -8.
	} else {
		return float64(musicVolume * masterVolume) / 1000. - 8.
	}
}

func getSoundVolume() float64 {
	if soundMuted || masterMuted {
		return -8.
	} else {
		return float64(soundVolume * masterVolume) / 1000. - 8.
	}
}

func getSfxVolume(key string) float64 {
	if sfxMuted[key] || masterMuted {
		return -8.
	} else {
		return float64(sfxVolume[key] * masterVolume) / 1000. - 8.
	}
}

func GetMasterVolume() int {
	if masterMuted {
		return 0
	} else {
		return int((masterVolume + 8.) * 10.)
	}
}

func GetMusicVolume() int {
	if musicMuted {
		return 0
	} else {
		return int((musicVolume + 8.) * 10.)
	}
}

func GetSoundVolume() int {
	if soundMuted {
		return 0
	} else {
		return int((soundVolume + 8.) * 10.)
	}
}

func GetSfxVolume(key string) int {
	if sfxMuted[key] {
		return 0
	} else {
		return int((sfxVolume[key] + 8.) * 10.)
	}
}

func SetMasterVolume(v int) {
	if v == 0 {
		masterMuted = true
	} else {
		masterMuted = false
	}
	masterVolume = v
}

func SetMusicVolume(v int) {
	if v == 0 {
		musicMuted = true
	} else {
		musicMuted = false
	}
	musicVolume = v
}

func SetSoundVolume(v int) {
	if v == 0 {
		soundMuted = true
	} else {
		soundMuted = false
	}
	soundVolume = v
}

func SetSfxVolume(v int, key string) {
	if v == 0 {
		sfxMuted[key] = true
	} else {
		sfxMuted[key] = false
	}
	sfxVolume[key] = v
}