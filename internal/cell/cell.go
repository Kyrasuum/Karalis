package cell

import (
	"karalis/internal/camera"
	"karalis/internal/object"
	pub_object "karalis/pkg/object"
)

var ()

type Cell struct {
	childs []pub_object.Object
}

func (c *Cell) Init() error {
	c.childs = []pub_object.Object{}

	sky, err := object.NewSkybox(nil)
	if err != nil {
		return err
	}
	c.AddChild(sky)

	return nil
}

func (c *Cell) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}
	for _, child := range c.childs {
		cmds = append(child.Prerender(cam), cmds...)
	}
	return cmds
}

func (c *Cell) Render(cam *camera.Cam) []func() {
	cmds := []func(){}
	for _, child := range c.childs {
		cmds = append(child.Render(cam), cmds...)
	}
	return cmds
}

func (c *Cell) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	for _, child := range c.childs {
		cmds = append(child.Postrender(cam), cmds...)
	}
	return cmds
}

func (c *Cell) Update(dt float32) {
	//check for collisions
	childs := c.GetChilds()
	for i, child := range childs {
		if !child.CanCollide() {
			continue
		}
		for _, pair := range childs[i+1:] {
			if !pair.CanCollide() {
				continue
			}
			//broad phase collision detection
			if pub_object.CheckCollisionSpheres(child.GetCollider().Sphere, pair.GetCollider().Sphere) {
				//handle collision
				coldata := pub_object.CollisionData{
					Obj1: child,
					Obj2: pair,
				}
				child.Collide(coldata)
				pair.Collide(coldata)
			}
		}
	}

	//perform update on objects
	for _, child := range c.childs {
		child.Update(dt)
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

func (c *Cell) AddChild(obj pub_object.Object) {
	c.childs = append(c.childs, obj)
	obj.OnAdd()
}

func (c *Cell) RemChild(obj pub_object.Object) {
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
	childs := c.childs
	grandchilds := []pub_object.Object{}
	for _, child := range childs {
		grandchilds = append(grandchilds, child.GetChilds()...)
	}

	return append(grandchilds, childs...)
}
