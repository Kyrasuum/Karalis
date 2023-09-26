package prim

import (
	"image/color"

	pub_object "godev/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Cube struct {
	mdl   raylib.Model
	pos   raylib.Vector3
	size  float32
	color color.RGBA
}

func (c *Cube) Init() {
	c.pos = raylib.NewVector3(0, 0, 0)
	c.size = 1
	c.color = raylib.White

	c.mdl = raylib.LoadModel("res/prim/cube.obj")
}

func (c *Cube) GetMaterials() *raylib.Material {
	return c.mdl.Materials
}

func (c *Cube) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	raylib.SetMaterialTexture(mat, raylib.MapDiffuse, tex)
}

func (c *Cube) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return mat.Maps.Texture
}

func (c *Cube) Prerender() []func() {
	return []func(){}
}

func (c *Cube) Render() []func() {
	raylib.SetTexture(c.mdl.Materials.Maps.Texture.ID)
	raylib.DrawModel(c.mdl, c.pos, c.size, c.color)
	raylib.SetTexture(0)
	return []func(){}
}

func (c *Cube) Postrender() []func() {
	return []func(){}
}

func (c *Cube) Update(dt float32) {
}

func (c *Cube) OnAdd() {
}

func (c *Cube) OnRemove() {
}

func (c *Cube) AddChild(obj pub_object.Object) {
}

func (c *Cube) RemChild(obj pub_object.Object) {
}
