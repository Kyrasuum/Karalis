package object

import (
	"karalis/internal/camera"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Object interface {
	GetModelMatrix() raylib.Matrix
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
	SetTexture(mat *raylib.Material, tex raylib.Texture2D)
	GetTexture(mat *raylib.Material) raylib.Texture2D
	Prerender(cam *camera.Cam) []func()
	Render(cam *camera.Cam) []func()
	Postrender(cam *camera.Cam) []func()
	Update(dt float32)
	OnAdd()
	OnRemove()
	AddChild(obj Object)
	RemChild(obj Object)
}
