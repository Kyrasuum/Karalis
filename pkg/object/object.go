package object

import (
	"image/color"
	"unsafe"

	"karalis/internal/camera"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

/*
#cgo windows CFLAGS: -I../../raylib/src
#cgo windows LDFLAGS: -L../../raylib/src
#include "raylib.h"

void UpdateModelUVs(Model* mdl) {
	UpdateMeshBuffer(mdl->meshes[0], 1, &(mdl->meshes->texcoords[0]), mdl->meshes->vertexCount*2*sizeof(float), 0);
}
*/
import "C"

var ()

type Object interface {
	CanCollide() bool
	GetCollider() Collider
	GetModelMatrix() raylib.Matrix
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
	Collide(CollisionData)
	RegCollideHandler(func(CollisionData) bool)
	OnAdd()
	OnRemove()
	AddChild(obj Object)
	RemChild(obj Object)
	GetChilds() []Object
}

func UpdateModelUVs(mdl *raylib.Model) {
	C.UpdateModelUVs((*C.Model)(unsafe.Pointer(mdl)))
}
