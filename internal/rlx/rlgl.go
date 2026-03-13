package rlx

import rl "github.com/gen2brain/raylib-go/raylib"

func DisableDepthMask() {
	Do(func() {
		rl.DisableDepthMask()
	})
}

func EnableDepthMask() {
	Do(func() {
		rl.EnableDepthMask()
	})
}

func DisableDepthTest() {
	Do(func() {
		rl.DisableDepthTest()
	})
}

func EnableDepthTest() {
	Do(func() {
		rl.EnableDepthTest()
	})
}

func PushMatrix() {
	Do(func() {
		rl.PushMatrix()
	})
}

func PopMatrix() {
	Do(func() {
		rl.PopMatrix()
	})
}

func Begin(mode int32) {
	Do(func() {
		rl.Begin(mode)
	})
}

func End() {
	Do(func() {
		rl.End()
	})
}

func EnableTextureCubemap(id uint32) {
	Do(func() {
		rl.EnableTextureCubemap(id)
	})
}

func DisableTextureCubemap() {
	Do(func() {
		rl.DisableTextureCubemap()
	})
}

func GetMatrixModelview() rl.Matrix {
	return Call(func() rl.Matrix {
		return rl.GetMatrixModelview()
	})
}

func GetMatrixProjection() rl.Matrix {
	return Call(func() rl.Matrix {
		return rl.GetMatrixProjection()
	})
}

func Color4ub(r uint8, g uint8, b uint8, a uint8) {
	Do(func() {
		rl.Color4ub(r, g, b, a)
	})
}

func TexCoord2f(x float32, y float32) {
	Do(func() {
		rl.TexCoord2f(x, y)
	})
}

func Vertex3f(x float32, y float32, z float32) {
	Do(func() {
		rl.Vertex3f(x, y, z)
	})
}
