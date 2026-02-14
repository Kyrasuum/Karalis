package object

import (
	"math"

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

type Penetration struct {
	Collides bool
	Depth    float32
	Normal   raylib.Vector3
	MTV      raylib.Vector3
}

type SweepCollision struct {
	Hit    bool
	Time   float32
	Point  raylib.Vector3
	Normal raylib.Vector3
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

func ComputeAABB(obj raylib.Model, mat raylib.Matrix) raylib.BoundingBox {
	box := raylib.GetModelBoundingBox(obj)

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

func ComputeOBB(aabb raylib.BoundingBox, transform raylib.Matrix) OrientedBox {
	var obb OrientedBox

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

// If your axes came from a transform matrix columns, they may include scale.
// This normalizes axes and pushes scale into HalfExtents so later math is correct.
func OrientedBoxNormalizeScale(obb OrientedBox) OrientedBox {
	const eps = 1e-8

	// X
	lx := raylib.Vector3Length(obb.AxisX)
	if lx > eps {
		inv := 1.0 / lx
		obb.AxisX = raylib.Vector3Scale(obb.AxisX, float32(inv))
		obb.HalfExtents.X *= lx
	}

	// Y
	ly := raylib.Vector3Length(obb.AxisY)
	if ly > eps {
		inv := 1.0 / ly
		obb.AxisY = raylib.Vector3Scale(obb.AxisY, float32(inv))
		obb.HalfExtents.Y *= ly
	}

	// Z
	lz := raylib.Vector3Length(obb.AxisZ)
	if lz > eps {
		inv := 1.0 / lz
		obb.AxisZ = raylib.Vector3Scale(obb.AxisZ, float32(inv))
		obb.HalfExtents.Z *= lz
	}

	return obb
}
