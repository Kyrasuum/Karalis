package prim

import (
	"fmt"
	"image/color"

	"karalis/internal/camera"
	"karalis/pkg/app"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Grid struct {
	spacing float32
	size    int32
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

	g.spacing = 1
	g.size = 10

	return nil
}

func (g *Grid) GetModelMatrix() raylib.Matrix {
	if g == nil {
		return raylib.Matrix{}
	}

	return raylib.MatrixIdentity()
}

func (g *Grid) GetModel() *raylib.Model {
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

	return raylib.White
}

func (g *Grid) SetScale(sc raylib.Vector3) {
	if g == nil {
		return
	}
}

func (g *Grid) GetScale() raylib.Vector3 {
	if g == nil {
		return raylib.Vector3{}
	}

	return raylib.NewVector3(1, 1, 1)
}

func (g *Grid) SetPos(p raylib.Vector3) {
	if g == nil {
		return
	}
}

func (g *Grid) GetPos() raylib.Vector3 {
	if g == nil {
		return raylib.Vector3{}
	}

	return raylib.NewVector3(0, 0, 0)
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

func (g *Grid) GetVertices() []raylib.Vector3 {
	if g == nil {
		return []raylib.Vector3{}
	}

	verts := []raylib.Vector3{}
	return verts
}

func (g *Grid) GetUVs() []raylib.Vector2 {
	if g == nil {
		return []raylib.Vector2{}
	}

	uvs := []raylib.Vector2{}
	return uvs
}

func (g *Grid) SetUVs(uvs []raylib.Vector2) {
	if g == nil {
		return
	}
}

func (g *Grid) GetMaterials() *raylib.Material {
	if g == nil {
		return nil
	}

	return &raylib.Material{}
}

func (g *Grid) SetTexture(tex raylib.Texture2D) {
	if g == nil {
		return
	}
}

func (g *Grid) GetTexture() raylib.Texture2D {
	if g == nil {
		return raylib.Texture2D{}
	}

	return raylib.Texture2D{}
}

func (g *Grid) Prerender(cam *camera.Cam) []func() {
	if g == nil {
		return []func(){}
	}

	return []func(){}
}

func (g *Grid) Render(cam *camera.Cam) []func() {
	if g == nil {
		return []func(){}
	}

	sh := app.CurApp.GetShader()
	sh.Begin()
	raylib.DrawGrid(g.size, g.spacing)
	sh.End()
	return []func(){}
}

func (g *Grid) Postrender(cam *camera.Cam) []func() {
	if g == nil {
		return []func(){}
	}

	return []func(){}
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

func (g *Grid) OnAdd() {
	if g == nil {
		return
	}
}

func (g *Grid) OnRemove() {
	if g == nil {
		return
	}
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
