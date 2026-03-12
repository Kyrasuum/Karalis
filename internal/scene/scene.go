package scene

import (
	"errors"
	"image/color"
	"slices"

	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Scene struct {
	childs []pub_object.Object
	parent pub_object.Object
}

func (s *Scene) Init() error {
	if s == nil {
		return errors.New("Invalid scene")
	}

	s.parent = nil
	s.childs = []pub_object.Object{}

	return nil
}

func (s *Scene) Prerender(cam pub_object.Camera) []func() {
	if s == nil {
		return []func(){}
	}

	cmds := []func(){}
	for _, child := range s.childs {
		switch child.(type) {
		default:
			cmds = append(cmds, child.Prerender(cam)...)
		}
	}
	return cmds
}

func (s *Scene) Render(cam pub_object.Camera) []func() {
	if s == nil {
		return []func(){}
	}

	cmds := []func(){}
	for _, child := range s.childs {
		switch child.(type) {
		default:
			cmds = append(cmds, child.Render(cam)...)
		}
	}
	return cmds
}

func (s *Scene) Postrender(cam pub_object.Camera) []func() {
	if s == nil {
		return []func(){}
	}

	cmds := []func(){}
	for _, child := range s.childs {
		cmds = append(cmds, child.Postrender(cam)...)
	}
	return cmds
}

func (s *Scene) OnResize(w int32, h int32) {
	if s == nil {
		return
	}

	//perform update on objects
	for _, child := range s.childs {
		child.OnResize(w, h)
	}
}

func (s *Scene) Update(dt float32) {
	if s == nil {
		return
	}

	//perform update on objects
	for _, child := range s.childs {
		child.Update(dt)
	}
}

func (s *Scene) OnAdd(obj pub_object.Object) {
	if s == nil {
		return
	}
	s.parent = obj
}

func (s *Scene) OnRemove() {
	if s == nil {
		return
	}
	s.parent = nil
}

func (s *Scene) AddChild(obj pub_object.Object) {
	if s == nil {
		return
	}

	s.childs = append(s.childs, obj)
	obj.OnAdd(s)
}

func (s *Scene) RemChild(obj pub_object.Object) {
	if s == nil {
		return
	}
	//find index of child
	index := -1
	for i, child := range s.childs {
		if obj == child {
			index = i
			break
		}
	}

	//remove child
	if index >= 0 {
		s.childs[index] = s.childs[len(s.childs)-1]
		s.childs = s.childs[:len(s.childs)-1]
		obj.OnRemove()
	}
}

func (s *Scene) GetChilds() []pub_object.Object {
	if s == nil {
		return []pub_object.Object{}
	}

	childs := s.childs
	grandchilds := []pub_object.Object{}
	for _, child := range childs {
		grandchilds = append(grandchilds, child.GetChilds()...)
	}

	return slices.Concat(grandchilds, childs)
}

func (s *Scene) GetCollider() pub_object.Collider {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Scene) GetModelMatrix() raylib.Matrix {
	if s == nil {
		return raylib.Matrix{}
	}
	return raylib.Matrix{}
}

func (s *Scene) GetModel() *raylib.Model {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Scene) SetColor(col color.Color) {
	if s == nil {
		return
	}
}

func (s *Scene) GetColor() color.Color {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Scene) GetScale() raylib.Vector3 {
	if s == nil {
		return raylib.Vector3{}
	}

	return raylib.Vector3{}
}

func (s *Scene) SetScale(sc raylib.Vector3) {
	if s == nil {
		return
	}
}

func (s *Scene) SetPos(pos raylib.Vector3) {
	if s == nil {
		return
	}
}

func (s *Scene) GetPos() raylib.Vector3 {
	if s == nil {
		return raylib.Vector3{}
	}

	return raylib.Vector3{}
}

func (s *Scene) GetPitch() float32 {
	if s == nil {
		return 0
	}

	return 0
}

func (s *Scene) SetPitch(pitch float32) {
	if s == nil {
		return
	}
}

func (s *Scene) GetYaw() float32 {
	if s == nil {
		return 0
	}

	return 0
}

func (s *Scene) SetYaw(yaw float32) {
	if s == nil {
		return
	}
}

func (s *Scene) GetRoll() float32 {
	if s == nil {
		return 0
	}

	return 0
}

func (s *Scene) SetRoll(roll float32) {
	if s == nil {
		return
	}
}

func (s *Scene) GetVertices() []raylib.Vector3 {
	if s == nil {
		return []raylib.Vector3{}
	}
	return []raylib.Vector3{}
}

func (s *Scene) GetUVs() []raylib.Vector2 {
	if s == nil {
		return []raylib.Vector2{}
	}
	return []raylib.Vector2{}
}

func (s *Scene) SetUVs(uvs []raylib.Vector2) {
	if s == nil {
		return
	}
}

func (s *Scene) GetMaterials() *raylib.Material {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Scene) SetTexture(tex raylib.Texture2D) {
	if s == nil {
		return
	}
}

func (s *Scene) GetTexture() *raylib.Texture2D {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Scene) GetParent() pub_object.Object {
	if s == nil {
		return nil
	}
	return s.parent
}
