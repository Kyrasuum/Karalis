package character

import (
	"karalis/pkg/object"

	pub_object "karalis/pkg/object"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Character struct {
	parent pub_object.Object
}

func NewCharacter() (c *Character, err error) {
	c = &Character{}
	err = c.Init()

	return c, err
}

func (c *Character) Init() error {
	if c == nil {
		return nil
	}
	c.parent = nil
	return nil
}

func (c *Character) Prerender(cam pub_object.Camera) []func() {
	if c == nil {
		return []func(){}
	}

	return []func(){}
}

func (c *Character) Render(cam pub_object.Camera) []func() {
	if c == nil {
		return []func(){}
	}

	return []func(){}
}

func (c *Character) Postrender(cam pub_object.Camera) []func() {
	if c == nil {
		return []func(){}
	}

	return []func(){}
}

func (c *Character) OnResize(w int32, h int32) {
	if c == nil {
		return
	}
}

func (c *Character) Update(dt float32) {
	if c == nil {
		return
	}
}

func (c *Character) GetCollider() object.Collider {
	if c == nil {
		return nil
	}

	return nil
}

func (c *Character) OnAdd(obj pub_object.Object) {
	if c == nil {
		return
	}
	c.parent = obj
}

func (c *Character) OnRemove() {
	if c == nil {
		return
	}
	c.parent = nil
}

func (c *Character) AddChild(obj object.Object) {
	if c == nil {
		return
	}
}

func (c *Character) RemChild(obj object.Object) {
	if c == nil {
		return
	}
}

func (c *Character) GetChilds() []object.Object {
	if c == nil {
		return []object.Object{}
	}

	return []object.Object{}
}

func (c *Character) GetVertices() []rl.Vector3 {
	if c == nil {
		return []rl.Vector3{}
	}

	verts := []rl.Vector3{}
	return verts
}

func (c *Character) GetUVs() []rl.Vector2 {
	if c == nil {
		return []rl.Vector2{}
	}

	uvs := []rl.Vector2{}
	return uvs
}

func (c *Character) SetUVs(uvs []rl.Vector2) {
	if c == nil {
		return
	}
}

func (c *Character) GetMaterials() *rl.Material {
	if c == nil {
		return &rl.Material{}
	}

	return &rl.Material{}
}

func (c *Character) SetTexture(tex rl.Texture2D) {
	if c == nil {
		return
	}
}

func (c *Character) GetTexture() *rl.Texture2D {
	if c == nil {
		return nil
	}

	return &rl.Texture2D{}
}

func (c *Character) GetModel() *rl.Model {
	if c == nil {
		return nil
	}

	return nil
}

func (c *Character) GetParent() pub_object.Object {
	if c == nil {
		return nil
	}
	return c.parent
}
