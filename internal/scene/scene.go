package scene

import (
	"errors"
	"image/color"
	"slices"

	pub_object "karalis/pkg/object"

	rl "github.com/gen2brain/raylib-go/raylib"
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

func (s *Scene) GetModelMatrix() rl.Matrix {
	if s == nil {
		return rl.Matrix{}
	}
	return rl.Matrix{}
}

func (s *Scene) GetModel() *rl.Model {
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

func (s *Scene) GetScale() rl.Vector3 {
	if s == nil {
		return rl.Vector3{}
	}

	return rl.Vector3{}
}

func (s *Scene) SetScale(sc rl.Vector3) {
	if s == nil {
		return
	}
}

func (s *Scene) SetPos(pos rl.Vector3) {
	if s == nil {
		return
	}
}

func (s *Scene) GetPos() rl.Vector3 {
	if s == nil {
		return rl.Vector3{}
	}

	return rl.Vector3{}
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

func (s *Scene) GetVertices() []rl.Vector3 {
	if s == nil {
		return []rl.Vector3{}
	}
	return []rl.Vector3{}
}

func (s *Scene) GetUVs() []rl.Vector2 {
	if s == nil {
		return []rl.Vector2{}
	}
	return []rl.Vector2{}
}

func (s *Scene) SetUVs(uvs []rl.Vector2) {
	if s == nil {
		return
	}
}

func (s *Scene) GetMaterials() *rl.Material {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Scene) SetTexture(tex rl.Texture2D) {
	if s == nil {
		return
	}
}

func (s *Scene) GetTexture() *rl.Texture2D {
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
