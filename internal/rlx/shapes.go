package rlx

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func SetShapesTexture(texture rl.Texture2D, source rl.Rectangle) {
	Do(func() {
		rl.SetShapesTexture(texture, source)
	})
}

func GetShapesTexture() rl.Texture2D {
	return Call(func() rl.Texture2D {
		return rl.GetShapesTexture()
	})
}

func GetShapesTextureRectangle() rl.Rectangle {
	return Call(func() rl.Rectangle {
		return rl.GetShapesTextureRectangle()
	})
}

func DrawPixel(posX int32, posY int32, col color.RGBA) {
	Do(func() {
		rl.DrawPixel(posX, posY, col)
	})
}

func DrawPixelV(position rl.Vector2, col color.RGBA) {
	Do(func() {
		rl.DrawPixelV(position, col)
	})
}

func DrawLine(startPosX int32, startPosY int32, endPosX int32, endPosY int32, col color.RGBA) {
	Do(func() {
		rl.DrawLine(startPosX, startPosY, endPosX, endPosY, col)
	})
}

func DrawLineV(startPos rl.Vector2, endPos rl.Vector2, col color.RGBA) {
	Do(func() {
		rl.DrawLineV(startPos, endPos, col)
	})
}

func DrawLineEx(startPos rl.Vector2, endPos rl.Vector2, thick float32, col color.RGBA) {
	Do(func() {
		rl.DrawLineEx(startPos, endPos, thick, col)
	})
}

func DrawLineStrip(points []rl.Vector2, col color.RGBA) {
	buf := append([]rl.Vector2(nil), points...)
	Do(func() {
		rl.DrawLineStrip(buf, col)
	})
}

func DrawCircle(centerX int32, centerY int32, radius float32, col color.RGBA) {
	Do(func() {
		rl.DrawCircle(centerX, centerY, radius, col)
	})
}

func DrawCircleV(center rl.Vector2, radius float32, col color.RGBA) {
	Do(func() {
		rl.DrawCircleV(center, radius, col)
	})
}

func DrawRectangle(posX int32, posY int32, width int32, height int32, col color.RGBA) {
	Do(func() {
		rl.DrawRectangle(posX, posY, width, height, col)
	})
}

func DrawRectangleRec(rec rl.Rectangle, col color.RGBA) {
	Do(func() {
		rl.DrawRectangleRec(rec, col)
	})
}

func DrawRectanglePro(rec rl.Rectangle, origin rl.Vector2, rotation float32, col color.RGBA) {
	Do(func() {
		rl.DrawRectanglePro(rec, origin, rotation, col)
	})
}

func DrawTriangle(v1 rl.Vector2, v2 rl.Vector2, v3 rl.Vector2, col color.RGBA) {
	Do(func() {
		rl.DrawTriangle(v1, v2, v3, col)
	})
}
