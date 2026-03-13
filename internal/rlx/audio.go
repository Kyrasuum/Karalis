package rlx

import rl "github.com/gen2brain/raylib-go/raylib"

func InitAudioDevice() {
	Do(func() {
		rl.InitAudioDevice()
	})
}

func CloseAudioDevice() {
	Do(func() {
		rl.CloseAudioDevice()
	})
}

func IsAudioDeviceReady() bool {
	return Call(func() bool {
		return rl.IsAudioDeviceReady()
	})
}

func SetMasterVolume(volume float32) {
	Do(func() {
		rl.SetMasterVolume(volume)
	})
}

func GetMasterVolume() float32 {
	return Call(func() float32 {
		return rl.GetMasterVolume()
	})
}

func LoadSound(fileName string) rl.Sound {
	return Call(func() rl.Sound {
		return rl.LoadSound(fileName)
	})
}

func LoadSoundFromWave(wave rl.Wave) rl.Sound {
	return Call(func() rl.Sound {
		return rl.LoadSoundFromWave(wave)
	})
}

func LoadSoundAlias(source rl.Sound) rl.Sound {
	return Call(func() rl.Sound {
		return rl.LoadSoundAlias(source)
	})
}

func IsSoundValid(sound rl.Sound) bool {
	return Call(func() bool {
		return rl.IsSoundValid(sound)
	})
}

func UpdateSound(sound rl.Sound, data []byte, sampleCount int32) {
	buf := append([]byte(nil), data...)
	Do(func() {
		rl.UpdateSound(sound, buf, sampleCount)
	})
}

func UnloadWave(wave rl.Wave) {
	Do(func() {
		rl.UnloadWave(wave)
	})
}

func UnloadSound(sound rl.Sound) {
	Do(func() {
		rl.UnloadSound(sound)
	})
}

func UnloadSoundAlias(alias rl.Sound) {
	Do(func() {
		rl.UnloadSoundAlias(alias)
	})
}

func PlaySound(sound rl.Sound) {
	Do(func() {
		rl.PlaySound(sound)
	})
}

func StopSound(sound rl.Sound) {
	Do(func() {
		rl.StopSound(sound)
	})
}

func PauseSound(sound rl.Sound) {
	Do(func() {
		rl.PauseSound(sound)
	})
}

func ResumeSound(sound rl.Sound) {
	Do(func() {
		rl.ResumeSound(sound)
	})
}

func IsSoundPlaying(sound rl.Sound) bool {
	return Call(func() bool {
		return rl.IsSoundPlaying(sound)
	})
}

func SetSoundVolume(sound rl.Sound, volume float32) {
	Do(func() {
		rl.SetSoundVolume(sound, volume)
	})
}

func SetSoundPitch(sound rl.Sound, pitch float32) {
	Do(func() {
		rl.SetSoundPitch(sound, pitch)
	})
}

func SetSoundPan(sound rl.Sound, pan float32) {
	Do(func() {
		rl.SetSoundPan(sound, pan)
	})
}
