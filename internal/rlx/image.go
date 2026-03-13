package rlx

import (
	"image"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadImageFromTexture(texture rl.Texture2D) *rl.Image {
	return Call(func() *rl.Image {
		return rl.LoadImageFromTexture(texture)
	})
}

func LoadImageFromScreen() *rl.Image {
	return Call(func() *rl.Image {
		return rl.LoadImageFromScreen()
	})
}

func LoadImage(fileName string) *rl.Image {
	return rl.LoadImage(fileName)
}

func LoadImageRaw(fileName string, width int32, height int32, format rl.PixelFormat, headerSize int32) *rl.Image {
	return rl.LoadImageRaw(fileName, width, height, format, headerSize)
}

func LoadImageAnim(fileName string, frames *int32) *rl.Image {
	return rl.LoadImageAnim(fileName, frames)
}

func LoadImageAnimFromMemory(fileType string, fileData []byte, dataSize int32, frames *int32) *rl.Image {
	buf := append([]byte(nil), fileData...)
	return rl.LoadImageAnimFromMemory(fileType, buf, dataSize, frames)
}

func LoadImageFromMemory(fileType string, fileData []byte, dataSize int32) *rl.Image {
	buf := append([]byte(nil), fileData...)
	return rl.LoadImageFromMemory(fileType, buf, dataSize)
}

func IsImageValid(image *rl.Image) bool {
	return rl.IsImageValid(image)
}

func UnloadImage(image *rl.Image) {
	rl.UnloadImage(image)
}

func ExportImage(image rl.Image, fileName string) bool {
	return rl.ExportImage(image, fileName)
}

func ExportImageToMemory(image rl.Image, fileType string) []byte {
	return rl.ExportImageToMemory(image, fileType)
}

func GenImageColor(width int, height int, col color.RGBA) *rl.Image {
	return rl.GenImageColor(width, height, col)
}

func ImageCopy(image *rl.Image) *rl.Image {
	return rl.ImageCopy(image)
}

func ImageFormat(image *rl.Image, newFormat rl.PixelFormat) {
	rl.ImageFormat(image, newFormat)
}

func ImageCrop(image *rl.Image, crop rl.Rectangle) {
	rl.ImageCrop(image, crop)
}

func ImageResize(image *rl.Image, newWidth int32, newHeight int32) {
	rl.ImageResize(image, newWidth, newHeight)
}

func ImageResizeNN(image *rl.Image, newWidth int32, newHeight int32) {
	rl.ImageResizeNN(image, newWidth, newHeight)
}

func ImageFlipVertical(image *rl.Image) {
	rl.ImageFlipVertical(image)
}

func ImageFlipHorizontal(image *rl.Image) {
	rl.ImageFlipHorizontal(image)
}

func ImageRotate(image *rl.Image, degrees int32) {
	rl.ImageRotate(image, degrees)
}

func ImageColorTint(image *rl.Image, col color.RGBA) {
	rl.ImageColorTint(image, col)
}

func LoadImageColors(image *rl.Image) []color.RGBA {
	return rl.LoadImageColors(image)
}

func UnloadImageColors(colors []color.RGBA) {
	rl.UnloadImageColors(colors)
}

func NewImageFromImage(img image.Image) *rl.Image {
	return rl.NewImageFromImage(img)
}
