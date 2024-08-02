package prim

import (
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
	g.spacing = 1
	g.size = 10

	return nil
}

func (c *Grid) GetModelMatrix() raylib.Matrix {
	return raylib.MatrixIdentity()
}

func (g *Grid) SetColor(col color.Color) {
}

func (g *Grid) GetColor() color.Color {
	return raylib.White
}

func (g *Grid) SetScale(sc raylib.Vector3) {
}

func (g *Grid) GetScale() raylib.Vector3 {
	return raylib.NewVector3(1, 1, 1)
}

func (g *Grid) SetPos(p raylib.Vector3) {
}

func (g *Grid) GetPos() raylib.Vector3 {
	return raylib.NewVector3(0, 0, 0)
}

func (c *Grid) GetPitch() float32 {
	return 0
}

func (c *Grid) SetPitch(p float32) {}

func (c *Grid) GetYaw() float32 {
	return 0
}

func (c *Grid) SetYaw(y float32) {}

func (c *Grid) GetRoll() float32 {
	return 0
}

func (c *Grid) SetRoll(r float32) {}

func (g *Grid) GetVertices() []raylib.Vector3 {
	verts := []raylib.Vector3{}
	return verts
}

func (g *Grid) GetUVs() []raylib.Vector2 {
	uvs := []raylib.Vector2{}
	return uvs
}

func (g *Grid) SetUVs(uvs []raylib.Vector2) {
}

func (c *Grid) GetMaterials() *raylib.Material {
	return &raylib.Material{}
}

func (c *Grid) SetTexture(tex raylib.Texture2D) {
}

func (c *Grid) GetTexture() raylib.Texture2D {
	return raylib.Texture2D{}
}

func (g *Grid) Prerender(cam *camera.Cam) []func() {
	return []func(){}
}

func (g *Grid) Render(cam *camera.Cam) []func() {
	sh := app.CurApp.GetShader()
	sh.Begin()
	raylib.DrawGrid(g.size, g.spacing)
	sh.End()
	return []func(){}
}

func (g *Grid) Postrender(cam *camera.Cam) []func() {
	return []func(){}
}

func (g *Grid) Update(dt float32) {
}

func (g *Grid) Collide(data pub_object.CollisionData) {
}

func (g *Grid) RegCollideHandler(handler func(pub_object.CollisionData) bool) {
}

func (g *Grid) GetCollidable() []pub_object.Object {
	return []pub_object.Object{}
}

func (g *Grid) GetCollider() pub_object.Collider {
	col := pub_object.Collider{}
	return col
}

func (g *Grid) OnAdd() {
}

func (g *Grid) OnRemove() {
}

func (g *Grid) AddChild(obj pub_object.Object) {
}

func (g *Grid) RemChild(obj pub_object.Object) {
}

func (g *Grid) GetChilds() []pub_object.Object {
	return []pub_object.Object{}
}
