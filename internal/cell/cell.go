package cell

import (
	"slices"

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

		switch child.(type) {
		case *object.Skybox:
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
	cmds := []func(){}
	for _, child := range c.childs {
		switch child.(type) {
		case *object.Skybox:
		default:
			cmds = append(cmds, child.Render(cam)...)
		}
	}
	return cmds
}

func (c *Cell) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	for _, child := range c.childs {
		cmds = append(cmds, child.Postrender(cam)...)
	}
	return cmds
}

func (c *Cell) Update(dt float32) {
	//check for collisions
	collision_check := func(child pub_object.Object, pairs []pub_object.Object) {
		for _, pair := range pairs {
			//check if pair can collide with child
			collidable := pair.GetCollidable()
			if collidable == nil || slices.Contains(collidable, child) {
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
	}

	//collisions check loop
	childs := c.GetChilds()
	for i, child := range childs {
		collidable := child.GetCollidable()
		if collidable == nil {
			//undefined collidable means all objects collidable
			collision_check(child, childs[i+1:])
		} else {
			//defined collidable means only collidable can collide with child
			collision_check(child, collidable)
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

	return slices.Concat(grandchilds, childs)
}
