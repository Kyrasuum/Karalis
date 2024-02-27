package object

import (
	"karalis/internal/camera"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Skybox struct {
}

func (s *Skybox) Init() error {
	return nil
}

func (s *Skybox) Prerender(cam *camera.Cam) []func() {
	return []func(){}
}

func (s *Skybox) Render(cam *camera.Cam) []func() {
	return []func(){}
}

func (s *Skybox) Postrender(cam *camera.Cam) []func() {
	return []func(){}
}

func (s *Skybox) Update(dt float32) {
}

func (s *Skybox) OnAdd() {
}

func (s *Skybox) OnRemove() {
}

func (s *Skybox) AddChild(obj pub_object.Object) {
}

func (s *Skybox) RemChild(obj pub_object.Object) {
}

func (s *Skybox) GetPos() raylib.Vector3 {
	return raylib.NewVector3(0, 0, 0)
}

func (s *Skybox) GetModelMatrix() raylib.Matrix {
	return raylib.MatrixTranslate(0, 0, 0)
}

func (s *Skybox) GetMaterials() *raylib.Material {
	return &raylib.Material{}
}

func (s *Skybox) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
}

func (s *Skybox) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return raylib.Texture2D{}
}
