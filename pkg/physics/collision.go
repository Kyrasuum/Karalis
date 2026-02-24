package physics

import (
	"math"

	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

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

func fmin(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func fmax(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

//
// ========================================
// COLLISION
// ========================================
//

func CheckCollisionSpheres(sp1 pub_object.Sphere, sp2 pub_object.Sphere) bool {
	return raylib.Vector3Length(raylib.Vector3Subtract(sp1.Center, sp2.Center)) < sp1.Radius+sp2.Radius
}

func CheckCollisionAABB(b1 raylib.BoundingBox, b2 raylib.BoundingBox) bool {
	if b1.Max.X <= b2.Min.X || b1.Min.X >= b2.Max.X {
		return false
	}
	if b1.Max.Y <= b2.Min.Y || b1.Min.Y >= b2.Max.Y {
		return false
	}
	if b1.Max.Z <= b2.Min.Z || b1.Min.Z >= b2.Max.Z {
		return false
	}
	return true
}

// CheckCollisionOBB performs OBB-OBB overlap using the standard Separating Axis Theorem.
// Fast, robust, and the usual go-to for oriented boxes.
func CheckCollisionOBB(a, b pub_object.OrientedBox) bool {
	a = pub_object.OrientedBoxNormalizeScale(a)
	b = pub_object.OrientedBoxNormalizeScale(b)

	A := [3]raylib.Vector3{a.AxisX, a.AxisY, a.AxisZ}
	B := [3]raylib.Vector3{b.AxisX, b.AxisY, b.AxisZ}
	ae := [3]float32{a.HalfExtents.X, a.HalfExtents.Y, a.HalfExtents.Z}
	be := [3]float32{b.HalfExtents.X, b.HalfExtents.Y, b.HalfExtents.Z}

	// Rotation matrix R[i][j] = Ai · Bj
	var R [3][3]float32
	var AbsR [3][3]float32

	const eps = 1e-6 // helps when axes are nearly parallel

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			R[i][j] = raylib.Vector3DotProduct(A[i], B[j])
			AbsR[i][j] = float32(math.Abs(float64(R[i][j]))) + eps
		}
	}

	// Translation vector in world, expressed in A's basis
	tWorld := raylib.Vector3Subtract(b.Center, a.Center)
	t := [3]float32{
		raylib.Vector3DotProduct(tWorld, A[0]),
		raylib.Vector3DotProduct(tWorld, A[1]),
		raylib.Vector3DotProduct(tWorld, A[2]),
	}

	// 1) Test axes A0, A1, A2
	for i := 0; i < 3; i++ {
		ra := ae[i]
		rb := be[0]*AbsR[i][0] + be[1]*AbsR[i][1] + be[2]*AbsR[i][2]
		if float32(math.Abs(float64(t[i]))) > ra+rb {
			return false
		}
	}

	// 2) Test axes B0, B1, B2
	for j := 0; j < 3; j++ {
		ra := ae[0]*AbsR[0][j] + ae[1]*AbsR[1][j] + ae[2]*AbsR[2][j]
		rb := be[j]
		// t in B basis: tB = [t·B0, t·B1, t·B2] = [tA·Rcol]
		tB := t[0]*R[0][j] + t[1]*R[1][j] + t[2]*R[2][j]
		if float32(math.Abs(float64(tB))) > ra+rb {
			return false
		}
	}

	// 3) Test the 9 cross products Ai x Bj
	// Formulas from "OBBTree" / Gottschalk et al. (same as many engine implementations).
	for i := 0; i < 3; i++ {
		ip1 := (i + 1) % 3
		ip2 := (i + 2) % 3
		for j := 0; j < 3; j++ {
			jp1 := (j + 1) % 3
			jp2 := (j + 2) % 3

			ra := ae[ip1]*AbsR[ip2][j] + ae[ip2]*AbsR[ip1][j]
			rb := be[jp1]*AbsR[i][jp2] + be[jp2]*AbsR[i][jp1]

			// |t · (Ai x Bj)| in A basis:
			// = | t[ip2]*R[ip1][j] - t[ip1]*R[ip2][j] |
			val := float32(math.Abs(float64(t[ip2]*R[ip1][j] - t[ip1]*R[ip2][j])))

			if val > ra+rb {
				return false
			}
		}
	}

	return true
}

func CheckCollisionCapsuleOBB(capsule pub_object.Capsule, obb pub_object.OrientedBox) bool {
	obb = pub_object.OrientedBoxNormalizeScale(obb)

	closestOnSegment := ClosestPointOnSegment(
		ClosestPointOnOBB(capsule.Start, obb),
		capsule.Start,
		capsule.End,
	)
	closestOnBox := ClosestPointOnOBB(closestOnSegment, obb)
	distSq := raylib.Vector3DistanceSqr(closestOnBox, closestOnSegment)
	return distSq <= capsule.Radius*capsule.Radius
}

func CheckCollisionCapsuleBox(capsule pub_object.Capsule, box raylib.BoundingBox) bool {
	closestOnSegment := ClosestPointOnSegment(
		ClosestPointOnBox(capsule.Start, box),
		capsule.Start,
		capsule.End,
	)
	closestOnBox := ClosestPointOnBox(closestOnSegment, box)
	distSq := raylib.Vector3DistanceSqr(closestOnBox, closestOnSegment)
	return distSq <= (capsule.Radius * capsule.Radius)
}

func CheckCollisionCapsuleSphere(capsule pub_object.Capsule, center raylib.Vector3, radius float32) bool {
	closestPoint := ClosestPointOnSegment(center, capsule.Start, capsule.End)
	distSq := raylib.Vector3DistanceSqr(center, closestPoint)
	radiusSum := capsule.Radius + radius
	return distSq <= (radiusSum * radiusSum)
}

func CheckCollisionCapsules(a, b pub_object.Capsule) bool {
	dirA := raylib.Vector3Subtract(a.End, a.Start)
	dirB := raylib.Vector3Subtract(b.End, b.Start)
	r := raylib.Vector3Subtract(a.Start, b.Start)

	dotAA := raylib.Vector3DotProduct(dirA, dirA)
	dotBB := raylib.Vector3DotProduct(dirB, dirB)
	dotAB := raylib.Vector3DotProduct(dirA, dirB)
	dotAR := raylib.Vector3DotProduct(dirA, r)
	dotBR := raylib.Vector3DotProduct(dirB, r)

	denom := dotAA*dotBB - dotAB*dotAB
	var s, t float32

	if denom > 1e-6 {
		s = (dotAB*dotBR - dotBB*dotAR) / denom
		s = fmax(0, fmin(1, s))
		t = (dotAB*s + dotBR) / dotBB

		if t < 0 {
			t = 0
			s = fmax(0, fmin(1, -dotAR/dotAA))
		} else if t > 1 {
			t = 1
			s = fmax(0, fmin(1, (dotAB-dotAR)/dotAA))
		}
	} else {
		s = 0.5
		t = fmax(0, fmin(1, dotBR/dotBB))
	}

	closestA := raylib.Vector3Add(a.Start, raylib.Vector3Scale(dirA, s))
	closestB := raylib.Vector3Add(b.Start, raylib.Vector3Scale(dirB, t))

	distSq := raylib.Vector3DistanceSqr(closestA, closestB)
	radiusSum := a.Radius + b.Radius
	return distSq <= (radiusSum * radiusSum)
}

func CheckCollisionCapsuleMesh(capsule pub_object.Capsule, mesh *raylib.Mesh, transform raylib.Matrix) bool {
	if mesh == nil || mesh.Vertices == nil || mesh.VertexCount <= 0 || mesh.TriangleCount <= 0 {
		return false
	}

	axis := raylib.Vector3Subtract(capsule.End, capsule.Start)
	radiusSq := capsule.Radius * capsule.Radius

	triCount := meshTriangleCount(mesh)

	const samples = 5
	for i := 0; i < triCount; i++ {
		v0, v1, v2 := meshTrianglePositions(mesh, i)

		a := v0.Transform(transform)
		b := v1.Transform(transform)
		c := v2.Transform(transform)

		for s := 0; s < samples; s++ {
			tt := float32(s) / float32(samples-1)
			p := raylib.Vector3Add(capsule.Start, raylib.Vector3Scale(axis, tt))

			closest := ClosestPointOnTriangle(p, a, b, c)
			if raylib.Vector3LengthSqr(raylib.Vector3Subtract(closest, p)) <= radiusSq {
				return true
			}
		}
	}
	return false
}

func CheckPenetrationCapsuleBox(capsule pub_object.Capsule, box raylib.BoundingBox) Penetration {
	var result Penetration

	closestOnSegment := ClosestPointOnSegment(
		ClosestPointOnBox(capsule.Start, box),
		capsule.Start,
		capsule.End,
	)

	closestOnBox := ClosestPointOnBox(closestOnSegment, box)
	delta := raylib.Vector3Subtract(closestOnSegment, closestOnBox)

	distSq := raylib.Vector3LengthSqr(delta)
	radiusSq := capsule.Radius * capsule.Radius
	if distSq >= radiusSq {
		return result
	}

	dist := float32(math.Sqrt(float64(distSq)))
	result.Collides = true
	result.Depth = capsule.Radius - dist

	if dist > 1e-6 {
		result.Normal = raylib.Vector3Scale(delta, 1.0/dist)
	} else {
		boxCenter := raylib.Vector3{
			X: (box.Min.X + box.Max.X) * 0.5,
			Y: (box.Min.Y + box.Max.Y) * 0.5,
			Z: (box.Min.Z + box.Max.Z) * 0.5,
		}
		toCenter := raylib.Vector3Subtract(closestOnSegment, boxCenter)

		ax := float32(math.Abs(float64(toCenter.X)))
		ay := float32(math.Abs(float64(toCenter.Y)))
		az := float32(math.Abs(float64(toCenter.Z)))

		if ax >= ay && ax >= az {
			if toCenter.X > 0 {
				result.Normal = raylib.Vector3{X: 1, Y: 0, Z: 0}
			} else {
				result.Normal = raylib.Vector3{X: -1, Y: 0, Z: 0}
			}
		} else if ay >= az {
			if toCenter.Y > 0 {
				result.Normal = raylib.Vector3{X: 0, Y: 1, Z: 0}
			} else {
				result.Normal = raylib.Vector3{X: 0, Y: -1, Z: 0}
			}
		} else {
			if toCenter.Z > 0 {
				result.Normal = raylib.Vector3{X: 0, Y: 0, Z: 1}
			} else {
				result.Normal = raylib.Vector3{X: 0, Y: 0, Z: -1}
			}
		}
	}

	result.MTV = raylib.Vector3Scale(result.Normal, result.Depth)
	return result
}

func CheckPenetrationCapsuleSphere(capsule pub_object.Capsule, center raylib.Vector3, radius float32) Penetration {
	var result Penetration

	closestOnSegment := ClosestPointOnSegment(center, capsule.Start, capsule.End)
	delta := raylib.Vector3Subtract(center, closestOnSegment)

	distSq := raylib.Vector3LengthSqr(delta)
	combinedRadius := capsule.Radius + radius
	combinedRadiusSq := combinedRadius * combinedRadius
	if distSq >= combinedRadiusSq {
		return result
	}

	dist := float32(math.Sqrt(float64(distSq)))
	result.Collides = true
	result.Depth = combinedRadius - dist

	if dist > 1e-6 {
		result.Normal = raylib.Vector3Scale(delta, 1.0/dist)
	} else {
		capsuleDir := raylib.Vector3Subtract(capsule.End, capsule.Start)
		capsuleLengthSq := raylib.Vector3LengthSqr(capsuleDir)

		if capsuleLengthSq > 1e-6 {
			// perp in XY
			result.Normal = raylib.Vector3{X: capsuleDir.Y, Y: -capsuleDir.X, Z: 0}
			normalLengthSq := raylib.Vector3LengthSqr(result.Normal)

			if normalLengthSq < 1e-6 {
				// perp in YZ
				result.Normal = raylib.Vector3{X: 0, Y: capsuleDir.Z, Z: -capsuleDir.Y}
				normalLengthSq = raylib.Vector3LengthSqr(result.Normal)
			}

			if normalLengthSq > 1e-6 {
				result.Normal = raylib.Vector3Normalize(result.Normal)
			} else {
				result.Normal = raylib.Vector3{X: 0, Y: 1, Z: 0}
			}
		} else {
			result.Normal = raylib.Vector3{X: 0, Y: 1, Z: 0}
		}
	}

	result.MTV = raylib.Vector3Scale(result.Normal, result.Depth)
	return result
}

func CheckPenetrationCapsules(a, b pub_object.Capsule) Penetration {
	var result Penetration

	dirA := raylib.Vector3Subtract(a.End, a.Start)
	dirB := raylib.Vector3Subtract(b.End, b.Start)
	r := raylib.Vector3Subtract(a.Start, b.Start)

	dotAA := raylib.Vector3DotProduct(dirA, dirA)
	dotBB := raylib.Vector3DotProduct(dirB, dirB)
	dotAB := raylib.Vector3DotProduct(dirA, dirB)
	dotAR := raylib.Vector3DotProduct(dirA, r)
	dotBR := raylib.Vector3DotProduct(dirB, r)

	denom := dotAA*dotBB - dotAB*dotAB
	var s, t float32

	if denom > 1e-6 {
		s = (dotAB*dotBR - dotBB*dotAR) / denom
		s = fmax(0, fmin(1, s))
		t = (dotAB*s + dotBR) / dotBB

		if t < 0 {
			t = 0
			s = fmax(0, fmin(1, -dotAR/dotAA))
		} else if t > 1 {
			t = 1
			s = fmax(0, fmin(1, (dotAB-dotAR)/dotAA))
		}
	} else {
		s = 0.5
		t = fmax(0, fmin(1, dotBR/dotBB))
	}

	closestA := raylib.Vector3Add(a.Start, raylib.Vector3Scale(dirA, s))
	closestB := raylib.Vector3Add(b.Start, raylib.Vector3Scale(dirB, t))

	delta := raylib.Vector3Subtract(closestA, closestB)
	distSq := raylib.Vector3LengthSqr(delta)
	combinedRadius := a.Radius + b.Radius
	combinedRadiusSq := combinedRadius * combinedRadius
	if distSq >= combinedRadiusSq {
		return result
	}

	dist := float32(math.Sqrt(float64(distSq)))
	result.Collides = true
	result.Depth = combinedRadius - dist

	if dist > 1e-6 {
		result.Normal = raylib.Vector3Scale(delta, 1.0/dist)
	} else {
		cross := raylib.Vector3CrossProduct(dirA, dirB)
		crossLenSq := raylib.Vector3LengthSqr(cross)

		if crossLenSq > 1e-6 {
			result.Normal = raylib.Vector3Normalize(cross)
		} else {
			perp := raylib.Vector3{X: dirA.Y, Y: -dirA.X, Z: 0}
			perpLenSq := raylib.Vector3LengthSqr(perp)

			if perpLenSq < 1e-6 {
				perp = raylib.Vector3{X: 0, Y: dirA.Z, Z: -dirA.Y}
				perpLenSq = raylib.Vector3LengthSqr(perp)
			}

			if perpLenSq > 1e-6 {
				result.Normal = raylib.Vector3Normalize(perp)
			} else {
				result.Normal = raylib.Vector3{X: 0, Y: 1, Z: 0}
			}
		}
	}

	result.MTV = raylib.Vector3Scale(result.Normal, result.Depth)
	return result
}

//
// ========================================
// SWEEP TESTS
// ========================================
//

func SweepSpherePoint(center raylib.Vector3, radius float32, velocity raylib.Vector3, point raylib.Vector3) SweepCollision {
	var result SweepCollision

	m := raylib.Vector3Subtract(center, point)
	c := raylib.Vector3DotProduct(m, m) - radius*radius

	if c <= 0 {
		result.Hit = true
		result.Time = 0
		result.Point = raylib.Vector3Add(point, raylib.Vector3Scale(raylib.Vector3Normalize(m), radius))
		result.Normal = raylib.Vector3Normalize(m)
		return result
	}

	a := raylib.Vector3DotProduct(velocity, velocity)
	b := raylib.Vector3DotProduct(m, velocity)
	discr := b*b - a*c
	if discr < 0 {
		return result
	}

	t := (-b - float32(math.Sqrt(float64(discr)))) / a
	if t < 0 || t > 1 {
		return result
	}

	hit := raylib.Vector3Add(center, raylib.Vector3Scale(velocity, t))

	result.Hit = true
	result.Time = t
	result.Point = hit
	result.Normal = raylib.Vector3Normalize(raylib.Vector3Subtract(hit, point))
	return result
}

func SweepSphereSegment(center raylib.Vector3, radius float32, velocity raylib.Vector3, aPt raylib.Vector3, bPt raylib.Vector3) SweepCollision {
	var result SweepCollision

	d := raylib.Vector3Subtract(bPt, aPt)
	m := raylib.Vector3Subtract(center, aPt)

	dd := raylib.Vector3DotProduct(d, d)
	md := raylib.Vector3DotProduct(m, d)
	nd := raylib.Vector3DotProduct(velocity, d)

	a0 := dd*raylib.Vector3DotProduct(velocity, velocity) - nd*nd
	b0 := dd*raylib.Vector3DotProduct(m, velocity) - md*nd
	c0 := dd*(raylib.Vector3DotProduct(m, m)-radius*radius) - md*md

	if float32(math.Abs(float64(a0))) < 1e-8 {
		return result
	}

	discr := b0*b0 - a0*c0
	if discr < 0 {
		return result
	}

	t := (-b0 - float32(math.Sqrt(float64(discr)))) / a0
	if t < 0 || t > 1 {
		return result
	}

	s := (md + t*nd) / dd
	if s < 0 || s > 1 {
		return result
	}

	hit := raylib.Vector3Add(center, raylib.Vector3Scale(velocity, t))
	closest := raylib.Vector3Add(aPt, raylib.Vector3Scale(d, s))

	result.Hit = true
	result.Time = t
	result.Point = hit
	result.Normal = raylib.Vector3Normalize(raylib.Vector3Subtract(hit, closest))
	return result
}

func SweepSphereTrianglePlane(center raylib.Vector3, radius float32, velocity raylib.Vector3, a, b, c raylib.Vector3) SweepCollision {
	var result SweepCollision

	ab := raylib.Vector3Subtract(b, a)
	ac := raylib.Vector3Subtract(c, a)
	normal := raylib.Vector3Normalize(raylib.Vector3CrossProduct(ab, ac))

	dist := raylib.Vector3DotProduct(raylib.Vector3Subtract(center, a), normal)
	denom := raylib.Vector3DotProduct(velocity, normal)

	// Moving away or parallel
	if denom >= 0 {
		return result
	}

	t := (radius - dist) / denom
	if t < 0 || t > 1 {
		return result
	}

	hitPoint := raylib.Vector3Add(center, raylib.Vector3Scale(velocity, t))
	projected := raylib.Vector3Subtract(hitPoint, raylib.Vector3Scale(normal, radius))

	closest := ClosestPointOnTriangle(projected, a, b, c)
	d2 := raylib.Vector3LengthSqr(raylib.Vector3Subtract(projected, closest))
	if d2 > 1e-6 {
		return result
	}

	result.Hit = true
	result.Time = t
	result.Point = hitPoint
	result.Normal = normal
	return result
}

func SweepSphereTriangle(center raylib.Vector3, radius float32, velocity raylib.Vector3, a, b, c raylib.Vector3) SweepCollision {
	var result SweepCollision
	result.Time = 1

	faceHit := SweepSphereTrianglePlane(center, radius, velocity, a, b, c)
	if faceHit.Hit && faceHit.Time < result.Time {
		result = faceHit
	}

	edges := [3][2]raylib.Vector3{{a, b}, {b, c}, {c, a}}
	for i := 0; i < 3; i++ {
		edgeHit := SweepSphereSegment(center, radius, velocity, edges[i][0], edges[i][1])
		if edgeHit.Hit && edgeHit.Time < result.Time {
			result = edgeHit
		}
	}

	verts := [3]raylib.Vector3{a, b, c}
	for i := 0; i < 3; i++ {
		vertHit := SweepSpherePoint(center, radius, velocity, verts[i])
		if vertHit.Hit && vertHit.Time < result.Time {
			result = vertHit
		}
	}

	return result
}

func SweepSphereBox(center raylib.Vector3, radius float32, velocity raylib.Vector3, box raylib.BoundingBox) SweepCollision {
	var collision SweepCollision

	velocityLength := raylib.Vector3Length(velocity)
	if velocityLength < 1e-6 {
		return collision
	}

	expanded := raylib.BoundingBox{
		Min: raylib.Vector3Subtract(box.Min, raylib.Vector3{X: radius, Y: radius, Z: radius}),
		Max: raylib.Vector3Add(box.Max, raylib.Vector3{X: radius, Y: radius, Z: radius}),
	}

	ray := raylib.Ray{
		Position:  center,
		Direction: raylib.Vector3Scale(velocity, 1.0/velocityLength),
	}
	hit := raylib.GetRayCollisionBox(ray, expanded)

	if hit.Hit && hit.Distance <= velocityLength {
		collision.Hit = true
		collision.Time = hit.Distance / velocityLength
		collision.Point = hit.Point
		collision.Normal = hit.Normal
	}

	return collision
}

// SweepSphereMesh matches R3D_SweepSphereMesh (refactored to *raylib.Mesh)
func SweepSphereMesh(center raylib.Vector3, radius float32, velocity raylib.Vector3, mesh *raylib.Mesh, transform raylib.Matrix) SweepCollision {
	var result SweepCollision
	result.Time = 1

	if mesh == nil || mesh.Vertices == nil || mesh.VertexCount <= 0 || mesh.TriangleCount <= 0 {
		return result
	}

	triCount := meshTriangleCount(mesh)
	for i := 0; i < triCount; i++ {
		v0, v1, v2 := meshTrianglePositions(mesh, i)

		a := v0.Transform(transform)
		b := v1.Transform(transform)
		c := v2.Transform(transform)

		hit := SweepSphereTriangle(center, radius, velocity, a, b, c)
		if hit.Hit && hit.Time < result.Time {
			result = hit
		}
	}

	return result
}

func SweepCapsuleBox(capsule pub_object.Capsule, velocity raylib.Vector3, box raylib.BoundingBox) SweepCollision {
	var collision SweepCollision

	velocityLength := raylib.Vector3Length(velocity)
	if velocityLength < 1e-6 {
		return collision
	}

	expanded := raylib.BoundingBox{
		Min: raylib.Vector3Subtract(box.Min, raylib.Vector3{X: capsule.Radius, Y: capsule.Radius, Z: capsule.Radius}),
		Max: raylib.Vector3Add(box.Max, raylib.Vector3{X: capsule.Radius, Y: capsule.Radius, Z: capsule.Radius}),
	}

	velocityDir := raylib.Vector3Scale(velocity, 1.0/velocityLength)
	capsuleAxis := raylib.Vector3Subtract(capsule.End, capsule.Start)

	bestHit := raylib.RayCollision{Distance: float32(math.Inf(1))}
	foundHit := false

	const samples = 3
	for i := 0; i < samples; i++ {
		t := float32(i) / float32(samples-1)
		samplePoint := raylib.Vector3Add(capsule.Start, raylib.Vector3Scale(capsuleAxis, t))

		ray := raylib.Ray{Position: samplePoint, Direction: velocityDir}
		hit := raylib.GetRayCollisionBox(ray, expanded)

		if hit.Hit && hit.Distance <= velocityLength && hit.Distance < bestHit.Distance {
			bestHit = hit
			foundHit = true
		}
	}

	if foundHit {
		collision.Hit = true
		collision.Time = bestHit.Distance / velocityLength
		collision.Point = bestHit.Point
		collision.Normal = bestHit.Normal
	}

	return collision
}

// SweepCapsuleMesh matches R3D_SweepCapsuleMesh (refactored to *raylib.Mesh)
//
// Note: the original C assumes indexed mesh (loop i += 3 over indices).
// This Go port supports both indexed and non-indexed by using meshTrianglePositions.
func SweepCapsuleMesh(capsule pub_object.Capsule, velocity raylib.Vector3, mesh *raylib.Mesh, transform raylib.Matrix) SweepCollision {
	var result SweepCollision
	result.Time = 1

	if mesh == nil || mesh.Vertices == nil || mesh.VertexCount <= 0 || mesh.TriangleCount <= 0 {
		return result
	}

	triCount := meshTriangleCount(mesh)
	for i := 0; i < triCount; i++ {
		va, vb, vc := meshTrianglePositions(mesh, i)
		a := va.Transform(transform)
		b := vb.Transform(transform)
		c := vc.Transform(transform)

		// Face plane tests (capsule endpoints as spheres)
		faceHit := SweepSphereTrianglePlane(capsule.Start, capsule.Radius, velocity, a, b, c)
		if faceHit.Hit && faceHit.Time < result.Time {
			result = faceHit
		}
		faceHit = SweepSphereTrianglePlane(capsule.End, capsule.Radius, velocity, a, b, c)
		if faceHit.Hit && faceHit.Time < result.Time {
			result = faceHit
		}

		// Edge segments (treat capsule endpoint sphere sweep vs each edge)
		segHit := SweepSphereSegment(capsule.Start, capsule.Radius, velocity, a, b)
		if segHit.Hit && segHit.Time < result.Time {
			result = segHit
		}
		segHit = SweepSphereSegment(capsule.Start, capsule.Radius, velocity, b, c)
		if segHit.Hit && segHit.Time < result.Time {
			result = segHit
		}
		segHit = SweepSphereSegment(capsule.Start, capsule.Radius, velocity, c, a)
		if segHit.Hit && segHit.Time < result.Time {
			result = segHit
		}

		// Vertex tests (start)
		vertHit := SweepSpherePoint(capsule.Start, capsule.Radius, velocity, a)
		if vertHit.Hit && vertHit.Time < result.Time {
			result = vertHit
		}
		vertHit = SweepSpherePoint(capsule.Start, capsule.Radius, velocity, b)
		if vertHit.Hit && vertHit.Time < result.Time {
			result = vertHit
		}
		vertHit = SweepSpherePoint(capsule.Start, capsule.Radius, velocity, c)
		if vertHit.Hit && vertHit.Time < result.Time {
			result = vertHit
		}

		// Vertex tests (end)
		vertHit = SweepSpherePoint(capsule.End, capsule.Radius, velocity, a)
		if vertHit.Hit && vertHit.Time < result.Time {
			result = vertHit
		}
		vertHit = SweepSpherePoint(capsule.End, capsule.Radius, velocity, b)
		if vertHit.Hit && vertHit.Time < result.Time {
			result = vertHit
		}
		vertHit = SweepSpherePoint(capsule.End, capsule.Radius, velocity, c)
		if vertHit.Hit && vertHit.Time < result.Time {
			result = vertHit
		}
	}

	return result
}
