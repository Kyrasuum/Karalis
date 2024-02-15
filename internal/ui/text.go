package ui

import (
	"image/color"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Text struct {
	pos raylib.Vector2

	text string
	font raylib.Font

	fontSize float32
	spacing  float32

	textCol color.RGBA
}

func (t *Text) defaults() {
	t.pos = raylib.Vector2{0, 0}

	t.text = ""
	t.font = raylib.GetFontDefault()

	t.fontSize = 16
	t.spacing = 1

	t.textCol = color.RGBA{255, 255, 255, 255}
}

func (t *Text) Init() error {
	t.defaults()
	return nil
}

func (t *Text) SetPosition(p raylib.Vector2) {
	t.pos = p
}

func (t *Text) SetText(s string) {
	t.text = s
}

func (t *Text) SetFont(f raylib.Font) {
	t.font = f
}

func (t *Text) SetFontSize(f float32) {
	t.fontSize = f
}

func (t *Text) SetTextColor(c color.RGBA) {
	t.textCol = c
}

func (t *Text) Render() {
	raylib.DrawTextEx(t.font, t.text, t.pos, t.fontSize, t.spacing, t.textCol)
}

func (t *Text) Update(dt float64) {
}

func (t *Text) OnInput() {
}

func (t *Text) OnAdd() {
}

func (t *Text) OnRemove() {
}

func (t *Text) OnResize(w int32, h int32) {
}
