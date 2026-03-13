package rlx

import (
	"image"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTexture(fileName string) rl.Texture2D {
	return Call(func() rl.Texture2D {
		return rl.LoadTexture(fileName)
	})
}

func LoadTextureFromImage(image *rl.Image) rl.Texture2D {
	return Call(func() rl.Texture2D {
		return rl.LoadTextureFromImage(image)
	})
}

func LoadTextureCubemap(image *rl.Image, layout int32) rl.Texture2D {
	return Call(func() rl.Texture2D {
		return rl.LoadTextureCubemap(image, layout)
	})
}

func LoadRenderTexture(width int32, height int32) rl.RenderTexture2D {
	return Call(func() rl.RenderTexture2D {
		return rl.LoadRenderTexture(width, height)
	})
}

func IsTextureValid(texture rl.Texture2D) bool {
	return Call(func() bool {
		return rl.IsTextureValid(texture)
	})
}

func UnloadTexture(texture rl.Texture2D) {
	Do(func() {
		rl.UnloadTexture(texture)
	})
}

func IsRenderTextureValid(target rl.RenderTexture2D) bool {
	return Call(func() bool {
		return rl.IsRenderTextureValid(target)
	})
}

func UnloadRenderTexture(target rl.RenderTexture2D) {
	Do(func() {
		rl.UnloadRenderTexture(target)
	})
}

func UpdateTexture(texture rl.Texture2D, pixels any) {
	switch p := pixels.(type) {
	case []color.RGBA:
		buf := append([]color.RGBA(nil), p...)
		Do(func() {
			rl.UpdateTexture(texture, buf)
		})
	case *image.RGBA:
		if p == nil {
			return
		}
		buf := image.NewRGBA(p.Rect)
		copy(buf.Pix, p.Pix)
		Do(func() {
			rl.UpdateTexture(texture, buf)
		})
	case []byte:
		buf := append([]byte(nil), p...)
		Do(func() {
			rl.UpdateTexture(texture, buf)
		})
	default:
		panic("rlx.UpdateTexture: unsupported pixels type")
	}
}

func UpdateTextureRec(texture rl.Texture2D, rec rl.Rectangle, pixels any) {
	switch p := pixels.(type) {
	case []color.RGBA:
		buf := append([]color.RGBA(nil), p...)
		Do(func() {
			rl.UpdateTextureRec(texture, rec, buf)
		})
	case *image.RGBA:
		if p == nil {
			return
		}
		buf := image.NewRGBA(p.Rect)
		copy(buf.Pix, p.Pix)
		Do(func() {
			rl.UpdateTextureRec(texture, rec, buf)
		})
	case []byte:
		buf := append([]byte(nil), p...)
		Do(func() {
			rl.UpdateTextureRec(texture, rec, buf)
		})
	default:
		panic("rlx.UpdateTextureRec: unsupported pixels type")
	}
}

func GenTextureMipmaps(texture *rl.Texture2D) {
	Do(func() {
		rl.GenTextureMipmaps(texture)
	})
}

func SetTextureFilter(texture rl.Texture2D, filter rl.TextureFilterMode) {
	Do(func() {
		rl.SetTextureFilter(texture, filter)
	})
}

func SetTextureWrap(texture rl.Texture2D, wrap rl.TextureWrapMode) {
	Do(func() {
		rl.SetTextureWrap(texture, wrap)
	})
}

func DrawTexture(texture rl.Texture2D, posX int32, posY int32, tint color.RGBA) {
	Do(func() {
		rl.DrawTexture(texture, posX, posY, tint)
	})
}

func DrawTextureV(texture rl.Texture2D, position rl.Vector2, tint color.RGBA) {
	Do(func() {
		rl.DrawTextureV(texture, position, tint)
	})
}

func DrawTextureEx(texture rl.Texture2D, position rl.Vector2, rotation float32, scale float32, tint color.RGBA) {
	Do(func() {
		rl.DrawTextureEx(texture, position, rotation, scale, tint)
	})
}

func DrawTextureRec(texture rl.Texture2D, source rl.Rectangle, position rl.Vector2, tint color.RGBA) {
	Do(func() {
		rl.DrawTextureRec(texture, source, position, tint)
	})
}

func DrawTexturePro(texture rl.Texture2D, source rl.Rectangle, dest rl.Rectangle, origin rl.Vector2, rotation float32, tint color.RGBA) {
	Do(func() {
		rl.DrawTexturePro(texture, source, dest, origin, rotation, tint)
	})
}

func DrawTextureNPatch(texture rl.Texture2D, nPatchInfo rl.NPatchInfo, dest rl.Rectangle, origin rl.Vector2, rotation float32, tint color.RGBA) {
	Do(func() {
		rl.DrawTextureNPatch(texture, nPatchInfo, dest, origin, rotation, tint)
	})
}
