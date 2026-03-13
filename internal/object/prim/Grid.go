package prim

import (
	"fmt"
	"image/color"

	"karalis/internal/rlx"
	"karalis/pkg/app"

	pub_object "karalis/pkg/object"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Grid struct {
	spacing float32
	size    int32

	parent pub_object.Object
}

func NewGrid() (g *Grid, err error) {
	g = &Grid{}
	err = g.Init()

	return g, err
}

func (g *Grid) Init() error {
	if g == nil {
		return fmt.Errorf("Invalid grid")
	}

	g.parent = nil
	g.spacing = 1
	g.size = 10

	return nil
}

func (g *Grid) GetModelMatrix() rl.Matrix {
	if g == nil {
		return rl.Matrix{}
	}

	return rl.MatrixIdentity()
}

func (g *Grid) GetModel() *rl.Model {
	if g == nil {
		return nil
	}

	return nil
}

func (g *Grid) SetColor(col color.Color) {
	if g == nil {
		return
	}
}

func (g *Grid) GetColor() color.Color {
	if g == nil {
		return nil
	}

	return rl.White
}

func (g *Grid) SetScale(sc rl.Vector3) {
	if g == nil {
		return
	}
}

func (g *Grid) GetScale() rl.Vector3 {
	if g == nil {
		return rl.Vector3{}
	}

	return rl.NewVector3(1, 1, 1)
}

func (g *Grid) SetPos(p rl.Vector3) {
	if g == nil {
		return
	}
}

func (g *Grid) GetPos() rl.Vector3 {
	if g == nil {
		return rl.Vector3{}
	}

	return rl.NewVector3(0, 0, 0)
}

func (g *Grid) GetPitch() float32 {
	if g == nil {
		return 0
	}

	return 0
}

func (g *Grid) SetPitch(p float32) {
	if g == nil {
		return
	}
}

func (g *Grid) GetYaw() float32 {
	if g == nil {
		return 0
	}

	return 0
}

func (g *Grid) SetYaw(y float32) {
	if g == nil {
		return
	}
}

func (g *Grid) GetRoll() float32 {
	if g == nil {
		return 0
	}

	return 0
}

func (g *Grid) SetRoll(r float32) {
	if g == nil {
		return
	}
}

func (g *Grid) GetVertices() []rl.Vector3 {
	if g == nil {
		return []rl.Vector3{}
	}

	verts := []rl.Vector3{}
	return verts
}

func (g *Grid) GetUVs() []rl.Vector2 {
	if g == nil {
		return []rl.Vector2{}
	}

	uvs := []rl.Vector2{}
	return uvs
}

func (g *Grid) SetUVs(uvs []rl.Vector2) {
	if g == nil {
		return
	}
}

func (g *Grid) GetMaterials() *rl.Material {
	if g == nil {
		return nil
	}

	return &rl.Material{}
}

func (g *Grid) SetTexture(tex rl.Texture2D) {
	if g == nil {
		return
	}
}

func (g *Grid) GetTexture() *rl.Texture2D {
	if g == nil {
		return nil
	}

	return &rl.Texture2D{}
}

func (g *Grid) Prerender(cam pub_object.Camera) []func() {
	if g == nil {
		return []func(){}
	}

	return []func(){}
}

func (g *Grid) Render(cam pub_object.Camera) []func() {
	if g == nil {
		return []func(){}
	}

	sh := app.CurApp.GetShader()
	sh.Begin()
	rlx.DrawGrid(g.size, g.spacing)
	sh.End()
	return []func(){}
}

func (g *Grid) Postrender(cam pub_object.Camera) []func() {
	if g == nil {
		return []func(){}
	}

	return []func(){}
}

func (g *Grid) OnResize(w int32, h int32) {
	if g == nil {
		return
	}
}

func (g *Grid) Update(dt float32) {
	if g == nil {
		return
	}
}

func (g *Grid) GetCollider() pub_object.Collider {
	if g == nil {
		return nil
	}

	return nil
}

func (g *Grid) OnAdd(obj pub_object.Object) {
	if g == nil {
		return
	}
	g.parent = obj
}

func (g *Grid) OnRemove() {
	if g == nil {
		return
	}
	g.parent = nil
}

func (g *Grid) AddChild(obj pub_object.Object) {
	if g == nil {
		return
	}
}

func (g *Grid) RemChild(obj pub_object.Object) {
	if g == nil {
		return
	}
}

func (g *Grid) GetChilds() []pub_object.Object {
	if g == nil {
		return []pub_object.Object{}
	}

	return []pub_object.Object{}
}

func (g *Grid) GetParent() pub_object.Object {
	if g == nil {
		return nil
	}
	return g.parent
}
