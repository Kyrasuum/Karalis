package rlx

import rl "github.com/gen2brain/raylib-go/raylib"

func IsKeyPressed(key int32) bool {
	return Call(func() bool {
		return rl.IsKeyPressed(key)
	})
}

func IsKeyPressedRepeat(key int32) bool {
	return Call(func() bool {
		return rl.IsKeyPressedRepeat(key)
	})
}

func IsKeyDown(key int32) bool {
	return Call(func() bool {
		return rl.IsKeyDown(key)
	})
}

func IsKeyReleased(key int32) bool {
	return Call(func() bool {
		return rl.IsKeyReleased(key)
	})
}

func IsKeyUp(key int32) bool {
	return Call(func() bool {
		return rl.IsKeyUp(key)
	})
}

func GetKeyPressed() int32 {
	return Call(func() int32 {
		return rl.GetKeyPressed()
	})
}

func GetCharPressed() int32 {
	return Call(func() int32 {
		return rl.GetCharPressed()
	})
}

func SetExitKey(key int32) {
	Do(func() {
		rl.SetExitKey(key)
	})
}

func IsMouseButtonPressed(button rl.MouseButton) bool {
	return Call(func() bool {
		return rl.IsMouseButtonPressed(button)
	})
}

func IsMouseButtonDown(button rl.MouseButton) bool {
	return Call(func() bool {
		return rl.IsMouseButtonDown(button)
	})
}

func IsMouseButtonReleased(button rl.MouseButton) bool {
	return Call(func() bool {
		return rl.IsMouseButtonReleased(button)
	})
}

func IsMouseButtonUp(button rl.MouseButton) bool {
	return Call(func() bool {
		return rl.IsMouseButtonUp(button)
	})
}

func GetMouseX() int32 {
	return Call(func() int32 {
		return rl.GetMouseX()
	})
}

func GetMouseY() int32 {
	return Call(func() int32 {
		return rl.GetMouseY()
	})
}

func GetMousePosition() rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetMousePosition()
	})
}

func GetMouseDelta() rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetMouseDelta()
	})
}

func SetMousePosition(x int, y int) {
	Do(func() {
		rl.SetMousePosition(x, y)
	})
}

func SetMouseOffset(offsetX int, offsetY int) {
	Do(func() {
		rl.SetMouseOffset(offsetX, offsetY)
	})
}

func SetMouseScale(scaleX float32, scaleY float32) {
	Do(func() {
		rl.SetMouseScale(scaleX, scaleY)
	})
}

func GetMouseWheelMove() float32 {
	return Call(func() float32 {
		return rl.GetMouseWheelMove()
	})
}

func GetMouseWheelMoveV() rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetMouseWheelMoveV()
	})
}

func SetMouseCursor(cursor int32) {
	Do(func() {
		rl.SetMouseCursor(cursor)
	})
}

func GetTouchX() int32 {
	return Call(func() int32 {
		return rl.GetTouchX()
	})
}

func GetTouchY() int32 {
	return Call(func() int32 {
		return rl.GetTouchY()
	})
}

func GetTouchPosition(index int32) rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetTouchPosition(index)
	})
}

func GetTouchPointId(index int32) int32 {
	return Call(func() int32 {
		return rl.GetTouchPointId(index)
	})
}

func GetTouchPointCount() int32 {
	return Call(func() int32 {
		return rl.GetTouchPointCount()
	})
}

func SetGesturesEnabled(flags uint32) {
	Do(func() {
		rl.SetGesturesEnabled(flags)
	})
}

func IsGestureDetected(gesture rl.Gestures) bool {
	return Call(func() bool {
		return rl.IsGestureDetected(gesture)
	})
}

func GetGestureDetected() rl.Gestures {
	return Call(func() rl.Gestures {
		return rl.GetGestureDetected()
	})
}

func GetGestureHoldDuration() float32 {
	return Call(func() float32 {
		return rl.GetGestureHoldDuration()
	})
}

func GetGestureDragVector() rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetGestureDragVector()
	})
}

func GetGestureDragAngle() float32 {
	return Call(func() float32 {
		return rl.GetGestureDragAngle()
	})
}

func GetGesturePinchVector() rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.GetGesturePinchVector()
	})
}

func GetGesturePinchAngle() float32 {
	return Call(func() float32 {
		return rl.GetGesturePinchAngle()
	})
}
