package prim

/*
#include "raylib.h"

void UpdateModelUVs(Model* mdl) {
	UpdateMeshBuffer(mdl->meshes[0], 1, &(mdl->meshes->texcoords[0]), mdl->meshes->vertexCount*2*sizeof(float), 0);
}
*/
import "C"
import (
	"image/color"
	"reflect"
	"unsafe"

	"karalis/internal/camera"
	pub_object "karalis/pkg/object"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Square struct {
	mdl   raylib.Model
	pos   raylib.Vector3
	size  float32
	color color.RGBA
}

func (c *Square) Init() error {
	c.pos = raylib.NewVector3(0, 0, 0)
	c.size = 1
	c.color = raylib.White

	mdl, err := res.GetRes("mdl/square.obj")
	if err != nil {
		return err
	}
	c.mdl = mdl.(raylib.Model)

	return nil
}

func (c *Square) GetVertices() []raylib.Vector3 {
	verts := []raylib.Vector3{}
	length := int(c.mdl.Meshes.VertexCount)

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Vertices))
	header.Len = length * 3
	header.Cap = length * 3

	for i := 0; i < length; i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (c *Square) GetUVs() []raylib.Vector2 {
	uvs := []raylib.Vector2{}
	length := int(c.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < length; i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (c *Square) SetUVs(uvs []raylib.Vector2) {
	length := int(c.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	C.UpdateModelUVs((*C.Model)(unsafe.Pointer(&c.mdl)))
}

func (c *Square) GetModelMatrix() raylib.Matrix {
	return c.mdl.Transform
}

func (c *Square) GetMaterials() *raylib.Material {
	return c.mdl.Materials
}

func (c *Square) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	raylib.SetMaterialTexture(mat, raylib.MapDiffuse, tex)
}

func (c *Square) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return mat.Maps.Texture
}

func (c *Square) Prerender(cam *camera.Cam) []func() {
	return []func(){}
}

func (c *Square) Render(cam *camera.Cam) []func() {
	raylib.SetTexture(c.mdl.Materials.Maps.Texture.ID)
	raylib.DrawModel(c.mdl, c.pos, c.size, c.color)
	raylib.SetTexture(0)
	return []func(){}
}

func (c *Square) Postrender(cam *camera.Cam) []func() {
	return []func(){}
}

func (c *Square) Update(dt float32) {
}

func (c *Square) OnAdd() {
}

func (c *Square) OnRemove() {
}

func (c *Square) AddChild(obj pub_object.Object) {
}

func (c *Square) RemChild(obj pub_object.Object) {
}
