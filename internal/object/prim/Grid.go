package prim

import (
	pub_object "godev/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Grid struct {
	pos  float32
	size int32
}

func (g *Grid) Init() {
	g.pos = 1
	g.size = 10
}

func (c *Grid) GetMaterials() *raylib.Material {
	return &raylib.Material{}
}

func (c *Grid) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
}

func (c *Grid) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return raylib.Texture2D{}
}

func (g *Grid) Prerender() []func() {
	return []func(){}
}

func (g *Grid) Render() []func() {
	raylib.DrawGrid(g.size, g.pos)
	return []func(){}
}

func (g *Grid) Postrender() []func() {
	return []func(){}
}

func (g *Grid) Update(dt float32) {
}

func (g *Grid) OnAdd() {
}

func (g *Grid) OnRemove() {
}

func (g *Grid) AddChild(obj pub_object.Object) {
}

func (g *Grid) RemChild(obj pub_object.Object) {
}
