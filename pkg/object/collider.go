package object

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

//
// ========================================
// HELPER TYPES
// ========================================
//

type Rect struct {
	X, Y int
	W, H int
}

type OrientedBox struct {
	Center      raylib.Vector3
	AxisX       raylib.Vector3
	AxisY       raylib.Vector3
	AxisZ       raylib.Vector3
	HalfExtents raylib.Vector3
}

type Capsule struct {
	Start  raylib.Vector3
	End    raylib.Vector3
	Radius float32
}

type Sphere struct {
	Center raylib.Vector3
	Radius float32
}

type Plane struct {
	Normal raylib.Vector3
	Offset float32
}

type CollisionData struct {
	Obj1 Object
	Obj2 Object
}

type Collider interface {
	GetObj() Object
	GetBoundingSphere() Sphere
	GetAABB() raylib.BoundingBox
	GetOOBB() OrientedBox
	GetCollidable() []Object
	Collide(CollisionData)
	Update(dt float32)
	RegHandler(string, interface{})
	GetTouching() []Object
}
