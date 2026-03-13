package rlx

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadModel(fileName string) rl.Model {
	return Call(func() rl.Model {
		return rl.LoadModel(fileName)
	})
}

func LoadModelFromMesh(mesh rl.Mesh) rl.Model {
	return Call(func() rl.Model {
		return rl.LoadModelFromMesh(mesh)
	})
}

func IsModelValid(model rl.Model) bool {
	return Call(func() bool {
		return rl.IsModelValid(model)
	})
}

func UnloadModel(model rl.Model) {
	Do(func() {
		rl.UnloadModel(model)
	})
}

func GetModelBoundingBox(model rl.Model) rl.BoundingBox {
	return Call(func() rl.BoundingBox {
		return rl.GetModelBoundingBox(model)
	})
}

func DrawModel(model rl.Model, position rl.Vector3, scale float32, tint color.RGBA) {
	Do(func() {
		rl.DrawModel(model, position, scale, tint)
	})
}

func DrawModelEx(model rl.Model, position rl.Vector3, rotationAxis rl.Vector3, rotationAngle float32, scale rl.Vector3, tint color.RGBA) {
	Do(func() {
		rl.DrawModelEx(model, position, rotationAxis, rotationAngle, scale, tint)
	})
}

func DrawModelWires(model rl.Model, position rl.Vector3, scale float32, tint color.RGBA) {
	Do(func() {
		rl.DrawModelWires(model, position, scale, tint)
	})
}

func DrawBoundingBox(box rl.BoundingBox, col color.RGBA) {
	Do(func() {
		rl.DrawBoundingBox(box, col)
	})
}

func UploadMesh(mesh *rl.Mesh, dynamic bool) {
	Do(func() {
		rl.UploadMesh(mesh, dynamic)
	})
}

func UpdateMeshBuffer(mesh rl.Mesh, index int, data []byte, offset int) {
	buf := append([]byte(nil), data...)
	Do(func() {
		rl.UpdateMeshBuffer(mesh, index, buf, offset)
	})
}

func UnloadMesh(mesh *rl.Mesh) {
	Do(func() {
		rl.UnloadMesh(mesh)
	})
}

func DrawMesh(mesh rl.Mesh, material rl.Material, transform rl.Matrix) {
	Do(func() {
		rl.DrawMesh(mesh, material, transform)
	})
}

func DrawMeshInstanced(mesh rl.Mesh, material rl.Material, transforms []rl.Matrix, instances int) {
	buf := append([]rl.Matrix(nil), transforms...)
	Do(func() {
		rl.DrawMeshInstanced(mesh, material, buf, instances)
	})
}

func DrawGrid(slices int32, spacing float32) {
	Do(func() {
		rl.DrawGrid(slices, spacing)
	})
}

func GenMeshCone(radius float32, height float32, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshCone(radius, height, slices)
	})
}

func GenMeshConeData(radius float32, height float32, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshConeData(radius, height, slices)
	})
}

func GenMeshCube(width float32, height float32, length float32) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshCube(width, height, length)
	})
}

func GenMeshCubeData(width float32, height float32, length float32) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshCubeData(width, height, length)
	})
}

func GenMeshCylinder(radius float32, height float32, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshCylinder(radius, height, slices)
	})
}

func GenMeshCylinderData(radius float32, height float32, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshCylinderData(radius, height, slices)
	})
}

func GenMeshHemiSphere(radius float32, rings int, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshHemiSphere(radius, rings, slices)
	})
}

func GenMeshHemiSphereData(radius float32, rings int, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshHemiSphereData(radius, rings, slices)
	})
}

func GenMeshPoly(sides int, radius float32) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshPoly(sides, radius)
	})
}

func GenMeshPolyData(sides int, radius float32) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshPolyData(sides, radius)
	})
}

func GenMeshSphere(radius float32, rings int, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshSphere(radius, rings, slices)
	})
}

func GenMeshSphereData(radius float32, rings int, slices int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshSphereData(radius, rings, slices)
	})
}

func GenMeshPlane(width float32, length float32, resX int, resZ int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshPlane(width, length, resX, resZ)
	})
}

func GenMeshPlaneData(width float32, length float32, resX int, resZ int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshPlaneData(width, length, resX, resZ)
	})
}

func GenMeshPlaneEx(origin, axisU, axisV rl.Vector3, resU, resV int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshPlaneEx(origin, axisU, axisV, resU, resV)
	})
}

func GenMeshPlaneExData(origin, axisU, axisV rl.Vector3, resU, resV int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshPlaneExData(origin, axisU, axisV, resU, resV)
	})
}

func GenMeshTorus(radius float32, size float32, radSeg int, sides int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshTorus(radius, size, radSeg, sides)
	})
}

func GenMeshTorusData(radius float32, size float32, radSeg int, sides int) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshTorusData(radius, size, radSeg, sides)
	})
}

func SetMaterialTexture(material *rl.Material, mapType int32, texture rl.Texture2D) {
	Do(func() {
		rl.SetMaterialTexture(material, mapType, texture)
	})
}

func GenMeshHeightmap(heightmap rl.Image, size rl.Vector3) rl.Mesh {
	return Call(func() rl.Mesh {
		return rl.GenMeshHeightmap(heightmap, size)
	})
}
