package rlx

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func GetFontDefault() rl.Font {
	return Call(func() rl.Font {
		return rl.GetFontDefault()
	})
}

func LoadFont(fileName string) rl.Font {
	return Call(func() rl.Font {
		return rl.LoadFont(fileName)
	})
}

func LoadFontEx(fileName string, fontSize int32, codepoints []rune, runesNumber ...int32) rl.Font {
	cp := append([]rune(nil), codepoints...)
	return Call(func() rl.Font {
		return rl.LoadFontEx(fileName, fontSize, cp, runesNumber...)
	})
}

func LoadFontFromImage(image rl.Image, key color.RGBA, firstChar rune) rl.Font {
	return Call(func() rl.Font {
		return rl.LoadFontFromImage(image, key, firstChar)
	})
}

func IsFontValid(font rl.Font) bool {
	return Call(func() bool {
		return rl.IsFontValid(font)
	})
}

func UnloadFont(font rl.Font) {
	Do(func() {
		rl.UnloadFont(font)
	})
}

func DrawFPS(posX int32, posY int32) {
	Do(func() {
		rl.DrawFPS(posX, posY)
	})
}

func DrawText(text string, posX int32, posY int32, fontSize int32, col color.RGBA) {
	Do(func() {
		rl.DrawText(text, posX, posY, fontSize, col)
	})
}

func DrawTextEx(font rl.Font, text string, position rl.Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	Do(func() {
		rl.DrawTextEx(font, text, position, fontSize, spacing, tint)
	})
}

func DrawTextPro(font rl.Font, text string, position rl.Vector2, origin rl.Vector2, rotation float32, fontSize float32, spacing float32, tint color.RGBA) {
	Do(func() {
		rl.DrawTextPro(font, text, position, origin, rotation, fontSize, spacing, tint)
	})
}

func DrawTextCodepoint(font rl.Font, codepoint rune, position rl.Vector2, fontSize float32, tint color.RGBA) {
	Do(func() {
		rl.DrawTextCodepoint(font, codepoint, position, fontSize, tint)
	})
}

func DrawTextCodepoints(font rl.Font, codepoints []rune, position rl.Vector2, fontSize float32, spacing float32, tint color.RGBA) {
	cp := append([]rune(nil), codepoints...)
	Do(func() {
		rl.DrawTextCodepoints(font, cp, position, fontSize, spacing, tint)
	})
}

func SetTextLineSpacing(spacing int) {
	Do(func() {
		rl.SetTextLineSpacing(spacing)
	})
}

func MeasureText(text string, fontSize int32) int32 {
	return Call(func() int32 {
		return rl.MeasureText(text, fontSize)
	})
}

func MeasureTextEx(font rl.Font, text string, fontSize float32, spacing float32) rl.Vector2 {
	return Call(func() rl.Vector2 {
		return rl.MeasureTextEx(font, text, fontSize, spacing)
	})
}
