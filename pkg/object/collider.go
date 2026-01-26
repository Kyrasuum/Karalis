package object

import (
	"math"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type OBB struct {
	Center raylib.Vector3
	Axis   [3]raylib.Vector3 // normalized
	Half   raylib.Vector3    // half sizes along each axis (x,y,z)
}

type BoundingSphere struct {
	Center raylib.Vector3
	Radius float32
}

type CollisionData struct {
	Obj1 Object
	Obj2 Object
}

type Collider interface {
	GetObj() Object
	GetBoundingSphere() BoundingSphere
	GetAABB() raylib.BoundingBox
	GetOOBB() OBB
	GetCollidable() []Object
	Collide(CollisionData)
	Update(dt float32)
	RegHandler(string, interface{})
	GetTouching() []Object
}

func CheckCollisionSpheres(sp1 BoundingSphere, sp2 BoundingSphere) bool {
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

func CheckCollisionOOBB(a OBB, b OBB) bool {
	// Convenience
	A := a.Axis
	B := b.Axis
	EA := a.Half
	EB := b.Half

	// Rotation matrix expressing B in A’s frame: R[i][j] = dot(Ai, Bj)
	var R [3][3]float32
	var AbsR [3][3]float32

	const eps = 1e-6 // helps stability when axes are nearly parallel

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			R[i][j] = raylib.Vector3DotProduct(A[i], B[j])
			AbsR[i][j] = float32(math.Abs(float64(R[i][j]))) + eps
		}
	}

	// Translation vector from A to B, in world, then expressed in A’s frame
	tWorld := raylib.Vector3Subtract(b.Center, a.Center)
	t := raylib.NewVector3(
		raylib.Vector3DotProduct(tWorld, A[0]),
		raylib.Vector3DotProduct(tWorld, A[1]),
		raylib.Vector3DotProduct(tWorld, A[2]),
	)

	// Helper for extents
	EAi := [3]float32{EA.X, EA.Y, EA.Z}
	EBi := [3]float32{EB.X, EB.Y, EB.Z}
	ti := [3]float32{t.X, t.Y, t.Z}

	// 1) Test axes L = A0, A1, A2
	for i := 0; i < 3; i++ {
		ra := EAi[i]
		rb := EBi[0]*AbsR[i][0] + EBi[1]*AbsR[i][1] + EBi[2]*AbsR[i][2]
		if float32(math.Abs(float64(ti[i]))) > ra+rb {
			return false
		}
	}

	// 2) Test axes L = B0, B1, B2
	for j := 0; j < 3; j++ {
		ra := EAi[0]*AbsR[0][j] + EAi[1]*AbsR[1][j] + EAi[2]*AbsR[2][j]
		rb := EBi[j]
		// projection of t onto Bj is dot(tWorld, Bj) = sum_i t_i * R[i][j]
		tproj := ti[0]*R[0][j] + ti[1]*R[1][j] + ti[2]*R[2][j]
		if float32(math.Abs(float64(tproj))) > ra+rb {
			return false
		}
	}

	// 3) Test axes L = Ai x Bj (9 tests)
	// A0 x B0
	{
		ra := EAi[1]*AbsR[2][0] + EAi[2]*AbsR[1][0]
		rb := EBi[1]*AbsR[0][2] + EBi[2]*AbsR[0][1]
		val := float32(math.Abs(float64(ti[2]*R[1][0] - ti[1]*R[2][0])))
		if val > ra+rb {
			return false
		}
	}
	// A0 x B1
	{
		ra := EAi[1]*AbsR[2][1] + EAi[2]*AbsR[1][1]
		rb := EBi[0]*AbsR[0][2] + EBi[2]*AbsR[0][0]
		val := float32(math.Abs(float64(ti[2]*R[1][1] - ti[1]*R[2][1])))
		if val > ra+rb {
			return false
		}
	}
	// A0 x B2
	{
		ra := EAi[1]*AbsR[2][2] + EAi[2]*AbsR[1][2]
		rb := EBi[0]*AbsR[0][1] + EBi[1]*AbsR[0][0]
		val := float32(math.Abs(float64(ti[2]*R[1][2] - ti[1]*R[2][2])))
		if val > ra+rb {
			return false
		}
	}

	// A1 x B0
	{
		ra := EAi[0]*AbsR[2][0] + EAi[2]*AbsR[0][0]
		rb := EBi[1]*AbsR[1][2] + EBi[2]*AbsR[1][1]
		val := float32(math.Abs(float64(ti[0]*R[2][0] - ti[2]*R[0][0])))
		if val > ra+rb {
			return false
		}
	}
	// A1 x B1
	{
		ra := EAi[0]*AbsR[2][1] + EAi[2]*AbsR[0][1]
		rb := EBi[0]*AbsR[1][2] + EBi[2]*AbsR[1][0]
		val := float32(math.Abs(float64(ti[0]*R[2][1] - ti[2]*R[0][1])))
		if val > ra+rb {
			return false
		}
	}
	// A1 x B2
	{
		ra := EAi[0]*AbsR[2][2] + EAi[2]*AbsR[0][2]
		rb := EBi[0]*AbsR[1][1] + EBi[1]*AbsR[1][0]
		val := float32(math.Abs(float64(ti[0]*R[2][2] - ti[2]*R[0][2])))
		if val > ra+rb {
			return false
		}
	}

	// A2 x B0
	{
		ra := EAi[0]*AbsR[1][0] + EAi[1]*AbsR[0][0]
		rb := EBi[1]*AbsR[2][2] + EBi[2]*AbsR[2][1]
		val := float32(math.Abs(float64(ti[1]*R[0][0] - ti[0]*R[1][0])))
		if val > ra+rb {
			return false
		}
	}
	// A2 x B1
	{
		ra := EAi[0]*AbsR[1][1] + EAi[1]*AbsR[0][1]
		rb := EBi[0]*AbsR[2][2] + EBi[2]*AbsR[2][0]
		val := float32(math.Abs(float64(ti[1]*R[0][1] - ti[0]*R[1][1])))
		if val > ra+rb {
			return false
		}
	}
	// A2 x B2
	{
		ra := EAi[0]*AbsR[1][2] + EAi[1]*AbsR[0][2]
		rb := EBi[0]*AbsR[2][1] + EBi[1]*AbsR[2][0]
		val := float32(math.Abs(float64(ti[1]*R[0][2] - ti[0]*R[1][2])))
		if val > ra+rb {
			return false
		}
	}

	return true
}
