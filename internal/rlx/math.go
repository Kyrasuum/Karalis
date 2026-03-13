package rlx

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Fade(col color.RGBA, alpha float32) color.RGBA {
	return rl.Fade(col, alpha)
}

func ColorToInt(col color.RGBA) int32 {
	return rl.ColorToInt(col)
}

func ColorNormalize(col color.RGBA) rl.Vector4 {
	return rl.ColorNormalize(col)
}

func ColorFromNormalized(normalized rl.Vector4) color.RGBA {
	return rl.ColorFromNormalized(normalized)
}

func ColorToHSV(col color.RGBA) rl.Vector3 {
	return rl.ColorToHSV(col)
}

func ColorFromHSV(hue float32, saturation float32, value float32) color.RGBA {
	return rl.ColorFromHSV(hue, saturation, value)
}

func ColorTint(col color.RGBA, tint color.RGBA) color.RGBA {
	return rl.ColorTint(col, tint)
}

func ColorBrightness(col color.RGBA, factor float32) color.RGBA {
	return rl.ColorBrightness(col, factor)
}

func ColorContrast(col color.RGBA, contrast float32) color.RGBA {
	return rl.ColorContrast(col, contrast)
}

func ColorAlpha(col color.RGBA, alpha float32) color.RGBA {
	return rl.ColorAlpha(col, alpha)
}

func ColorAlphaBlend(dst color.RGBA, src color.RGBA, tint color.RGBA) color.RGBA {
	return rl.ColorAlphaBlend(dst, src, tint)
}

func ColorLerp(col1, col2 color.RGBA, factor float32) color.RGBA {
	return rl.ColorLerp(col1, col2, factor)
}

func GetColor(hexValue uint) color.RGBA {
	return rl.GetColor(hexValue)
}

func GetPixelDataSize(width int32, height int32, format int32) int32 {
	return rl.GetPixelDataSize(width, height, format)
}
