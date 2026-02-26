package object

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	_ "time"
	"unsafe"

	"karalis/internal/camera"
	"karalis/internal/shader"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Grass struct {
	parent  pub_object.Object
	compute shader.Compute
	shader  shader.Shader

	visibleSSBO uint32
	counterSSBO uint32

	MaxVisible int32
	radius     float32
	cellSize   float32

	mesh raylib.Mesh
	mat  raylib.Material
	tex  raylib.Texture2D
}

// BladeGPU is std430-friendly if we keep it simple.
// Here: position.xyz + "height" in w.
// (std430: vec4 aligned to 16 bytes -> maps nicely to 4 float32)
type BladeGPU struct {
	Px, Py, Pz, H float32
}

func fillBlades(blades []BladeGPU) {
	radius := float32(25)
	for i := range blades {
		// Polar sampling
		a := rand.Float32() * 2 * math.Pi
		r := radius * float32(math.Sqrt(float64(rand.Float32())))
		x := r * float32(math.Cos(float64(a)))
		z := r * float32(math.Sin(float64(a)))

		h := 0.35 + rand.Float32()*0.75

		blades[i] = BladeGPU{
			Px: x,
			Py: 0,
			Pz: z,
			H:  h,
		}
	}
}

func NewGrass(parent pub_object.Object) (g *Grass, err error) {
	g = &Grass{parent: parent}
	err = g.Init()
	return g, err
}

func (g *Grass) Init() error {
	if g == nil {
		return fmt.Errorf("Invalid grass")
	}
	// Worst-case visible blades we will store in the "visible" SSBO.
	// If you exceed this, compute will clamp.
	g.MaxVisible = int32(150_000)
	g.radius = float32(300)
	g.cellSize = float32(5)

	// --- SSBOs ---
	// Output visible blades buffer (binding=1 in compute; also binding=1 in vertex shader)
	emptyVisible := make([]BladeGPU, g.MaxVisible)
	g.visibleSSBO = raylib.LoadShaderBuffer(
		uint32(len(emptyVisible))*uint32(unsafe.Sizeof(emptyVisible[0])),
		unsafe.Pointer(&emptyVisible[0]),
		int32(raylib.DynamicDraw),
	)

	// Counter buffer (binding=2): uint visibleCount
	var zero uint32 = 0
	g.counterSSBO = raylib.LoadShaderBuffer(
		uint32(unsafe.Sizeof(zero)),
		unsafe.Pointer(&zero),
		int32(raylib.DynamicDraw),
	)

	g.mesh = raylib.GenMeshPlane(1, 1, 1, 1)
	raylib.UploadMesh(&g.mesh, false)

	g.compute = shader.Compute{}
	g.compute.Init("grass")
	g.shader = shader.Shader{}
	g.shader.Init("grass")

	g.mat = raylib.LoadMaterialDefault()
	g.mat.Shader = *g.shader.GetShader()

	return nil
}

func (g *Grass) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if g == nil {
		return cmds
	}

	return cmds
}

func (g *Grass) Render(cam *camera.Cam) []func() {
	cmds := []func(){}
	if g == nil {
		return cmds
	}
	zero := uint32(0)
	raylib.UpdateShaderBuffer(g.counterSSBO, unsafe.Pointer(&zero), 4, 0)

	g.compute.Begin()
	raylib.BindShaderBuffer(g.visibleSSBO, 1)
	raylib.BindShaderBuffer(g.counterSSBO, 2)
	// Upload uniforms to compute
	err := g.compute.SetUniform("uCameraPos", cam.GetPos())
	if err != nil {
		fmt.Printf("(compute)uCameraPos: %+v\n", err)
	}
	err = g.compute.SetUniform("uMaxVisible", []int32{g.MaxVisible})
	if err != nil {
		fmt.Printf("(compute)uMaxVisible: %+v\n", err)
	}
	err = g.compute.SetUniform("uSeed", []int32{123456})
	if err != nil {
		fmt.Printf("(compute)uSeed: %+v\n", err)
	}
	err = g.compute.SetUniform("uRadius", []float32{g.radius})
	if err != nil {
		fmt.Printf("(compute)uRadius: %+v\n", err)
	}
	err = g.compute.SetUniform("uCellSize", []float32{g.cellSize})
	if err != nil {
		fmt.Printf("(compute)uCellSize: %+v\n", err)
	}

	rCells := int(math.Ceil(float64(g.radius / g.cellSize)))
	size := rCells * 2

	gx := uint32((size + 15) / 16)
	gy := uint32((size + 15) / 16)
	raylib.ComputeShaderDispatch(gx, gy, 1)
	g.compute.End()

	var visibleCount uint32
	raylib.ReadShaderBuffer(g.counterSSBO, unsafe.Pointer(&visibleCount), 4, 0)
	raylib.BindShaderBuffer(g.visibleSSBO, 1)
	err = g.shader.SetUniform("uRot", raylib.MatrixRotate(raylib.Vector3{1, 0, 0}, math.Pi/2))
	if err != nil {
		fmt.Printf("(shader)uRot: %+v\n", err)
	}

	instances := int(visibleCount)
	if instances > 0 {
		mdl := g.GetModelMatrix()
		svec := g.GetScale()
		scale := raylib.MatrixScale(1/svec.X*0.2, 1/svec.Y*0.2, 1/svec.Z*0.2)
		mat := raylib.MatrixMultiply(mdl, scale)
		transforms := make([]raylib.Matrix, instances)
		for i := range transforms {
			transforms[i] = mat
		}
		raylib.DrawMeshInstanced(g.mesh, g.mat, transforms, instances)
	}

	return cmds
}

