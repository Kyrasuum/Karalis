package collider

import (
	"fmt"
	"math"
	"slices"

	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Collider struct {
	obj pub_object.Object

	box    raylib.BoundingBox
	sphere pub_object.BoundingSphere

	touching      []pub_object.Object
	last_touching []pub_object.Object

	collidable []pub_object.Object
	childs     []pub_object.Object

	event_handlers map[string][]interface{}
}

func NewCollider(obj pub_object.Object) (*Collider, error) {
	mdl := obj.GetModel()
	if mdl == nil {
		return nil, fmt.Errorf("Invalid object for collisions")
	}

	col := Collider{
		obj:    obj,
		box:    raylib.GetModelBoundingBox(*mdl),
		sphere: pub_object.BoundingSphere{},
	}

	col.sphere.Center = raylib.NewVector3((col.box.Min.X+col.box.Max.X)/2, (col.box.Min.Y+col.box.Max.Y)/2, (col.box.Min.Z+col.box.Max.Z)/2)
	col.sphere.Radius = raylib.Vector3Distance(col.sphere.Center, col.box.Min)

	col.touching = []pub_object.Object{}
	col.last_touching = []pub_object.Object{}
	col.collidable = nil

	col.event_handlers = map[string][]interface{}{
		"collision":       []interface{}{},
		"start_collision": []interface{}{},
		"end_collision":   []interface{}{},
	}

	return &col, nil
}

func (c *Collider) Collide(data pub_object.CollisionData) {
	if data.Obj1 != c.obj {
		c.last_touching = append(c.last_touching, data.Obj1)
	}
	if data.Obj2 != c.obj {
		c.last_touching = append(c.last_touching, data.Obj2)
	}

	for _, handler := range c.event_handlers["collision"] {
		switch thandler := handler.(type) {
		case func(pub_object.CollisionData) bool:
			if !thandler(data) {
				break
			}
		}
	}
}
func (c *Collider) Update(dt float32) {
	for _, obj := range c.last_touching {
		if !slices.Contains(c.touching, obj) {
			data := pub_object.CollisionData{
				Obj1: c.obj,
				Obj2: obj,
			}
			for _, handler := range c.event_handlers["start_collision"] {
				switch thandler := handler.(type) {
				case func(pub_object.CollisionData) bool:
					if !thandler(data) {
						break
					}
				}
			}
		}
	}
	for _, obj := range c.touching {
		if !slices.Contains(c.last_touching, obj) {
			data := pub_object.CollisionData{
				Obj1: c.obj,
				Obj2: obj,
			}
			for _, handler := range c.event_handlers["end_collision"] {
				switch thandler := handler.(type) {
				case func(pub_object.CollisionData) bool:
					if !thandler(data) {
						break
					}
				}
			}
		}
	}

	c.touching = c.last_touching
	c.last_touching = []pub_object.Object{}
}

func (c *Collider) RegHandler(event string, handler interface{}) {
	if _, ok := c.event_handlers[event]; ok {
		c.event_handlers[event] = append(c.event_handlers[event], handler)
	}
}

func (c *Collider) GetCollidable() []pub_object.Object {
	return c.collidable
}

func (c *Collider) GetTouching() []pub_object.Object {
	return c.touching
}

func (c *Collider) GetBoundingSphere() pub_object.BoundingSphere {
	sp := c.sphere
	box := c.box
	mat := c.obj.GetModelMatrix()
	box.Max = raylib.Vector3Transform(box.Max, mat)
	box.Min = raylib.Vector3Transform(box.Min, mat)

	min := raylib.NewVector3(
		float32(math.Min(float64(box.Min.X), float64(box.Max.X))),
		float32(math.Min(float64(box.Min.Y), float64(box.Max.Y))),
		float32(math.Min(float64(box.Min.Z), float64(box.Max.Z))),
	)
	max := raylib.NewVector3(
		float32(math.Max(float64(box.Min.X), float64(box.Max.X))),
		float32(math.Max(float64(box.Min.Y), float64(box.Max.Y))),
		float32(math.Max(float64(box.Min.Z), float64(box.Max.Z))),
	)
	box.Min = min
	box.Max = max
	sp.Center = raylib.NewVector3((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2, (box.Min.Z+box.Max.Z)/2)
	sp.Radius = raylib.Vector3Distance(box.Max, box.Min) / 2
	return sp
}

func (c *Collider) GetBoundingBox() raylib.BoundingBox {
	box := c.box
	mat := c.obj.GetModelMatrix()
	box.Max = raylib.Vector3Transform(box.Max, mat)
	box.Min = raylib.Vector3Transform(box.Min, mat)

	min := raylib.NewVector3(
		float32(math.Min(float64(box.Min.X), float64(box.Max.X))),
		float32(math.Min(float64(box.Min.Y), float64(box.Max.Y))),
		float32(math.Min(float64(box.Min.Z), float64(box.Max.Z))),
	)
	max := raylib.NewVector3(
		float32(math.Max(float64(box.Min.X), float64(box.Max.X))),
		float32(math.Max(float64(box.Min.Y), float64(box.Max.Y))),
		float32(math.Max(float64(box.Min.Z), float64(box.Max.Z))),
	)
	box.Min = min
	box.Max = max
	return box
}
