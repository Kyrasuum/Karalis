package world

import (
	"errors"
	"image/color"
	"slices"

	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Cell struct {
	parent  pub_object.Object
	terrain pub_object.Object
	childs  []pub_object.Object
}

func NewTerrainCell(pos, sc raylib.Vector3) (*Cell, error) {
	c := &Cell{}
	err := c.Init()
	if err != nil {
		return nil, err
	}

	ter, err := RandTerrain(float64(pos.X), float64(pos.Z), 256, 256, 1234567)
	if err != nil {
		return nil, err
	}
	ter.SetScale(sc)
	ter.SetPos(pos)
	c.terrain = ter

	return c, nil
}

func NewCityCell(pos, sc raylib.Vector3) (*Cell, error) {
	c := &Cell{}
	err := c.Init()
	if err != nil {
		return nil, err
	}

	ter, err := RandDungeon(int(pos.X), int(pos.Z), 256, 256, 1234567)
	if err != nil {
		return nil, err
	}
	ter.SetScale(sc)
	ter.SetPos(pos)
	c.terrain = ter

	return c, nil
}

func NewDungeonCell(pos, sc raylib.Vector3) (*Cell, error) {
	c := &Cell{}
	err := c.Init()
	if err != nil {
		return nil, err
	}

	ter, err := RandCity(int(pos.X), int(pos.Z), 256, 256, 1234567)
	if err != nil {
		return nil, err
	}
	ter.SetScale(pos)
	ter.SetPos(sc)
	c.terrain = ter

	return c, nil
}

func (c *Cell) Init() error {
	if c == nil {
		return errors.New("Invalid cell")
	}
	c.parent = nil
	c.terrain = nil
	c.childs = []pub_object.Object{}

	return nil
}

func (c *Cell) Prerender(cam pub_object.Camera) []func() {
	if c == nil {
		return []func(){}
	}

	cmds := []func(){}
	if c.terrain != nil {
		cmds = c.terrain.Prerender(cam)
	}
	for _, child := range c.childs {
		switch child.(type) {
		default:
			cmds = append(cmds, child.Prerender(cam)...)
		}
	}
	return cmds
}

func (c *Cell) Render(cam pub_object.Camera) []func() {
	if c == nil {
		return []func(){}
	}

	cmds := []func(){}
	if c.terrain != nil {
		cmds = c.terrain.Render(cam)
	}
	for _, child := range c.childs {
		switch child.(type) {
		default:
			cmds = append(cmds, child.Render(cam)...)
		}
	}
	return cmds
}

func (c *Cell) Postrender(cam pub_object.Camera) []func() {
	if c == nil {
		return []func(){}
	}

	cmds := []func(){}
	if c.terrain != nil {
		cmds = c.terrain.Postrender(cam)
	}
	for _, child := range c.childs {
		cmds = append(cmds, child.Postrender(cam)...)
	}
	return cmds
}

func (c *Cell) OnResize(w int32, h int32) {
	if c == nil {
		return
	}

	//perform update on objects
	for _, child := range c.childs {
		child.OnResize(w, h)
	}
	if c.terrain != nil {
		c.terrain.OnResize(w, h)
	}
}

func (c *Cell) Update(dt float32) {
	if c == nil {
		return
	}

	//perform update on objects
	for _, child := range c.childs {
		child.Update(dt)
	}
	if c.terrain != nil {
		c.terrain.Update(dt)
	}
}

func (c *Cell) OnAdd(obj pub_object.Object) {
	if c == nil {
		return
	}
	c.parent = obj
}

func (c *Cell) OnRemove() {
	if c == nil {
		return
	}
	c.parent = nil
}

func (c *Cell) AddChild(obj pub_object.Object) {
	if c == nil {
		return
	}

	c.childs = append(c.childs, obj)
	obj.OnAdd(c)
}

func (c *Cell) RemChild(obj pub_object.Object) {
	if c == nil {
		return
	}
	//find index of child
	index := -1
	for i, child := range c.childs {
		if obj == child {
			index = i
			break
		}
	}

	//remove child
	if index >= 0 {
		c.childs[index] = c.childs[len(c.childs)-1]
		c.childs = c.childs[:len(c.childs)-1]
		obj.OnRemove()
	}
}

func (c *Cell) GetChilds() []pub_object.Object {
	if c == nil {
		return []pub_object.Object{}
	}

	childs := c.childs
	if c.terrain != nil {
		childs = append(c.childs, c.terrain)
	}
	grandchilds := []pub_object.Object{}
	for _, child := range childs {
		grandchilds = append(grandchilds, child.GetChilds()...)
	}

	return slices.Concat(grandchilds, childs)
}

func (c *Cell) GetCollider() pub_object.Collider {
	if c == nil {
		return nil
	}

	return nil
}

func (c *Cell) GetModelMatrix() raylib.Matrix {
	if c == nil || c.terrain == nil {
		return raylib.Matrix{}
	}
	return c.terrain.GetModelMatrix()
}

func (c *Cell) GetModel() *raylib.Model {
	if c == nil || c.terrain == nil {
		return nil
	}

	return c.terrain.GetModel()
}

func (c *Cell) SetColor(col color.Color) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetColor(col)
}

