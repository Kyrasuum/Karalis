package object

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type BoundingSphere struct {
	Center raylib.Vector3
	Radius float32
}

type CollisionData struct {
	Obj1 Object
	Obj2 Object
}

type Collider interface {
	GetBoundingSphere() BoundingSphere
	GetBoundingBox() raylib.BoundingBox
	GetCollidable() []Object
	Collide(CollisionData)
	Update(dt float32)
	RegHandler(string, interface{})
	GetTouching() []Object
}

func CheckCollisionSpheres(sp1 BoundingSphere, sp2 BoundingSphere) bool {
	return raylib.Vector3Length(raylib.Vector3Subtract(sp1.Center, sp2.Center)) < sp1.Radius+sp2.Radius
}
