package object

import (
	"image/color"
	"math"
	"unsafe"

	"karalis/internal/camera"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

/*
#include "raylib.h"

void UpdateModelUVs(Model* mdl) {
	UpdateMeshBuffer(mdl->meshes[0], 1, &(mdl->meshes->texcoords[0]), mdl->meshes->vertexCount*2*sizeof(float), 0);
}
*/
import "C"

var ()

type Object interface {
	GetCollider() Collider
	GetModelMatrix() raylib.Matrix
	GetModel() *raylib.Model
	SetColor(c color.Color)
	GetColor() color.Color
	SetScale(s raylib.Vector3)
	GetScale() raylib.Vector3
	SetPos(p raylib.Vector3)
	GetPos() raylib.Vector3
	GetPitch() float32
	SetPitch(p float32)
	GetYaw() float32
	SetYaw(y float32)
	GetRoll() float32
	SetRoll(r float32)
	GetVertices() []raylib.Vector3
	GetUVs() []raylib.Vector2
	SetUVs(uvs []raylib.Vector2)
	GetMaterials() *raylib.Material
	SetTexture(tex raylib.Texture2D)
	GetTexture() raylib.Texture2D
	Prerender(cam *camera.Cam) []func()
	Render(cam *camera.Cam) []func()
	Postrender(cam *camera.Cam) []func()
	Update(dt float32)
	OnAdd()
	OnRemove()
	AddChild(obj Object)
	RemChild(obj Object)
	GetChilds() []Object
}

func UpdateModelUVs(mdl *raylib.Model) {
	C.UpdateModelUVs((*C.Model)(unsafe.Pointer(mdl)))
}

func meshVerticesXYZ(mesh *raylib.Mesh) []float32 {
	if mesh == nil || mesh.Vertices == nil || mesh.VertexCount <= 0 {
		return nil
	}
	n := int(mesh.VertexCount) * 3
	return unsafe.Slice(mesh.Vertices, n)
}

func meshIndices(mesh *raylib.Mesh) []uint16 {
	if mesh == nil || mesh.Indices == nil || mesh.TriangleCount <= 0 {
		return nil
	}
	n := int(mesh.TriangleCount) * 3
	return unsafe.Slice(mesh.Indices, n)
}

func meshHasIndices(mesh *raylib.Mesh) bool {
	return mesh != nil && mesh.Indices != nil
}

func meshTriangleCount(mesh *raylib.Mesh) int {
	if mesh == nil {
		return 0
	}
	// raylib already gives triangle count; trust it.
	return int(mesh.TriangleCount)
}

func meshPositionAt(mesh *raylib.Mesh, vertexIndex int) raylib.Vector3 {
	v := meshVerticesXYZ(mesh)
	base := vertexIndex * 3
	return raylib.Vector3{X: v[base], Y: v[base+1], Z: v[base+2]}
}

func meshTrianglePositions(mesh *raylib.Mesh, tri int) (a, b, c raylib.Vector3) {
	base := tri * 3
	if meshHasIndices(mesh) {
		idx := meshIndices(mesh)
		i0 := int(idx[base])
		i1 := int(idx[base+1])
		i2 := int(idx[base+2])
		return meshPositionAt(mesh, i0), meshPositionAt(mesh, i1), meshPositionAt(mesh, i2)
	}
	// Non-indexed triangle list: vertices laid out as triples
	return meshPositionAt(mesh, base), meshPositionAt(mesh, base+1), meshPositionAt(mesh, base+2)
}

// For "FLT_MAX" style comparisons.
func f32Inf() float32 { return float32(math.Inf(1)) }
