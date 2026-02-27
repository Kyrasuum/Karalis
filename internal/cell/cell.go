package cell

import (
	"errors"
	"slices"

	"karalis/internal/camera"
	"karalis/internal/object/world"
	pub_object "karalis/pkg/object"
)

var ()

type Cell struct {
	childs []pub_object.Object
}

func (c *Cell) Init() error {
	if c == nil {
		return errors.New("Invalid cell")
	}

	c.childs = []pub_object.Object{}

	sky, err := world.NewSkybox(nil)
	if err != nil {
		return err
	}
	c.AddChild(sky)

	return nil
}

func (c *Cell) Prerender(cam *camera.Cam) []func() {
	if c == nil {
		return []func(){}
	}

	cmds := []func(){}
	for _, child := range c.childs {

		switch child.(type) {
		case *world.Skybox:
			cmds = append(cmds, child.Prerender(cam)...)
			cmds = append(cmds, cam.Render()...)
			cmds = append(cmds, child.Render(cam)...)
			cmds = append(cmds, cam.Prerender()...)
		default:
			cmds = append(cmds, child.Prerender(cam)...)
		}
	}
	return cmds
}

func (c *Cell) Render(cam *camera.Cam) []func() {
	if c == nil {
		return []func(){}
	}

	cmds := []func(){}
	for _, child := range c.childs {
		switch child.(type) {
		case *world.Skybox:
		default:
			cmds = append(cmds, child.Render(cam)...)
		}
	}
	return cmds
}

func (c *Cell) Postrender(cam *camera.Cam) []func() {
	if c == nil {
		return []func(){}
	}

	cmds := []func(){}
	for _, child := range c.childs {
		cmds = append(cmds, child.Postrender(cam)...)
	}
	return cmds
}

func (c *Cell) Update(dt float32) {
	if c == nil {
		return
	}

	//perform update on objects
	for _, child := range c.childs {
		child.Update(dt)
	}
}

func (c *Cell) OnAdd() {
	if c == nil {
		return
	}

	for _, child := range c.childs {
		child.OnAdd()
	}
}

func (c *Cell) OnRemove() {
	if c == nil {
		return
	}

	for _, child := range c.childs {
		child.OnRemove()
	}
}

func (c *Cell) AddChild(obj pub_object.Object) {
	if c == nil {
		return
	}

	c.childs = append(c.childs, obj)
	obj.OnAdd()
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
	grandchilds := []pub_object.Object{}
	for _, child := range childs {
		grandchilds = append(grandchilds, child.GetChilds()...)
	}

	return slices.Concat(grandchilds, childs)
}
