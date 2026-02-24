package physics

import (
	"math"

	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

//
// ========================================
// RESOLUTION
// ========================================
//

func DepenetrateSphereBox(center *raylib.Vector3, radius float32, box raylib.BoundingBox, outPenetration *float32) bool {
	closestPoint := ClosestPointOnBox(*center, box)
	delta := raylib.Vector3Subtract(*center, closestPoint)
	distSq := raylib.Vector3LengthSqr(delta)
	radiusSq := radius * radius

	if distSq >= radiusSq {
		return false
	}

	dist := float32(math.Sqrt(float64(distSq)))
	penetration := radius - dist

	var direction raylib.Vector3
	if dist > 1e-6 {
		direction = raylib.Vector3Scale(delta, 1.0/dist)
	} else {
		direction = raylib.Vector3{X: 0, Y: 1, Z: 0}
	}

	*center = raylib.Vector3Add(*center, raylib.Vector3Scale(direction, penetration))
	if outPenetration != nil {
		*outPenetration = penetration
	}
	return true
}

func DepenetrateCapsuleBox(capsule *pub_object.Capsule, box raylib.BoundingBox, outPenetration *float32) bool {
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
		return false
	}

	dist := float32(math.Sqrt(float64(distSq)))
	penetration := capsule.Radius - dist

	var direction raylib.Vector3
	if dist > 1e-6 {
		direction = raylib.Vector3Scale(delta, 1.0/dist)
	} else {
		direction = raylib.Vector3{X: 0, Y: 1, Z: 0}
	}

	correction := raylib.Vector3Scale(direction, penetration)
	capsule.Start = raylib.Vector3Add(capsule.Start, correction)
	capsule.End = raylib.Vector3Add(capsule.End, correction)

	if outPenetration != nil {
		*outPenetration = penetration
	}
	return true
}

func SlideVelocity(velocity, normal raylib.Vector3) raylib.Vector3 {
	dot := raylib.Vector3DotProduct(velocity, normal)
	return raylib.Vector3Subtract(velocity, raylib.Vector3Scale(normal, dot))
}

func BounceVelocity(velocity, normal raylib.Vector3, bounciness float32) raylib.Vector3 {
	dot := raylib.Vector3DotProduct(velocity, normal)
	reflection := raylib.Vector3Subtract(velocity, raylib.Vector3Scale(normal, 2*dot))
	return raylib.Vector3Scale(reflection, bounciness)
}

func SlideSphereBox(center raylib.Vector3, radius float32, velocity raylib.Vector3, box raylib.BoundingBox, outNormal *raylib.Vector3) raylib.Vector3 {
	collision := SweepSphereBox(center, radius, velocity, box)
	if !collision.Hit {
		if outNormal != nil {
			*outNormal = raylib.Vector3{}
		}
		return velocity
	}

	if outNormal != nil {
		*outNormal = collision.Normal
	}

	safeTime := fmax(0, collision.Time-0.001)
	safeVelocity := raylib.Vector3Scale(velocity, safeTime)
	remainingVelocity := raylib.Vector3Scale(velocity, 1-safeTime)
	slidedRemaining := SlideVelocity(remainingVelocity, collision.Normal)

	return raylib.Vector3Add(safeVelocity, slidedRemaining)
}

func SlideSphereMesh(center raylib.Vector3, radius float32, velocity raylib.Vector3, mesh *raylib.Mesh, transform raylib.Matrix, outNormal *raylib.Vector3) raylib.Vector3 {
	collision := SweepSphereMesh(center, radius, velocity, mesh, transform)
	if !collision.Hit {
		if outNormal != nil {
			*outNormal = raylib.Vector3{}
		}
		return velocity
	}

	if outNormal != nil {
		*outNormal = collision.Normal
	}

	safeTime := fmax(0, collision.Time-0.001)
	safeVelocity := raylib.Vector3Scale(velocity, safeTime)
	remainingVelocity := raylib.Vector3Scale(velocity, 1-safeTime)
	slidedRemaining := SlideVelocity(remainingVelocity, collision.Normal)

	return raylib.Vector3Add(safeVelocity, slidedRemaining)
}

func SlideCapsuleBox(capsule pub_object.Capsule, velocity raylib.Vector3, box raylib.BoundingBox, outNormal *raylib.Vector3) raylib.Vector3 {
	collision := SweepCapsuleBox(capsule, velocity, box)
	if !collision.Hit {
		if outNormal != nil {
			*outNormal = raylib.Vector3{}
		}
		return velocity
	}

	if outNormal != nil {
		*outNormal = collision.Normal
	}

	safeTime := fmax(0, collision.Time-0.001)
	safeVelocity := raylib.Vector3Scale(velocity, safeTime)
	remainingVelocity := raylib.Vector3Scale(velocity, 1-safeTime)
	slidedRemaining := SlideVelocity(remainingVelocity, collision.Normal)

	return raylib.Vector3Add(safeVelocity, slidedRemaining)
}

func SlideCapsuleMesh(capsule pub_object.Capsule, velocity raylib.Vector3, mesh *raylib.Mesh, transform raylib.Matrix, outNormal *raylib.Vector3) raylib.Vector3 {
	collision := SweepCapsuleMesh(capsule, velocity, mesh, transform)
	if !collision.Hit {
		if outNormal != nil {
			*outNormal = raylib.Vector3{}
		}
		return velocity
	}

	if outNormal != nil {
		*outNormal = collision.Normal
	}

	safeTime := fmax(0, collision.Time-0.001)
	safeVelocity := raylib.Vector3Scale(velocity, safeTime)
	remainingVelocity := raylib.Vector3Scale(velocity, 1-safeTime)
	slidedRemaining := SlideVelocity(remainingVelocity, collision.Normal)

	return raylib.Vector3Add(safeVelocity, slidedRemaining)
}
