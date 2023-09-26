package cell

import (
	"godev/internal/camera"
	"godev/pkg/object"
)

var ()

type Cell struct {
	childs []object.Object
	cams   []*camera.Cam
}

func (c *Cell) Init() {
	c.childs = []object.Object{}
	c.cams = []*camera.Cam{}
}

func (c *Cell) AddCam(obj *camera.Cam) {
	c.cams = append(c.cams, obj)
	obj.OnAdd()
}

func (c *Cell) RemCam(obj *camera.Cam) {
	index := -1
	for i, cam := range c.cams {
		if obj == cam {
			index = i
			break
		}
	}

	if index >= 0 {
		c.cams[index] = c.cams[len(c.cams)-1]
		c.cams = c.cams[:len(c.cams)-1]
		obj.OnRemove()
	}

}

func (c *Cell) GetCams() []*camera.Cam {
	return c.cams
}

func (c *Cell) Prerender() []func() {
	for _, cam := range c.cams {
		cmds := cam.Prerender()
		for _, child := range c.childs {
			cmds = append(child.Prerender(), cmds...)
		}
		for _, cmd := range cmds {
			cmd()
		}
	}
	return []func(){}
}

func (c *Cell) Render() []func() {
	for _, cam := range c.cams {
		cmds := cam.Render()
		for _, child := range c.childs {
			cmds = append(child.Render(), cmds...)
		}
		for _, cmd := range cmds {
			cmd()
		}
	}
	return []func(){}
}

func (c *Cell) Postrender() []func() {
	for _, cam := range c.cams {
		cmds := cam.Postrender()
		for _, child := range c.childs {
			cmds = append(child.Postrender(), cmds...)
		}
		for _, cmd := range cmds {
			cmd()
		}
	}
	return []func(){}
}

func (c *Cell) Update(dt float32) {
	for _, child := range c.childs {
		child.Update(dt)
	}
}

func (c *Cell) OnResize(w int32, h int32) {
	for _, cam := range c.cams {
		cam.OnResize(w, h)
	}
}

func (c *Cell) OnAdd() {
	for _, child := range c.childs {
		child.OnAdd()
	}
}

func (c *Cell) OnRemove() {
	for _, child := range c.childs {
		child.OnRemove()
	}
}

func (c *Cell) AddChild(obj object.Object) {
	c.childs = append(c.childs, obj)
	obj.OnAdd()
}

func (c *Cell) RemChild(obj object.Object) {
	index := -1
	for i, child := range c.childs {
		if obj == child {
			index = i
			break
		}
	}

	if index >= 0 {
		c.childs[index] = c.childs[len(c.childs)-1]
		c.childs = c.childs[:len(c.childs)-1]
		obj.OnRemove()
	}
}
