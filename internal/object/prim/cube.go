package prim

import (
	"image/color"
	"reflect"
	"unsafe"

	"karalis/internal/camera"
	pub_object "karalis/pkg/object"

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

func (c *Cube) GetVertices() []raylib.Vector3 {
	verts := []raylib.Vector3{}
	length := c.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (c *Cube) GetUVs() []raylib.Vector2 {
	uvs := []raylib.Vector2{}
	length := c.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (c *Cube) SetUVs(uvs []raylib.Vector2) {
	var mdluvs []float32
	for i := 0; i < len(uvs); i++ {
		mdluvs = append(mdluvs, uvs[i].X, uvs[i].Y)
	}
	c.mdl.Meshes.Texcoords = &mdluvs[0]
}

func (c *Cube) GetModelMatrix() raylib.Matrix {
	return c.mdl.Transform
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

func (c *Cube) Prerender(cam *camera.Cam) []func() {
	return []func(){}
}

func (c *Cube) Render(cam *camera.Cam) []func() {
	raylib.SetTexture(c.mdl.Materials.Maps.Texture.ID)
	raylib.DrawModel(c.mdl, c.pos, c.size, c.color)
	raylib.SetTexture(0)
	return []func(){}
}

func (c *Cube) Postrender(cam *camera.Cam) []func() {
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
