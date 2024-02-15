package character

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Object interface {
	Init() error
	GetMaterials() *raylib.Material
	SetTexture(mat *raylib.Material, tex raylib.Texture2D)
	GetTexture(mat *raylib.Material) raylib.Texture2D
	Prerender() []func()
	Render() []func()
	Postrender() []func()
	Update(dt float32)
	OnAdd()
	OnRemove()
	AddChild(obj Object)
	RemChild(obj Object)
}
