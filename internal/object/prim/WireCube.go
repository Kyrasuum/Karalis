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

func (c *WireCube) GetVertices() []raylib.Vector3 {
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

func (c *WireCube) GetUVs() []raylib.Vector2 {
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

func (c *WireCube) SetUVs(uvs []raylib.Vector2) {
	var mdluvs []float32
	for i := 0; i < len(uvs); i++ {
		mdluvs = append(mdluvs, uvs[i].X, uvs[i].Y)
	}
	c.mdl.Meshes.Texcoords = &mdluvs[0]
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

func (c *WireCube) Prerender(cam *camera.Cam) []func() {
	return []func(){}
}

func (c *WireCube) Render(cam *camera.Cam) []func() {
	raylib.SetTexture(c.mdl.Materials.Maps.Texture.ID)
	raylib.DrawModel(c.mdl, c.pos, c.size, c.color)
	raylib.SetTexture(0)
	return []func(){}
}

func (c *WireCube) Postrender(cam *camera.Cam) []func() {
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
