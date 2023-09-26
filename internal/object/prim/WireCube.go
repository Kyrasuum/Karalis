package prim

import (
	"image/color"

	pub_object "godev/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type WireCube struct {
	mdl   raylib.Model
	pos   raylib.Vector3
	size  float32
	color color.RGBA
}

func (c *WireCube) Init() {
	c.pos = raylib.NewVector3(0, 0, 0)
	c.size = 1
	c.color = raylib.Red
	c.mdl = raylib.LoadModel("res/prim/cube.obj")
}

func (c *WireCube) GetMaterials() *raylib.Material {
	return c.mdl.Materials
}

func (c *WireCube) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	raylib.SetMaterialTexture(mat, raylib.MapDiffuse, tex)
}

func (c *WireCube) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return mat.Maps.Texture
}

func (c *WireCube) Prerender() []func() {
	return []func(){}
}

func (c *WireCube) Render() []func() {
	raylib.SetTexture(c.mdl.Materials.Maps.Texture.ID)
	raylib.DrawModel(c.mdl, c.pos, c.size, c.color)
	raylib.SetTexture(0)
	return []func(){}
}

func (c *WireCube) Postrender() []func() {
	return []func(){}
}

func (c *WireCube) Update(dt float32) {
}

func (c *WireCube) OnAdd() {
}

func (c *WireCube) OnRemove() {
}

func (c *WireCube) AddChild(obj pub_object.Object) {
}

func (c *WireCube) RemChild(obj pub_object.Object) {
}
