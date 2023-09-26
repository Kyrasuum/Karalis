package character

import (
	"godev/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Character struct {
}

func (c *Character) Init() {
}

func (c *Character) Prerender() []func() {
	return []func(){}
}

func (c *Character) Render() []func() {
	return []func(){}
}

func (c *Character) Postrender() []func() {
	return []func(){}
}

func (c *Character) Update(dt float32) {
}

func (c *Character) OnAdd() {
}

func (c *Character) OnRemove() {
}

func (c *Character) AddChild(obj object.Object) {
}

func (c *Character) RemChild(obj object.Object) {
}

func (c *Character) GetMaterials() *raylib.Material {
	return &raylib.Material{}
}

func (c *Character) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
}

func (c *Character) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return raylib.Texture2D{}
}
