package physics

import (
	"unsafe"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func meshVerticesXYZ(mesh *raylib.Mesh) []float32 {
	if mesh == nil || mesh.Vertices == nil || mesh.VertexCount <= 0 {
		return nil
	}
	n := int(mesh.VertexCount) * 3
	return unsafe.Slice(mesh.Vertices, n)
}

func meshIndices(mesh *raylib.Mesh) []uint16 {
	if mesh == nil || mesh.Indices == nil || mesh.TriangleCount <= 0 {
		return nil
	}
	n := int(mesh.TriangleCount) * 3
	return unsafe.Slice(mesh.Indices, n)
}

func meshHasIndices(mesh *raylib.Mesh) bool {
	return mesh != nil && mesh.Indices != nil
}

func meshTriangleCount(mesh *raylib.Mesh) int {
	if mesh == nil {
		return 0
	}
	// raylib already gives triangle count; trust it.
	return int(mesh.TriangleCount)
}

func meshPositionAt(mesh *raylib.Mesh, vertexIndex int) raylib.Vector3 {
	v := meshVerticesXYZ(mesh)
	base := vertexIndex * 3
	return raylib.Vector3{X: v[base], Y: v[base+1], Z: v[base+2]}
}

func meshTrianglePositions(mesh *raylib.Mesh, tri int) (a, b, c raylib.Vector3) {
	base := tri * 3
	if meshHasIndices(mesh) {
		idx := meshIndices(mesh)
		i0 := int(idx[base])
		i1 := int(idx[base+1])
		i2 := int(idx[base+2])
		return meshPositionAt(mesh, i0), meshPositionAt(mesh, i1), meshPositionAt(mesh, i2)
	}
	// Non-indexed triangle list: vertices laid out as triples
	return meshPositionAt(mesh, base), meshPositionAt(mesh, base+1), meshPositionAt(mesh, base+2)
}
