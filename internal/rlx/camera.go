package rlx

import rl "github.com/gen2brain/raylib-go/raylib"

func GetCameraMatrix(camera rl.Camera) rl.Matrix {
	return Call(func() rl.Matrix {
		return rl.GetCameraMatrix(camera)
	})
}

func GetWorldToScreen(position rl.Vector3, camera rl.Camera) rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetWorldToScreen(position, camera)
	})
}

// If your raylib-go version exposes these:
func GetCameraProjectionMatrix(camera *rl.Camera, aspect float32) rl.Matrix {
	return Call(func() rl.Matrix {
		return rl.GetCameraProjectionMatrix(camera, aspect)
	})
}

func GetCameraViewMatrix(camera *rl.Camera) rl.Matrix {
	return Call(func() rl.Matrix {
		return rl.GetCameraViewMatrix(camera)
	})
}
