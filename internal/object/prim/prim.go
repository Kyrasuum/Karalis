package prim

import (
	"fmt"
	"image/color"
	"reflect"
	"slices"
	"unsafe"

	"karalis/internal/camera"
	"karalis/internal/collider"
	"karalis/pkg/app"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

type Prim struct {
	mdl   raylib.Model
	pos   raylib.Vector3
	rot   raylib.Vector3
	scale raylib.Vector3
	color color.RGBA

	childs []pub_object.Object
	col    pub_object.Collider
}

func (p *Prim) init() error {
	if p == nil {
		return fmt.Errorf("Invalid prim")
	}

	p.pos = raylib.NewVector3(0, 0, 0)
	p.rot = raylib.NewVector3(0, 0, 0)
	p.scale = raylib.NewVector3(1, 1, 1)
	p.color = raylib.White
	p.mdl = raylib.Model{}

	p.childs = []pub_object.Object{}

	col, err := collider.NewCollider(p)
	if err != nil {
		return err
	}
	p.col = col

	return nil
}

func (p *Prim) GetModelMatrix() raylib.Matrix {
	if p == nil {
		return raylib.Matrix{}
	}

	matScale := raylib.MatrixScale(p.scale.X, p.scale.Y, p.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(p.GetPitch()), float64(p.GetYaw()), float64(p.GetRoll()))
	matRotation := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := raylib.MatrixTranslate(p.pos.X, p.pos.Y, p.pos.Z)
	matTransform := raylib.MatrixMultiply(raylib.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = raylib.MatrixMultiply(p.mdl.Transform, matTransform)
	return matTransform
}

func (p *Prim) GetModel() *raylib.Model {
	if p == nil {
		return nil
	}

	return &p.mdl
}

func (p *Prim) GetColor() color.Color {
	if p == nil {
		return nil
	}

	return p.color
}

func (p *Prim) SetColor(col color.Color) {
	if p == nil {
		return
	}

	switch color := col.(type) {
	case color.RGBA:
		p.color = color
	}
}

func (p *Prim) GetScale() raylib.Vector3 {
	if p == nil {
		return raylib.Vector3{}
	}

	return p.scale
}

func (p *Prim) SetScale(sc raylib.Vector3) {
	if p == nil {
		return
	}

	p.scale = sc
}

func (p *Prim) GetPos() raylib.Vector3 {
	if p == nil {
		return raylib.Vector3{}
	}

	return p.pos
}

func (p *Prim) SetPos(pos raylib.Vector3) {
	if p == nil {
		return
	}

	p.pos = pos
}

func (p *Prim) GetPitch() float32 {
	if p == nil {
		return 0
	}

	return p.rot.X
}

func (p *Prim) SetPitch(pitch float32) {
	if p == nil {
		return
	}

	p.rot.X = pitch
}

func (p *Prim) GetYaw() float32 {
	if p == nil {
		return 0
	}

	return p.rot.Y
}

func (p *Prim) SetYaw(yaw float32) {
	if p == nil {
		return
	}

	p.rot.Y = yaw
}

func (p *Prim) GetRoll() float32 {
	if p == nil {
		return 0
	}

	return p.rot.Z
}

func (p *Prim) SetRoll(roll float32) {
	if p == nil {
		return
	}

	p.rot.Z = roll
}

func (p *Prim) GetVertices() []raylib.Vector3 {
	if p == nil {
		return []raylib.Vector3{}
	}

	verts := []raylib.Vector3{}
	length := p.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(p.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (p *Prim) GetUVs() []raylib.Vector2 {
	if p == nil {
		return []raylib.Vector2{}
	}

	uvs := []raylib.Vector2{}
	length := p.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(p.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (p *Prim) SetUVs(uvs []raylib.Vector2) {
	if p == nil {
		return
	}

	length := int(p.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(p.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	pub_object.UpdateModelUVs(&p.mdl)
}

func (p *Prim) GetMaterials() *raylib.Material {
	if p == nil {
		return nil
	}

	return p.mdl.Materials
}

func (p *Prim) SetTexture(tex raylib.Texture2D) {
	raylib.SetMaterialTexture(p.mdl.Materials, raylib.MapDiffuse, tex)
}

func (p *Prim) GetTexture() *raylib.Texture2D {
	if p == nil {
		return nil
	}

	return &p.mdl.Materials.Maps.Texture
}

func (p *Prim) Prerender(cam *camera.Cam) []func() {
	if p == nil {
		return []func(){}
	}

	return []func(){}
}

func (p *Prim) Render(cam *camera.Cam) []func() {
	if p == nil {
		return []func(){}
	}

	matTransform := p.GetModelMatrix()
	sh := app.CurApp.GetShader()
	p.mdl.Materials.Shader = *sh.GetShader()
	raylib.DrawMesh(*p.mdl.Meshes, *p.mdl.Materials, matTransform)

	return []func(){}
}

func (p *Prim) Postrender(cam *camera.Cam) []func() {
	if p == nil {
		return []func(){}
	}

	return []func(){}
}

func (p *Prim) Update(dt float32) {
	if p == nil {
		return
	}

	if p.col != nil {
		p.col.Update(dt)
	}
}

func (p *Prim) GetCollider() pub_object.Collider {
	if p == nil {
		return nil
	}

	return p.col
}

func (p *Prim) OnAdd() {
	if p == nil {
		return
	}
}

func (p *Prim) OnRemove() {
	if p == nil {
		return
	}
}

func (p *Prim) AddChild(obj pub_object.Object) {
	if p == nil {
		return
	}
}

func (p *Prim) RemChild(obj pub_object.Object) {
	if p == nil {
		return
	}
}

func (p *Prim) GetChilds() []pub_object.Object {
	if p == nil {
		return []pub_object.Object{}
	}

	childs := p.childs
	grandchilds := []pub_object.Object{}
	for _, child := range childs {
		grandchilds = append(grandchilds, child.GetChilds()...)
	}

	return slices.Concat(grandchilds, childs)
}
