package collider

import (
	"fmt"
	"math"
	"slices"

	pub_object "karalis/pkg/object"
	_ "karalis/pkg/physics"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Collider struct {
	obj pub_object.Object

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
		obj: obj,
	}

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

func (c *Collider) GetObj() pub_object.Object {
	if c == nil {
		return nil
	}

	return c.obj
}

func (c *Collider) Collide(data pub_object.CollisionData) {
	if c == nil {
		return
	}

	// add other object as touching
	var other pub_object.Object
	if data.Obj1 == c.obj {
		other = data.Obj2
	} else if data.Obj2 == c.obj {
		other = data.Obj1
	}
	c.touching = append(c.touching, other)

	// initial collision handlers
	if !slices.Contains(c.last_touching, other) {
		for _, handler := range c.event_handlers["start_collision"] {
			switch thandler := handler.(type) {
			case func(pub_object.CollisionData) bool:
				if !thandler(data) {
					break
				}
			}
		}
	}

	// ongoing collision handlers
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
	if c == nil {
		return
	}

	// look for objects that stopped colliding and call handler for
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

	c.last_touching = c.touching
	c.touching = []pub_object.Object{}
}

func (c *Collider) RegHandler(event string, handler interface{}) {
	if c == nil {
		return
	}

	if _, ok := c.event_handlers[event]; ok {
		c.event_handlers[event] = append(c.event_handlers[event], handler)
	}
}

func (c *Collider) GetCollidable() []pub_object.Object {
	if c == nil {
		return []pub_object.Object{}
	}

	return c.collidable
}

func (c *Collider) GetTouching() []pub_object.Object {
	if c == nil {
		return []pub_object.Object{}
	}

	if len(c.touching) == 0 {
		return c.last_touching
	}
	return c.touching
}

func (c *Collider) GetBoundingSphere() pub_object.Sphere {
	sp := pub_object.Sphere{}
	if c == nil {
		return sp
	}

	box := c.GetAABB()
	sp.Center = raylib.NewVector3((box.Min.X+box.Max.X)/2, (box.Min.Y+box.Max.Y)/2, (box.Min.Z+box.Max.Z)/2)
	sp.Radius = raylib.Vector3Distance(box.Max, box.Min) / 2
	return sp
}

func (c *Collider) GetAABB() raylib.BoundingBox {
	if c == nil {
		return raylib.BoundingBox{}
	}

	mat := c.obj.GetModelMatrix()
	box := raylib.GetModelBoundingBox(*c.obj.GetModel())

	corners := [8]raylib.Vector3{
		{box.Min.X, box.Min.Y, box.Min.Z},
		{box.Min.X, box.Min.Y, box.Max.Z},
		{box.Min.X, box.Max.Y, box.Min.Z},
		{box.Min.X, box.Max.Y, box.Max.Z},
		{box.Max.X, box.Min.Y, box.Min.Z},
		{box.Max.X, box.Min.Y, box.Max.Z},
		{box.Max.X, box.Max.Y, box.Min.Z},
		{box.Max.X, box.Max.Y, box.Max.Z},
	}

	min := raylib.Vector3Transform(corners[0], mat)
	max := min
	for i := 1; i < 8; i++ {
		p := raylib.Vector3Transform(corners[i], mat)
		min = raylib.NewVector3(
			float32(math.Min(float64(min.X), float64(p.X))),
			float32(math.Min(float64(min.Y), float64(p.Y))),
			float32(math.Min(float64(min.Z), float64(p.Z))),
		)
		max = raylib.NewVector3(
			float32(math.Max(float64(max.X), float64(p.X))),
			float32(math.Max(float64(max.Y), float64(p.Y))),
			float32(math.Max(float64(max.Z), float64(p.Z))),
		)
	}
	box.Min = min
	box.Max = max
	return box
}

func (c *Collider) GetOOBB() pub_object.OrientedBox {
	if c == nil {
		return pub_object.OrientedBox{}
	}

	aabb := raylib.GetModelBoundingBox(*c.obj.GetModel())
	transform := c.obj.GetModelMatrix()

	var obb pub_object.OrientedBox

	obb.HalfExtents = raylib.Vector3{
		X: (aabb.Max.X - aabb.Min.X) * 0.5,
		Y: (aabb.Max.Y - aabb.Min.Y) * 0.5,
		Z: (aabb.Max.Z - aabb.Min.Z) * 0.5,
	}

	localCenter := raylib.Vector3{
		X: (aabb.Min.X + aabb.Max.X) * 0.5,
		Y: (aabb.Min.Y + aabb.Max.Y) * 0.5,
		Z: (aabb.Min.Z + aabb.Max.Z) * 0.5,
	}

	obb.Center = raylib.Vector3{
		X: transform.M0*localCenter.X + transform.M4*localCenter.Y + transform.M8*localCenter.Z + transform.M12,
		Y: transform.M1*localCenter.X + transform.M5*localCenter.Y + transform.M9*localCenter.Z + transform.M13,
		Z: transform.M2*localCenter.X + transform.M6*localCenter.Y + transform.M10*localCenter.Z + transform.M14,
	}

	obb.AxisX = raylib.Vector3{X: transform.M0, Y: transform.M1, Z: transform.M2}
	obb.AxisY = raylib.Vector3{X: transform.M4, Y: transform.M5, Z: transform.M6}
	obb.AxisZ = raylib.Vector3{X: transform.M8, Y: transform.M9, Z: transform.M10}

	return obb
}
