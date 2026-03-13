package rlx

import rl "github.com/gen2brain/raylib-go/raylib"

func SetTargetFPS(fps int32) {
	Do(func() {
		rl.SetTargetFPS(fps)
	})
}

func GetFrameTime() float32 {
	return Call(func() float32 {
		return rl.GetFrameTime()
	})
}

func GetTime() float64 {
	return Call(func() float64 {
		return rl.GetTime()
	})
}

func GetFPS() int32 {
	return Call(func() int32 {
		return rl.GetFPS()
	})
}

func SwapScreenBuffer() {
	Do(func() {
		rl.SwapScreenBuffer()
	})
}

func PollInputEvents() {
	Do(func() {
		rl.PollInputEvents()
	})
}

func WaitTime(seconds float64) {
	Do(func() {
		rl.WaitTime(seconds)
	})
}
