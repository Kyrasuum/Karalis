package ui

import (
	"image/color"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Box struct {
	pos    raylib.Rectangle
	inset  raylib.Vector2
	border float32

	bgCol   color.RGBA
	bordCol color.RGBA
}

func (b *Box) defaults() {
	b.pos = raylib.Rectangle{0, 0, 0, 0}
	b.inset = raylib.Vector2{0, 0}
	b.border = float32(3)

	b.bgCol = color.RGBA{0, 0, 0, 255}
	b.bordCol = color.RGBA{255, 255, 255, 255}
}

func (b *Box) Init() error {
	if b == nil {
		return fmt.Errorf("Invalid box")
	}

	b.defaults()

	return nil
}

func (b *Box) SetPosition(p raylib.Rectangle) {
	if b == nil {
		return
	}

	b.pos = p
}

func (b *Box) SetInset(i raylib.Vector2) {
	if b == nil {
		return
	}

	b.inset = i
}

func (b *Box) SetBorder(bo float32) {
	if b == nil {
		return
	}

	b.border = bo
}

func (b *Box) SetBGColor(c color.RGBA) {
	if b == nil {
		return
	}

	b.bgCol = c
}

func (b *Box) SetBorderColor(c color.RGBA) {
	if b == nil {
		return
	}

	b.bordCol = c
}

func (b *Box) Render() {
	raylib.DrawRectangleLinesEx(b.pos, b.border, b.bordCol)
	if b == nil {
		return
	}

	raylib.DrawRectangleRec(b.pos, b.bgCol)
}

func (b *Box) Update(dt float64) {
	if b == nil {
		return
	}
}

func (b *Box) OnInput() {
	if b == nil {
		return
	}
}

func (b *Box) OnResize(w int32, h int32) {
	if b == nil {
		return
	}
}

func (b *Box) OnAdd() {
	if b == nil {
		return
	}
}

func (b *Box) OnRemove() {
	if b == nil {
		return
	}
}
