package physics

import (
	"math"

	"karalis/pkg/lmath"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

//
// ========================================
// RAYCAST
// ========================================
//

// Möller-Trumbore intersection
func RaycastTriangle(outT *float32, outEdge1, outEdge2 *raylib.Vector3, localOrigin, localDirection, v0, v1, v2 raylib.Vector3) bool {
	edge1 := raylib.Vector3Subtract(v1, v0)
	edge2 := raylib.Vector3Subtract(v2, v0)

	h := raylib.Vector3CrossProduct(localDirection, edge2)
	a := raylib.Vector3DotProduct(edge1, h)
	if a < 1e-5 {
		return false
	}

	f := 1 / a
	s := raylib.Vector3Subtract(localOrigin, v0)
	u := f * raylib.Vector3DotProduct(s, h)
	if u < 0 || u > 1 {
		return false
	}

	q := raylib.Vector3CrossProduct(s, edge1)
	v := f * raylib.Vector3DotProduct(localDirection, q)
	if v < 0 || u+v > 1 {
		return false
	}

	t := f * raylib.Vector3DotProduct(edge2, q)
	if t < 1e-5 {
		return false
	}

	*outT = t
	*outEdge1 = edge1
	*outEdge2 = edge2
	return true
}

func RaycastMeshVertices(closestT *float32, closestEdge1, closestEdge2 *raylib.Vector3, vertices []raylib.Vector3, triangleCount int, localOrigin, localDirection raylib.Vector3) {
	for i := 0; i < triangleCount; i++ {
		base := i * 3
		v0 := vertices[base]
		v1 := vertices[base+1]
		v2 := vertices[base+2]

		var t float32
		var e1, e2 raylib.Vector3
		if RaycastTriangle(&t, &e1, &e2, localOrigin, localDirection, v0, v1, v2) {
			if t < *closestT {
				*closestT = t
				*closestEdge1 = e1
				*closestEdge2 = e2
			}
		}
	}
}

func RaycastMeshIndexed(closestT *float32, closestEdge1, closestEdge2 *raylib.Vector3, vertices []raylib.Vector3, indices []uint32, triangleCount int, localOrigin, localDirection raylib.Vector3) {
	for i := 0; i < triangleCount; i++ {
		base := i * 3
		v0 := vertices[indices[base]]
		v1 := vertices[indices[base+1]]
		v2 := vertices[indices[base+2]]

		var t float32
		var e1, e2 raylib.Vector3
		if RaycastTriangle(&t, &e1, &e2, localOrigin, localDirection, v0, v1, v2) {
			if t < *closestT {
				*closestT = t
				*closestEdge1 = e1
				*closestEdge2 = e2
			}
		}
	}
}

func RaycastMesh(ray raylib.Ray, mesh *raylib.Mesh, transform raylib.Matrix) raylib.RayCollision {
	collision := raylib.RayCollision{Distance: float32(math.Inf(1))}

	if mesh == nil || mesh.Vertices == nil || mesh.VertexCount <= 0 || mesh.TriangleCount <= 0 {
		return collision
	}

	invTransform := raylib.MatrixInvert(transform)
	localOrigin := ray.Position.Transform(invTransform)
	localDirection := raylib.Vector3Normalize(ray.Direction.Transform(invTransform))

	closestT := float32(math.Inf(1))
	var closestEdge1, closestEdge2 raylib.Vector3

	triCount := meshTriangleCount(mesh)
	for i := 0; i < triCount; i++ {
		v0, v1, v2 := meshTrianglePositions(mesh, i)

		var t float32
		var e1, e2 raylib.Vector3
		if RaycastTriangle(&t, &e1, &e2, localOrigin, localDirection, v0, v1, v2) {
			if t < closestT {
				closestT = t
				closestEdge1 = e1
				closestEdge2 = e2
			}
		}
	}

	if closestT < float32(math.Inf(1)) {
		closestHitLocal := raylib.Vector3Add(localOrigin, raylib.Vector3Scale(localDirection, closestT))
		normalLocal := raylib.Vector3Normalize(raylib.Vector3CrossProduct(closestEdge1, closestEdge2))
		normalMatrix := raylib.MatrixTranspose(invTransform)

		collision.Hit = true
		collision.Point = closestHitLocal.Transform(transform)
		collision.Distance = raylib.Vector3Distance(ray.Position, collision.Point)
		collision.Normal = raylib.Vector3Normalize(normalLocal.Transform(normalMatrix))
	}

	return collision
}

func RaycastModel(ray raylib.Ray, model raylib.Model, transform raylib.Matrix) raylib.RayCollision {
	collision := raylib.RayCollision{Distance: float32(math.Inf(1))}

	if model.MeshCount <= 0 || model.Meshes == nil {
		return collision
	}

	for meshIdx := int32(0); meshIdx < model.MeshCount; meshIdx++ {
		mesh := &model.GetMeshes()[meshIdx]
		test := RaycastMesh(ray, mesh, transform)
		if test.Distance < collision.Distance {
			collision = test
		}
	}

	return collision
}

func RaycastObject(ray raylib.Ray, obj pub_object.Object) raylib.RayCollision {
	model := obj.GetModel()
	transform := obj.GetModelMatrix()

	aabb := obj.GetCollider().GetAABB()
	meshBoxCol := raylib.GetRayCollisionBox(ray, aabb)
	if !meshBoxCol.Hit {
		return raylib.RayCollision{Distance: float32(math.Inf(1))}
	}

	return RaycastModel(ray, *model, transform)
}

func RaycastCell(ray raylib.Ray, cell []pub_object.Object, depth int) raylib.RayCollision {
	collision := raylib.RayCollision{Distance: float32(math.Inf(1))}
	var hit interface{}
	hit = nil

	for _, obj := range cell {
		test := RaycastObject(ray, obj)
		if test.Distance < collision.Distance {
			collision = test
			hit = &obj
		}
	}
	switch hit.(type) {
	case *pub_object.Portal:
		collision = RaycastPortal(ray, hit.(pub_object.Portal), collision, depth-1)
	}

	return collision
}

func RaycastPortal(inRay raylib.Ray, p pub_object.Portal, entryHit raylib.RayCollision, depth int) raylib.RayCollision {
	collision := raylib.RayCollision{Distance: float32(math.Inf(1))}
	if !entryHit.Hit || entryHit.Distance == collision.Distance {
		return collision
	}

	// Prevent infinite recursion if portals look into each other, etc.
	if depth <= 0 {
		return collision
	}

	exit := p.GetPair()
	if exit == nil || exit.GetScene() == nil {
		return collision
	}

	// Entry/exit portal transforms.
	entryM := p.GetModelMatrix()
	exitM := exit.GetModelMatrix()

	entryInv := raylib.MatrixInvert(entryM)

	// 1) Take the *world-space* hit point and incoming direction.
	hitWorld := entryHit.Point
	dirWorld := raylib.Vector3Normalize(inRay.Direction)

	// 2) Convert hit point and direction into entry portal local space.
	hitLocal := hitWorld.Transform(entryInv)
	dirLocal := dirWorld.Transform(entryInv)
	dirLocal = raylib.Vector3Normalize(dirLocal)

	// 3) Mirror across the portal plane (assumes plane is local Z=0, normal +Z).
	// This is the "through the portal" mapping.
	hitLocal.Z = -hitLocal.Z
	dirLocal.Z = -dirLocal.Z

	// 4) Convert mapped point/direction into exit portal world space.
	exitHitWorld := hitLocal.Transform(exitM)
	exitDirWorld := dirLocal.Transform(exitM)
	exitDirWorld = raylib.Vector3Normalize(exitDirWorld)

	subRay := raylib.Ray{
		Position:  exitHitWorld,
		Direction: exitDirWorld,
	}

	// 6) Raycast into the exit portal scene.
	childs := lmath.Filter(exit.GetScene().GetChilds(), func(i int, v pub_object.Object) bool { return v != exit })
	col := RaycastCell(subRay, childs, depth)
	col.Distance += entryHit.Distance

	return col
}

//
// ========================================
// GROUNDED TESTS
// ========================================
//

func IsSphereGroundedBox(center raylib.Vector3, radius, checkDistance float32, ground raylib.BoundingBox, outGround *raylib.RayCollision) bool {
	ray := raylib.Ray{
		Position:  center,
		Direction: raylib.Vector3{X: 0, Y: -1, Z: 0},
	}
	collision := raylib.GetRayCollisionBox(ray, ground)
	grounded := collision.Hit && collision.Distance <= (radius+checkDistance)

	if outGround != nil {
		*outGround = collision
	}
	return grounded
}

func IsSphereGroundedMesh(center raylib.Vector3, radius, checkDistance float32, mesh *raylib.Mesh, transform raylib.Matrix, outGround *raylib.RayCollision) bool {
	ray := raylib.Ray{
		Position:  center,
		Direction: raylib.Vector3{X: 0, Y: -1, Z: 0},
	}
	collision := RaycastMesh(ray, mesh, transform)
	grounded := collision.Hit && collision.Distance <= (radius+checkDistance)

	if outGround != nil {
		*outGround = collision
	}
	return grounded
}

func IsCapsuleGroundedBox(capsule pub_object.Capsule, checkDistance float32, ground raylib.BoundingBox, outGround *raylib.RayCollision) bool {
	ray := raylib.Ray{
		Position:  capsule.Start,
		Direction: raylib.Vector3{X: 0, Y: -1, Z: 0},
	}
	collision := raylib.GetRayCollisionBox(ray, ground)
	grounded := collision.Hit && collision.Distance <= (capsule.Radius+checkDistance)

	if outGround != nil {
		*outGround = collision
	}
	return grounded
}

func IsCapsuleGroundedMesh(capsule pub_object.Capsule, checkDistance float32, mesh *raylib.Mesh, transform raylib.Matrix, outGround *raylib.RayCollision) bool {
	ray := raylib.Ray{
		Position:  capsule.Start,
		Direction: raylib.Vector3{X: 0, Y: -1, Z: 0},
	}
	collision := RaycastMesh(ray, mesh, transform)
	grounded := collision.Hit && collision.Distance <= (capsule.Radius+checkDistance)

	if outGround != nil {
		*outGround = collision
	}
	return grounded
}

//
// ========================================
// CLOSEST POINTS
// ========================================
//

func ClosestPointOnOBB(p raylib.Vector3, obb pub_object.OrientedBox) raylib.Vector3 {
	// Assumes obb.AxisX/Y/Z are unit length and HalfExtents are in those units.
	d := raylib.Vector3Subtract(p, obb.Center)

	// Project d onto each axis to get local coordinates.
	x := raylib.Vector3DotProduct(d, obb.AxisX)
	y := raylib.Vector3DotProduct(d, obb.AxisY)
	z := raylib.Vector3DotProduct(d, obb.AxisZ)

	// Clamp to box extents.
	x = float32(lmath.Clamp(float64(x), float64(-obb.HalfExtents.X), float64(obb.HalfExtents.X)))
	y = float32(lmath.Clamp(float64(y), float64(-obb.HalfExtents.Y), float64(obb.HalfExtents.Y)))
	z = float32(lmath.Clamp(float64(z), float64(-obb.HalfExtents.Z), float64(obb.HalfExtents.Z)))

	// Convert back to world.
	q := obb.Center
	q = raylib.Vector3Add(q, raylib.Vector3Scale(obb.AxisX, x))
	q = raylib.Vector3Add(q, raylib.Vector3Scale(obb.AxisY, y))
	q = raylib.Vector3Add(q, raylib.Vector3Scale(obb.AxisZ, z))
	return q
}

func ClosestPointOnSegment(point, start, end raylib.Vector3) raylib.Vector3 {
	dir := raylib.Vector3Subtract(end, start)
	lenSq := raylib.Vector3LengthSqr(dir)

	if lenSq < 1e-10 {
		return start
	}

	t := raylib.Vector3DotProduct(raylib.Vector3Subtract(point, start), dir) / lenSq
	t = fmax(0, fmin(1, t))

	return raylib.Vector3Add(start, raylib.Vector3Scale(dir, t))
}

func ClosestPointOnTriangle(p, a, b, c raylib.Vector3) raylib.Vector3 {
	ab := raylib.Vector3Subtract(b, a)
	ac := raylib.Vector3Subtract(c, a)
	ap := raylib.Vector3Subtract(p, a)

	d1 := raylib.Vector3DotProduct(ab, ap)
	d2 := raylib.Vector3DotProduct(ac, ap)
	if d1 <= 0 && d2 <= 0 {
		return a
	}

	bp := raylib.Vector3Subtract(p, b)
	d3 := raylib.Vector3DotProduct(ab, bp)
	d4 := raylib.Vector3DotProduct(ac, bp)
	if d3 >= 0 && d4 <= d3 {
		return b
	}

	cp := raylib.Vector3Subtract(p, c)
	d5 := raylib.Vector3DotProduct(ab, cp)
	d6 := raylib.Vector3DotProduct(ac, cp)
	if d6 >= 0 && d5 <= d6 {
		return c
	}

	vc := d1*d4 - d3*d2
	if vc <= 0 && d1 >= 0 && d3 <= 0 {
		v := d1 / (d1 - d3)
		return raylib.Vector3Add(a, raylib.Vector3Scale(ab, v))
	}

	vb := d5*d2 - d1*d6
	if vb <= 0 && d2 >= 0 && d6 <= 0 {
		v := d2 / (d2 - d6)
		return raylib.Vector3Add(a, raylib.Vector3Scale(ac, v))
	}

	va := d3*d6 - d5*d4
	if va <= 0 && (d4-d3) >= 0 && (d5-d6) >= 0 {
		v := (d4 - d3) / ((d4 - d3) + (d5 - d6))
		return raylib.Vector3Add(b, raylib.Vector3Scale(raylib.Vector3Subtract(c, b), v))
	}

	denom := 1 / (va + vb + vc)
	v := vb * denom
	w := vc * denom

	return raylib.Vector3Add(a, raylib.Vector3Add(raylib.Vector3Scale(ab, v), raylib.Vector3Scale(ac, w)))
}

func ClosestPointOnBox(point raylib.Vector3, box raylib.BoundingBox) raylib.Vector3 {
	return raylib.Vector3{
		X: fmax(box.Min.X, fmin(point.X, box.Max.X)),
		Y: fmax(box.Min.Y, fmin(point.Y, box.Max.Y)),
		Z: fmax(box.Min.Z, fmin(point.Z, box.Max.Z)),
	}
}
