package world

import (
	"fmt"
	"image/color"
	"math"
	"unsafe"

	"karalis/internal/camera"
	"karalis/internal/shader"
	pub_object "karalis/pkg/object"
	"karalis/pkg/rng"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var (
	GrassLOD = 30.0
)

type Grass struct {
	parent  pub_object.Object
	compute shader.Compute
	shader  shader.Shader

	inSSBO      uint32
	visibleSSBO uint32
	counterSSBO uint32
	transforms  []raylib.Matrix

	maxVisible uint32
	blades     uint32
	radius     float32
	cellSize   float32
	seed       uint32

	mdl raylib.Model
}

// BladeGPU is std430-friendly if we keep it simple.
// Here: position.xyz + "height" in w.
// (std430: vec4 aligned to 16 bytes -> maps nicely to 4 float32)
type BladeGPU struct {
	Px, Py, Pz, H float32
}

func NewGrass(parent *Terrain, seed uint32) (g *Grass, err error) {
	g = &Grass{}
	// Worst-case visible blades we will store in the "visible" SSBO.
	// If you exceed this, compute will clamp.
	g.maxVisible = uint32(0)
	err = g.Init(parent, seed)
	return g, err
}

func (g *Grass) Init(parent *Terrain, seed uint32) error {
	if g == nil {
		return fmt.Errorf("Invalid grass")
	}
	g.seed = seed
	g.parent = parent

	// Counter buffer (binding=2): uint visibleCount
	var zero uint32 = 0
	g.counterSSBO = raylib.LoadShaderBuffer(
		uint32(unsafe.Sizeof(zero)),
		unsafe.Pointer(&zero),
		int32(raylib.DynamicDraw),
	)

	// Shaders
	g.compute = shader.Compute{}
	err := g.compute.Init("grass")
	if err != nil {
		return err
	}
	g.shader = shader.Shader{}
	err = g.shader.Init("grass")
	if err != nil {
		return err
	}

	// Grass model
	mdl, err := res.GetRes("mdl/grass/grass_2.obj")
	if err != nil {
		return err
	}
	g.mdl = mdl.(raylib.Model)
	g.mdl.Materials.Shader = *g.shader.GetShader()

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
	if g == nil || g.parent == nil || g.inSSBO <= 0 {
		return cmds
	}
	zero := uint32(0)
	raylib.UpdateShaderBuffer(g.counterSSBO, unsafe.Pointer(&zero), 4, 0)

	g.compute.Begin()
	raylib.BindShaderBuffer(g.inSSBO, 0)
	raylib.BindShaderBuffer(g.visibleSSBO, 1)
	raylib.BindShaderBuffer(g.counterSSBO, 2)
	// Upload uniforms to compute
	err := g.compute.SetUniform("uCameraPos", cam.GetPos())
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = g.compute.SetUniform("uMaxVisible", []uint32{g.maxVisible})
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = g.compute.SetUniform("uSeed", []uint32{g.seed})
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = g.compute.SetUniform("uRadius", []float32{g.radius})
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	groups := uint32((g.blades + 256 - 1) / 256)
	raylib.ComputeShaderDispatch(groups, 1, 1)
	g.compute.End()

	var visibleCount uint32
	raylib.ReadShaderBuffer(g.counterSSBO, unsafe.Pointer(&visibleCount), uint32(unsafe.Sizeof(visibleCount)), 0)
	raylib.BindShaderBuffer(g.visibleSSBO, 1)
	err = g.shader.SetUniform("uTime", float32(raylib.GetTime()))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	instances := int(visibleCount)
	raylib.DrawMeshInstanced(*g.mdl.Meshes, *g.mdl.Materials, g.transforms, instances)

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
	if g == nil || g.parent == nil {
		return
	}
	if g.inSSBO > 0 {
		raylib.UnloadShaderBuffer(g.inSSBO)
	}
	if g.visibleSSBO > 0 {
		raylib.UnloadShaderBuffer(g.visibleSSBO)
	}
	var blades []BladeGPU
	hm := g.parent.(*Terrain).GetHeightMap()
	if hm == nil {
		return
	}

	// calculate grass blades
	hmW := hm.Bounds().Max.X - hm.Bounds().Min.X
	hmH := hm.Bounds().Max.Y - hm.Bounds().Min.Y

	Offset := g.parent.GetPos()
	Scale := g.parent.GetScale()
	g.radius = float32(Scale.X / 2)
	g.maxVisible = uint32(g.radius * float32(GrassLOD))

	length := int32(math.Sqrt(float64(g.maxVisible)))
	cellSize := float32(g.radius) / float32(length) * 2

	for gz := int32(0); gz < length; gz++ {
		for gx := int32(0); gx < length; gx++ {
			cellX := gx
			cellZ := gz
			h := rng.Hash2(cellX, cellZ, g.seed)

			// jitter in the cell
			jx := rng.U01(rng.HashU32(h ^ 0xA511E9B3))
			jz := rng.U01(rng.HashU32(h ^ 0x63D83595))

			x := float32(cellX)*cellSize + jx*cellSize
			z := float32(cellZ)*cellSize + jz*cellSize

			// heightmap lookup
			u := int(float32(hmW-1) * (x) / (2 * g.radius))
			v := int(float32(hmH-1) * (z) / (2 * g.radius))
			if u < 0 || u >= hmW || v < 0 || v >= hmH {
				continue
			}
			y := float32(hm.RGBAAt(u, v).R) / 255.0
			if float64(y) < rng.SeaLevel+rng.SandBand {
				continue
			}

			H := 0.35 + (1.10-0.35)*rng.U01(rng.HashU32(h^0xC2B2AE35))

			blades = append(blades, BladeGPU{
				Px: x + Offset.X,
				Py: y*Scale.Y + Offset.Y,
				Pz: z + Offset.Z,
				H:  H,
			})
		}
	}
	g.blades = uint32(len(blades))
	if g.blades < 1 {
		return
	}

	// input blades buffer (binding=0 in compute)
	g.inSSBO = raylib.LoadShaderBuffer(uint32(len(blades))*uint32(unsafe.Sizeof(BladeGPU{})),
		unsafe.Pointer(&blades[0]),
		int32(raylib.StaticDraw),
	)

	// Output visible blades buffer (binding=1 in compute; also binding=1 in vertex shader)
	emptyVisible := make([]BladeGPU, g.maxVisible)
	g.visibleSSBO = raylib.LoadShaderBuffer(
		uint32(len(emptyVisible))*uint32(unsafe.Sizeof(emptyVisible[0])),
		unsafe.Pointer(&emptyVisible[0]),
		int32(raylib.DynamicDraw),
	)

	g.transforms = make([]raylib.Matrix, g.blades)
	for i := range g.transforms {
		g.transforms[i] = raylib.MatrixIdentity()
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

	return g.mdl.Materials
}

func (g *Grass) SetTexture(tex raylib.Texture2D) {
	if g == nil {
		return
	}
	raylib.SetMaterialTexture(g.mdl.Materials, raylib.MapDiffuse, tex)
}

func (g *Grass) GetTexture() *raylib.Texture2D {
	if g == nil {
		return nil
	}

	return &g.mdl.Materials.Maps.Texture
}
