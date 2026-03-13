package rlx

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ClearBackground(col color.RGBA) {
	Do(func() {
		rl.ClearBackground(col)
	})
}

func BeginDrawing() {
	Do(func() {
		rl.BeginDrawing()
	})
}

func EndDrawing() {
	Do(func() {
		rl.EndDrawing()
	})
}

func BeginMode2D(camera rl.Camera2D) {
	Do(func() {
		rl.BeginMode2D(camera)
	})
}

func EndMode2D() {
	Do(func() {
		rl.EndMode2D()
	})
}

func BeginMode3D(camera rl.Camera3D) {
	Do(func() {
		rl.BeginMode3D(camera)
	})
}

func EndMode3D() {
	Do(func() {
		rl.EndMode3D()
	})
}

func BeginTextureMode(target rl.RenderTexture2D) {
	Do(func() {
		rl.BeginTextureMode(target)
	})
}

func EndTextureMode() {
	Do(func() {
		rl.EndTextureMode()
	})
}

func BeginShaderMode(shader rl.Shader) {
	Do(func() {
		rl.BeginShaderMode(shader)
	})
}

func EndShaderMode() {
	Do(func() {
		rl.EndShaderMode()
	})
}

func BeginBlendMode(mode rl.BlendMode) {
	Do(func() {
		rl.BeginBlendMode(mode)
	})
}

func EndBlendMode() {
	Do(func() {
		rl.EndBlendMode()
	})
}

func BeginScissorMode(x int32, y int32, width int32, height int32) {
	Do(func() {
		rl.BeginScissorMode(x, y, width, height)
	})
}

func EndScissorMode() {
	Do(func() {
		rl.EndScissorMode()
	})
}

func BeginVrStereoMode(config rl.VrStereoConfig) {
	Do(func() {
		rl.BeginVrStereoMode(config)
	})
}

func EndVrStereoMode() {
	Do(func() {
		rl.EndVrStereoMode()
	})
}

func LoadVrStereoConfig(device rl.VrDeviceInfo) rl.VrStereoConfig {
	return Call(func() rl.VrStereoConfig {
		return rl.LoadVrStereoConfig(device)
	})
}

func UnloadVrStereoConfig(config rl.VrStereoConfig) {
	Do(func() {
		rl.UnloadVrStereoConfig(config)
	})
}
