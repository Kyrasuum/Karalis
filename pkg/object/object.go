package object

import (
	"image/color"
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

type Cell interface {
	Init() error
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

type Portal interface {
	Init(scene Cell, exit *Portal, cam *camera.Cam, obj Object) error
	GetScene() Cell
	Pair(portal Portal)
	GetPair() Portal
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
	GetCollider() Collider
	GetModelMatrix() raylib.Matrix
	GetModel() *raylib.Model
	GetVertices() []raylib.Vector3
	GetUVs() []raylib.Vector2
	SetUVs(uvs []raylib.Vector2)
	GetMaterials() *raylib.Material
	SetTexture(tex raylib.Texture2D)
	GetTexture() raylib.Texture2D
	SetPortal(obj Portal)
	GetPortal() Portal
	SetCam(cam *camera.Cam)
	GetCam() *camera.Cam
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

// If your axes came from a transform matrix columns, they may include scale.
// This normalizes axes and pushes scale into HalfExtents so later math is correct.
func OrientedBoxNormalizeScale(obb OrientedBox) OrientedBox {
	const eps = 1e-8

	// X
	lx := raylib.Vector3Length(obb.AxisX)
	if lx > eps {
		inv := 1.0 / lx
		obb.AxisX = raylib.Vector3Scale(obb.AxisX, float32(inv))
		obb.HalfExtents.X *= lx
	}

	// Y
	ly := raylib.Vector3Length(obb.AxisY)
	if ly > eps {
		inv := 1.0 / ly
		obb.AxisY = raylib.Vector3Scale(obb.AxisY, float32(inv))
		obb.HalfExtents.Y *= ly
	}

	// Z
	lz := raylib.Vector3Length(obb.AxisZ)
	if lz > eps {
		inv := 1.0 / lz
		obb.AxisZ = raylib.Vector3Scale(obb.AxisZ, float32(inv))
		obb.HalfExtents.Z *= lz
	}

	return obb
}