func (c *Cell) GetColor() color.Color {
	if c == nil || c.terrain == nil {
		return nil
	}

	return c.terrain.GetColor()
}

func (c *Cell) GetScale() raylib.Vector3 {
	if c == nil || c.terrain == nil {
		return raylib.Vector3{}
	}

	return c.terrain.GetScale()
}

func (c *Cell) SetScale(sc raylib.Vector3) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetScale(sc)
}

func (c *Cell) SetPos(pos raylib.Vector3) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetPos(pos)
}

func (c *Cell) GetPos() raylib.Vector3 {
	if c == nil || c.terrain == nil {
		return raylib.Vector3{}
	}

	return c.terrain.GetPos()
}

func (c *Cell) GetPitch() float32 {
	if c == nil || c.terrain == nil {
		return 0
	}

	return c.terrain.GetPitch()
}

func (c *Cell) SetPitch(pitch float32) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetPitch(pitch)
}

func (c *Cell) GetYaw() float32 {
	if c == nil || c.terrain == nil {
		return 0
	}

	return c.terrain.GetYaw()
}

func (c *Cell) SetYaw(yaw float32) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetYaw(yaw)
}

func (c *Cell) GetRoll() float32 {
	if c == nil || c.terrain == nil {
		return 0
	}

	return c.terrain.GetRoll()
}

func (c *Cell) SetRoll(roll float32) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetRoll(roll)
}

func (c *Cell) GetVertices() []raylib.Vector3 {
	if c == nil || c.terrain == nil {
		return []raylib.Vector3{}
	}
	return c.terrain.GetVertices()
}

func (c *Cell) GetUVs() []raylib.Vector2 {
	if c == nil || c.terrain == nil {
		return []raylib.Vector2{}
	}
	return c.terrain.GetUVs()
}

func (c *Cell) SetUVs(uvs []raylib.Vector2) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetUVs(uvs)
}

func (c *Cell) GetMaterials() *raylib.Material {
	if c == nil || c.terrain == nil {
		return nil
	}

	return c.terrain.GetMaterials()
}

func (c *Cell) SetTexture(tex raylib.Texture2D) {
	if c == nil || c.terrain == nil {
		return
	}
	c.terrain.SetTexture(tex)
}

func (c *Cell) GetTexture() *raylib.Texture2D {
	if c == nil || c.terrain == nil {
		return nil
	}

	return c.terrain.GetTexture()
}

func (c *Cell) GetParent() pub_object.Object {
	if c == nil {
		return nil
	}
	return c.parent
}
