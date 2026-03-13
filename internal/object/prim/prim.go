package prim

import (
	"fmt"
	"image/color"
	"reflect"
	"runtime"
	"slices"
	"unsafe"

	"karalis/internal/collider"
	"karalis/internal/rlx"
	"karalis/pkg/lmath"

	pub_object "karalis/pkg/object"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Prim struct {
	mdl   rl.Model
	pos   rl.Vector3
	rot   rl.Vector3
	scale rl.Vector3
	color color.RGBA

	parent  pub_object.Object
	childs  []pub_object.Object
	col     pub_object.Collider
	cleaner *runtime.Cleanup
}

func (p *Prim) init() error {
	if p == nil {
		return fmt.Errorf("Invalid prim")
	}

	p.pos = rl.NewVector3(0, 0, 0)
	p.rot = rl.NewVector3(0, 0, 0)
	p.scale = rl.NewVector3(1, 1, 1)
	p.color = rl.White
	p.mdl = rl.Model{}
	if p.cleaner != nil {
		p.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(p, func(mdl rl.Model) {
		rlx.UnloadModel(mdl)
	}, p.mdl)
	p.cleaner = &cleaner

	p.parent = nil
	p.childs = []pub_object.Object{}

	col, err := collider.NewCollider(p)
	if err != nil {
		return err
	}
	p.col = col

	return nil
}

func (p *Prim) Prerender(cam pub_object.Camera) []func() {
	if p == nil {
		return []func(){}
	}

	return []func(){}
}

func (p *Prim) Render(cam pub_object.Camera) []func() {
	if p == nil {
		return []func(){}
	}

	matTransform := p.GetModelMatrix()
	rlx.DrawMesh(*p.mdl.Meshes, *p.mdl.Materials, matTransform)

	return []func(){}
}

func (p *Prim) Postrender(cam pub_object.Camera) []func() {
	if p == nil {
		return []func(){}
	}

	return []func(){}
}

func (p *Prim) OnResize(w int32, h int32) {
	if p == nil {
		return
	}
}

func (p *Prim) Update(dt float32) {
	if p == nil {
		return
	}

	if p.col != nil {
		p.col.Update(dt)
	}
}

func (p *Prim) GetModelMatrix() rl.Matrix {
	if p == nil {
		return rl.Matrix{}
	}

	matScale := rl.MatrixScale(p.scale.X, p.scale.Y, p.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(p.GetPitch()), float64(p.GetYaw()), float64(p.GetRoll()))
	matRotation := rl.QuaternionToMatrix(rl.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := rl.MatrixTranslate(p.pos.X, p.pos.Y, p.pos.Z)
	matTransform := rl.MatrixMultiply(rl.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = rl.MatrixMultiply(p.mdl.Transform, matTransform)
	return matTransform
}

func (p *Prim) GetModel() *rl.Model {
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

func (p *Prim) GetScale() rl.Vector3 {
	if p == nil {
		return rl.Vector3{}
	}

	return p.scale
}

func (p *Prim) SetScale(sc rl.Vector3) {
	if p == nil {
		return
	}

	p.scale = sc
}

func (p *Prim) GetPos() rl.Vector3 {
	if p == nil {
		return rl.Vector3{}
	}

	return p.pos
}

func (p *Prim) SetPos(pos rl.Vector3) {
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

func (p *Prim) GetVertices() []rl.Vector3 {
	if p == nil {
		return []rl.Vector3{}
	}

	verts := []rl.Vector3{}
	length := p.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(p.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, rl.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (p *Prim) GetUVs() []rl.Vector2 {
	if p == nil {
		return []rl.Vector2{}
	}

	uvs := []rl.Vector2{}
	length := p.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(p.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, rl.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (p *Prim) SetUVs(uvs []rl.Vector2) {
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

func (p *Prim) GetMaterials() *rl.Material {
	if p == nil {
		return nil
	}

	return p.mdl.Materials
}

func (p *Prim) SetTexture(tex rl.Texture2D) {
	rlx.SetMaterialTexture(p.mdl.Materials, rl.MapDiffuse, tex)
}

func (p *Prim) GetTexture() *rl.Texture2D {
	if p == nil {
		return nil
	}

	return &p.mdl.Materials.Maps.Texture
}

func (p *Prim) GetCollider() pub_object.Collider {
	if p == nil {
		return nil
	}

	return p.col
}

func (p *Prim) OnAdd(obj pub_object.Object) {
	if p == nil {
		return
	}
	p.parent = obj
}

func (p *Prim) OnRemove() {
	if p == nil {
		return
	}
	p.parent = nil
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

func (p *Prim) GetParent() pub_object.Object {
	if p == nil {
		return nil
	}
	return p.parent
}
