package prim

import (
	"image/color"
	"reflect"
	"unsafe"

	"karalis/internal/camera"
	pub_object "karalis/pkg/object"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

var ()

type WireCube struct {
	mdl   raylib.Model
	pos   raylib.Vector3
	rot   raylib.Vector3
	scale raylib.Vector3
	color color.RGBA
}

func (c *WireCube) Init() error {
	c.pos = raylib.NewVector3(0, 0, 0)
	c.rot = raylib.NewVector3(0, 0, 0)
	c.scale = raylib.NewVector3(1, 1, 1)
	c.color = raylib.Red

	mdl, err := res.GetRes("mdl/cube.obj")
	if err != nil {
		return err
	}
	c.mdl = mdl.(raylib.Model)

	return nil
}

func (c *WireCube) GetModelMatrix() raylib.Matrix {
	matScale := raylib.MatrixScale(c.scale.X, c.scale.Y, c.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(c.GetPitch()), float64(c.GetYaw()), float64(c.GetRoll()))
	matRotation := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := raylib.MatrixTranslate(c.pos.X, c.pos.Y, c.pos.Z)
	matTransform := raylib.MatrixMultiply(raylib.MatrixMultiply(matScale, matRotation), matTranslation)
	return matTransform
}

func (c *WireCube) GetPos() raylib.Vector3 {
	return c.pos
}

func (c *WireCube) GetPitch() float32 {
	return c.rot.X
}

func (c *WireCube) SetPitch(p float32) {
	c.rot.X = p
}

func (c *WireCube) GetYaw() float32 {
	return c.rot.Y
}

func (c *WireCube) SetYaw(y float32) {
	c.rot.Y = y
}

func (c *WireCube) GetRoll() float32 {
	return c.rot.Z
}

func (c *WireCube) SetRoll(r float32) {
	c.rot.Z = r
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
	matTransform := c.GetModelMatrix()
	raylib.DrawMesh(*c.mdl.Meshes, *c.mdl.Materials, matTransform)
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
