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

func (c *Collider) GetBoundingSphere() pub_object.BoundingSphere {
	sp := pub_object.BoundingSphere{}
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

	box := raylib.GetModelBoundingBox(*c.obj.GetModel())
	mat := c.obj.GetModelMatrix()

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

func (c *Collider) GetOOBB() pub_object.OBB {
	if c == nil {
		return pub_object.OBB{}
	}

	box := raylib.GetModelBoundingBox(*c.obj.GetModel())
	mat := c.obj.GetModelMatrix()
	// local center + half extents
	localCenter := raylib.NewVector3(
		(box.Min.X+box.Max.X)/2,
		(box.Min.Y+box.Max.Y)/2,
		(box.Min.Z+box.Max.Z)/2,
	)
	localHalf := raylib.NewVector3(
		(box.Max.X-box.Min.X)/2,
		(box.Max.Y-box.Min.Y)/2,
		(box.Max.Z-box.Min.Z)/2,
	)

	// Extract basis vectors from matrix columns (raylib matrices are column-major in concept)
	// In raylib.Matrix: M0 M4 M8  M12
	//                  M1 M5 M9  M13
	//                  M2 M6 M10 M14
	//                  M3 M7 M11 M15
	xCol := raylib.NewVector3(mat.M0, mat.M1, mat.M2)
	yCol := raylib.NewVector3(mat.M4, mat.M5, mat.M6)
	zCol := raylib.NewVector3(mat.M8, mat.M9, mat.M10)

	// Axis directions (normalized)
	ax := raylib.Vector3Normalize(xCol)
	ay := raylib.Vector3Normalize(yCol)
	az := raylib.Vector3Normalize(zCol)

	// Scale magnitudes baked into columns
	sx := raylib.Vector3Length(xCol)
	sy := raylib.Vector3Length(yCol)
	sz := raylib.Vector3Length(zCol)

	// World half extents = local half * scale along each axis
	worldHalf := raylib.NewVector3(localHalf.X*sx, localHalf.Y*sy, localHalf.Z*sz)

	// World center: transform local center by m
	worldCenter := raylib.Vector3Transform(localCenter, mat)

	return pub_object.OBB{
		Center: worldCenter,
		Axis:   [3]raylib.Vector3{ax, ay, az},
		Half:   worldHalf,
	}
}
