package ui

import (
	"image/color"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Button struct {
	pos    raylib.Rectangle
	inset  raylib.Vector2
	border float32

	hovered bool
	clicked bool

	text    string
	subtext string
	font    raylib.Font

	fontSize float32
	spacing  float32

	textCol color.RGBA
	bgCol   color.RGBA
	bordCol color.RGBA

	onPress func()
}

func (b *Button) defaults() {
	b.pos = raylib.Rectangle{0, 0, 0, 0}
	b.inset = raylib.Vector2{0, 0}
	b.border = float32(3)

	b.hovered = false
	b.clicked = false

	b.text = ""
	b.subtext = ""
	b.font = raylib.GetFontDefault()

	b.fontSize = 16
	b.spacing = 1

	b.textCol = color.RGBA{255, 255, 255, 255}
	b.bgCol = color.RGBA{0, 0, 0, 255}
	b.bordCol = color.RGBA{255, 255, 255, 255}

	b.onPress = func() {}
}

func (b *Button) Init() error {
	b.defaults()

	return nil
}

func (b *Button) SetOnPress(f func()) {
	b.onPress = f
}

func (b *Button) SetPosition(p raylib.Rectangle) {
	b.pos = p
	b.calcText()
}

func (b *Button) SetInset(i raylib.Vector2) {
	b.inset = i
	b.calcText()
}

func (b *Button) SetBorder(bo float32) {
	b.border = bo
	b.calcText()
}

func (b *Button) SetText(t string) {
	b.text = t
}

func (b *Button) SetFont(f raylib.Font) {
	b.font = f
}

func (b *Button) SetFontSize(f float32) {
	b.fontSize = f
}

func (b *Button) SetTextColor(c color.RGBA) {
	b.textCol = c
}

func (b *Button) SetBGColor(c color.RGBA) {
	b.bgCol = c
}

func (b *Button) SetBorderColor(c color.RGBA) {
	b.bordCol = c
}

func (b *Button) calcText() {
	b.subtext = b.text[:0]
	for i := 1; 1 < len(b.text); i++ {
		w := raylib.MeasureTextEx(b.font, b.text[:i], b.fontSize, b.spacing)
		if w.X > b.pos.Width-b.inset.X-b.border {
			break
		}
		b.subtext = b.text[:i]
	}
}

func (b *Button) Render() {
	raylib.DrawRectangleLinesEx(b.pos, b.border, b.bordCol)
	raylib.DrawRectangleRec(b.pos, b.bgCol)
	raylib.DrawTextEx(b.font, b.subtext, raylib.Vector2{(b.pos.X + b.inset.X + b.border), (b.pos.Y + b.inset.Y + b.border)}, b.fontSize, b.spacing, b.textCol)
}

func (b *Button) Update(dt float64) {
}

func (b *Button) OnInput() {
	mp := raylib.GetMousePosition()
	if raylib.CheckCollisionPointRec(mp, b.pos) {
		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			if !b.clicked {
				b.onPress()
			}
			b.clicked = true
		} else {
			b.clicked = false
		}
		b.hovered = true
	} else {
		b.clicked = false
		b.hovered = false
	}
}

func (b *Button) OnResize(w int32, h int32) {
}

func (b *Button) OnAdd() {
}

func (b *Button) OnRemove() {
}
