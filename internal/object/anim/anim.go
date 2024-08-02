package anim

import (
	"fmt"
	"image/color"
	"reflect"
	"unsafe"

	"karalis/internal/camera"
	"karalis/pkg/app"
	pub_object "karalis/pkg/object"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

type Anim struct {
	mdl   raylib.Model
	pos   raylib.Vector3
	rot   raylib.Vector3
	scale raylib.Vector3
	color color.RGBA
}

func NewAnim(i interface{}) (a *Anim, err error) {
	a = &Anim{}

	err = a.Init()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Anim) Init() error {
	a.pos = raylib.NewVector3(0, 0, 0)
	a.rot = raylib.NewVector3(0, 0, 0)
	a.scale = raylib.NewVector3(1, 1, 1)
	a.color = raylib.White
	a.mdl = raylib.Model{}

	return nil
}

func (a *Anim) GetModelMatrix() raylib.Matrix {
	matScale := raylib.MatrixScale(a.scale.X, a.scale.Y, a.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(a.GetPitch()), float64(a.GetYaw()), float64(a.GetRoll()))
	matRotation := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := raylib.MatrixTranslate(a.pos.X, a.pos.Y, a.pos.Z)
	matTransform := raylib.MatrixMultiply(raylib.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = raylib.MatrixMultiply(a.mdl.Transform, matTransform)
	return matTransform
}

func (a *Anim) SetColor(c color.Color) {
	switch color := c.(type) {
	case color.RGBA:
		a.color = color
	}
}

func (a *Anim) GetColor() color.Color {
	return a.color
}

func (a *Anim) SetScale(s raylib.Vector3) {
	a.scale = s
}

func (a *Anim) GetScale() raylib.Vector3 {
	return a.scale
}

func (a *Anim) SetPos(p raylib.Vector3) {
	a.pos = p
}

func (a *Anim) GetPos() raylib.Vector3 {
	return a.pos
}

func (a *Anim) GetPitch() float32 {
	return a.rot.X
}

func (a *Anim) SetPitch(p float32) {
	a.rot.X = p
}

func (a *Anim) GetYaw() float32 {
	return a.rot.Y
}

func (a *Anim) SetYaw(y float32) {
	a.rot.Y = y
}

func (a *Anim) GetRoll() float32 {
	return a.rot.Z
}

func (a *Anim) SetRoll(r float32) {
	a.rot.Z = r
}

func (a *Anim) GetVertices() []raylib.Vector3 {
	verts := []raylib.Vector3{}
	length := a.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(a.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (a *Anim) GetUVs() []raylib.Vector2 {
	uvs := []raylib.Vector2{}
	length := a.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(a.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (a *Anim) SetUVs(uvs []raylib.Vector2) {
	length := int(a.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(a.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	pub_object.UpdateModelUVs(&a.mdl)
}

func (a *Anim) GetMaterials() *raylib.Material {
	return a.mdl.Materials
}

func (a *Anim) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	raylib.SetMaterialTexture(mat, raylib.MapDiffuse, tex)
}

func (a *Anim) GetTexture() raylib.Texture2D {
	return a.mdl.Materials.Maps.Texture
}

func (a *Anim) LoadModel(name string) error {
	mdl, err := res.GetRes(name)
	if err != nil {
		return err
	}
	switch model := mdl.(type) {
	case raylib.Model:
		a.mdl = model
	default:
		return fmt.Errorf("Invalid model object\n")
	}

	return nil
}

func (a *Anim) Prerender(cam *camera.Cam) []func() {
	return []func(){}
}

func (a *Anim) Render(cam *camera.Cam) []func() {
	matTransform := a.GetModelMatrix()
	sh := app.CurApp.GetShader()
	a.mdl.Materials.Shader = *sh.GetShader()
	raylib.DrawMesh(*a.mdl.Meshes, *a.mdl.Materials, matTransform)

	return []func(){}
}

func (a *Anim) Postrender(cam *camera.Cam) []func() {
	return []func(){}
}

func (a *Anim) Update(dt float32) {}

func (a *Anim) Collide(data pub_object.CollisionData) {}

func (a *Anim) RegCollideHandler(handler func(pub_object.CollisionData) bool) {}

func (a *Anim) GetCollidable() []pub_object.Object {
	return []pub_object.Object{}
}

func (a *Anim) GetCollider() pub_object.Collider {
	col := pub_object.Collider{}
	return col
}

func (a *Anim) OnAdd() {}

func (a *Anim) OnRemove() {}

func (a *Anim) AddChild(obj pub_object.Object) {}

func (a *Anim) RemChild(obj pub_object.Object) {}
