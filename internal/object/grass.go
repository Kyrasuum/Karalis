package object

import (
	"unsafe"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Grass struct {
	program      uint32
	ssbo         uint32
	maxInstances uint32

	instances []raylib.Matrix // cached for DrawMeshInstanced
	bladeMesh raylib.Mesh
	bladeMat  raylib.Material
}

func (g *Grass) InitCompute(compCode string, maxInstances uint32) {
	g.maxInstances = maxInstances

	sh := raylib.CompileShader(compCode, int32(raylib.ComputeShader))
	g.program = raylib.LoadComputeShaderProgram(sh)

	// SSBO layout: uint count + maxInstances * instance
	// instance = 2 vec4 = 8 floats = 32 bytes
	ssboSize := uint32(4) + maxInstances*32
	g.ssbo = raylib.LoadShaderBuffer(ssboSize, nil, int32(raylib.DynamicCopy))

	// Create a simple grass blade mesh (start with a vertical quad or 3 crossed quads)
	// (You can replace this with your "3 flat sprites" later.)
	g.bladeMesh = raylib.GenMeshPlane(0.08, 0.30, 1, 1)
	// rotate plane upright in your vertex shader or bake vertices (either way works)
	g.bladeMat = raylib.LoadMaterialDefault()
}

func (g *Grass) GenerateOnce(t *Terrain, mask raylib.Texture2D, terrainSize raylib.Vector3) {
	// reset counter to 0
	var zero uint32 = 0
	raylib.UpdateShaderBuffer(g.ssbo, unsafe.Pointer(&zero), 4, 0)

	// bind images for imageLoad (heightImg binding=0, maskImg binding=1)
	raylib.BindImageTexture(t.hgt.ID, 0, int32(raylib.UncompressedR8g8b8a8), true)
	raylib.BindImageTexture(mask.ID, 1, int32(raylib.UncompressedR8g8b8a8), true)

	// bind SSBO at binding=2 (matches shader)
	raylib.EnableShader(g.program)
	raylib.BindShaderBuffer(g.ssbo, 2)

	// TODO: set uniforms (heightSize, origin, size, density...)
	// If you stick with rlgl helpers only, easiest path is:
	// - encode params into a tiny params SSBO at binding=3
	// - read them in compute
	// (raylib's compute example avoids uniforms and uses SSBOs heavily) :contentReference[oaicite:5]{index=5}

	// dispatch
	w := uint32(t.hgt.Width)
	h := uint32(t.hgt.Height)
	gx := (w + 15) / 16
	gy := (h + 15) / 16
	raylib.ComputeShaderDispatch(gx, gy, 1)
	raylib.DisableShader()

	// read back count
	var count uint32
	raylib.ReadShaderBuffer(g.ssbo, unsafe.Pointer(&count), 4, 0)

	// read back instances (pos_scale + nrm_yaw)
	type inst struct {
		PosScale [4]float32
		NrmYaw   [4]float32
	}
	cpu := make([]inst, count)
	if count > 0 {
		raylib.ReadShaderBuffer(g.ssbo, unsafe.Pointer(&cpu[0]), count*32, 4)
	}

	// build matrices on CPU (position + align to normal + yaw + scale)
	g.instances = g.instances[:0]
	for i := 0; i < int(count); i++ {
		p := cpu[i].PosScale
		n := cpu[i].NrmYaw
		pos := raylib.NewVector3(p[0], p[1], p[2])
		sc := p[3]
		yaw := n[3]
		normal := raylib.NewVector3(n[0], n[1], n[2])

		mrot := raylib.MatrixRotate(normal, yaw)
		mpos := raylib.MatrixTranslate(pos.X, pos.Y, pos.Z)
		msca := raylib.MatrixScale(sc, sc, sc)
		m := raylib.MatrixMultiply(raylib.MatrixMultiply(msca, mrot), mpos)

		g.instances = append(g.instances, m)
	}
}

func (g *Grass) Draw(modelMatrix raylib.Matrix) {
	if len(g.instances) == 0 {
		return
	}
	// NOTE: DrawMeshInstanced uses per-instance matrices (fast) :contentReference[oaicite:6]{index=6}
	raylib.DrawMeshInstanced(g.bladeMesh, g.bladeMat, g.instances, len(g.instances))
}
