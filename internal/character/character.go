package character

import (
	"karalis/internal/camera"
	"karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Character struct {
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
	return nil
}

func (c *Character) Prerender(cam *camera.Cam) []func() {
	if c == nil {
		return []func(){}
	}

	return []func(){}
}

func (c *Character) Render(cam *camera.Cam) []func() {
	if c == nil {
		return []func(){}
	}

	return []func(){}
}

func (c *Character) Postrender(cam *camera.Cam) []func() {
	if c == nil {
		return []func(){}
	}

	return []func(){}
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

func (c *Character) OnAdd() {
	if c == nil {
		return
	}
}

func (c *Character) OnRemove() {
	if c == nil {
		return
	}
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

func (c *Character) GetVertices() []raylib.Vector3 {
	if c == nil {
		return []raylib.Vector3{}
	}

	verts := []raylib.Vector3{}
	return verts
}

func (c *Character) GetUVs() []raylib.Vector2 {
	if c == nil {
		return []raylib.Vector2{}
	}

	uvs := []raylib.Vector2{}
	return uvs
}

func (c *Character) SetUVs(uvs []raylib.Vector2) {
	if c == nil {
		return
	}
}

func (c *Character) GetMaterials() *raylib.Material {
	if c == nil {
		return &raylib.Material{}
	}

	return &raylib.Material{}
}

func (c *Character) SetTexture(tex raylib.Texture2D) {
	if c == nil {
		return
	}
}

func (c *Character) GetTexture() *raylib.Texture2D {
	if c == nil {
		return nil
	}

	return &raylib.Texture2D{}
}

func (c *Character) GetModel() *raylib.Model {
	if c == nil {
		return nil
	}

	return nil
}