func (g *Grass) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if g == nil {
		return cmds
	}

	return cmds
}

func (g *Grass) Update(dt float32) {
	if g == nil {
		return
	}
}

func (g *Grass) GetCollider() pub_object.Collider {
	if g == nil || g.parent == nil {
		return nil
	}

	return g.parent.GetCollider()
}

func (g *Grass) OnAdd() {
	if g == nil {
		return
	}
}

func (g *Grass) OnRemove() {
	if g == nil {
		return
	}
}

func (g *Grass) AddChild(obj pub_object.Object) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.AddChild(obj)
}

func (g *Grass) RemChild(obj pub_object.Object) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.RemChild(obj)
}

func (g *Grass) GetChilds() []pub_object.Object {
	if g == nil || g.parent == nil {
		return []pub_object.Object{}
	}

	return []pub_object.Object{}
}

func (g *Grass) GetModelMatrix() raylib.Matrix {
	if g == nil || g.parent == nil {
		return raylib.Matrix{}
	}

	return g.parent.GetModelMatrix()
}

func (g *Grass) GetModel() *raylib.Model {
	if g == nil || g.parent == nil {
		return nil
	}

	return nil
}

func (g *Grass) SetColor(col color.Color) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.SetColor(col)
}

func (g *Grass) GetColor() color.Color {
	if g == nil || g.parent == nil {
		return nil
	}

	return g.parent.GetColor()
}

func (g *Grass) GetScale() raylib.Vector3 {
	if g == nil || g.parent == nil {
		return raylib.Vector3{}
	}

	return g.parent.GetScale()
}

func (g *Grass) SetScale(sc raylib.Vector3) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.SetScale(sc)
}

func (g *Grass) SetPos(pos raylib.Vector3) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.SetPos(pos)
}

func (g *Grass) GetPos() raylib.Vector3 {
	if g == nil || g.parent == nil {
		return raylib.Vector3{}
	}

	return g.parent.GetPos()
}

func (g *Grass) GetPitch() float32 {
	if g == nil || g.parent == nil {
		return 0
	}

	return g.parent.GetPitch()
}

func (g *Grass) SetPitch(pitch float32) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.SetPitch(pitch)
}

func (g *Grass) GetYaw() float32 {
	if g == nil || g.parent == nil {
		return 0
	}

	return g.parent.GetYaw()
}

func (g *Grass) SetYaw(yaw float32) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.SetYaw(yaw)
}

func (g *Grass) GetRoll() float32 {
	if g == nil || g.parent == nil {
		return 0
	}

	return g.parent.GetRoll()
}

func (g *Grass) SetRoll(roll float32) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.SetRoll(roll)
}

func (g *Grass) GetVertices() []raylib.Vector3 {
	if g == nil || g.parent == nil {
		return []raylib.Vector3{}
	}

	return g.parent.GetVertices()
}

func (g *Grass) GetUVs() []raylib.Vector2 {
	if g == nil || g.parent == nil {
		return []raylib.Vector2{}
	}

	return g.parent.GetUVs()
}

func (g *Grass) SetUVs(uvs []raylib.Vector2) {
	if g == nil || g.parent == nil {
		return
	}
	g.parent.SetUVs(uvs)
}

func (g *Grass) GetMaterials() *raylib.Material {
	if g == nil {
		return nil
	}

	return &g.mat
}

func (g *Grass) SetTexture(tex raylib.Texture2D) {
	if g == nil {
		return
	}

	g.tex = tex
}

func (g *Grass) GetTexture() *raylib.Texture2D {
	if g == nil {
		return nil
	}

	return &g.tex
}
