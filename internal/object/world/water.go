package world

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"unsafe"

	"karalis/internal/camera"
	"karalis/internal/shader"
	pub_object "karalis/pkg/object"
	"karalis/pkg/rng"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var (
	WaterDetail = float32(2.0)
)

type Water struct {
	parent pub_object.Object
	depth  float32
	shader shader.Shader
	volume shader.Shader
	underw shader.Shader

	resize bool

	waterColor raylib.Vector4
	waveAmp    float32
	waveFreq   float32
	waveSpeed  float32
	fresPow    float32
	specPow    float32
	detailStr  float32

	rt  raylib.RenderTexture2D
	mdl raylib.Model
}

func NewWater(p pub_object.Object, d float32) (water *Water, err error) {
	water = &Water{
		parent: p,
		depth:  d,
		resize: false,
	}
	err = water.Init()

	return water, err
}

func (w *Water) Init() error {
	if w == nil {
		return fmt.Errorf("Invalid Water")
	}

	siz := w.GetScale()
	w.waterColor = raylib.NewVector4(0.05, 0.18, 0.25, 0.75)
	w.waveAmp = float32(0.1) // world units
	w.waveFreq = float32(1 / (siz.X + 0.1))
	w.waveSpeed = float32(0.45) // time scale
	w.fresPow = float32(4.0)
	w.specPow = float32(96.0)
	w.detailStr = float32(0.1)

	w.rt = raylib.LoadRenderTexture(800, 512)

	w.shader = shader.Shader{}
	err := w.shader.Init("water")
	if err != nil {
		return err
	}
	w.volume = *(w.shader.Extend("watervolume").(*shader.Shader))

	w.underw = shader.Shader{}
	err = w.underw.Init("underwater")
	if err != nil {
		return err
	}

	w.Update(0.0)
	return nil
}

func (w *Water) Prerender(cam *camera.Cam) []func() {
	if w == nil {
		return []func(){}
	}

	return []func(){}
}

func (w *Water) Render(cam *camera.Cam) []func() {
	if w == nil {
		return []func(){}
	}
	if w.resize {
		raylib.UnloadRenderTexture(w.rt)
		w.rt = raylib.LoadRenderTexture(int32(raylib.GetRenderWidth()), int32(raylib.GetRenderHeight()))
		w.resize = false
	}

	return []func(){
		func() {
			//enables cutting into water
			raylib.DisableDepthMask()

			//update water shader uniforms
			err := w.shader.SetUniform("uCameraPos", cam.GetPos())
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uTime", float32(raylib.GetTime()))
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uWaterColor", w.waterColor)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uWaveSpeed", w.waveSpeed)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uFresnelPower", w.fresPow)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uSpecPower", w.specPow)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uDetailStrength", w.detailStr)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uWaterHeight", w.depth+float32(rng.SandBand))
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uWaveAmp", w.waveAmp)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.shader.SetUniform("uWaveFreq", w.waveFreq)
			if err != nil {
				log.Printf("%+v\n", err)
			}

			//render water
			siz := w.GetScale()
			pos := w.GetPos()
			scl := raylib.MatrixScale(siz.X, siz.Y, siz.Z)
			dwn := raylib.MatrixTranslate(pos.X, pos.Y+w.depth*siz.Y, pos.Z)
			rot := raylib.MatrixRotate(raylib.Vector3{1, 0, 0}, 0)
			mat := raylib.MatrixMultiply(raylib.MatrixMultiply(rot, scl), dwn)
			w.mdl.Materials.Shader = *w.shader.GetShader()
			raylib.DrawMesh(*w.mdl.Meshes, *w.mdl.Materials, mat)

			//enables cutting into water
			raylib.EnableDepthMask()
			w.shader.End()
		},
		func() {
			//render to texture so we can sample for underwater effect
			raylib.BeginTextureMode(w.rt)
			raylib.ClearBackground(raylib.Black)
			cmds := cam.Render()

			//update water shader uniforms
			err := w.volume.SetUniform("uTime", float32(raylib.GetTime()))
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.volume.SetUniform("uWaveSpeed", w.waveSpeed)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.volume.SetUniform("uWaterHeight", w.depth+float32(rng.SandBand))
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.volume.SetUniform("uWaveAmp", w.waveAmp)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.volume.SetUniform("uWaveFreq", w.waveFreq)
			if err != nil {
				log.Printf("%+v\n", err)
			}

			//render water volume
			siz := w.GetScale()
			pos := w.GetPos()
			scl := raylib.MatrixScale(siz.X, siz.Y, siz.Z)
			dwn := raylib.MatrixTranslate(pos.X, pos.Y+w.depth*siz.Y, pos.Z)
			rot := raylib.MatrixRotate(raylib.Vector3{1, 0, 0}, 0)
			mat := raylib.MatrixMultiply(raylib.MatrixMultiply(rot, scl), dwn)
			w.mdl.Materials.Shader = *w.volume.GetShader()
			raylib.DrawMesh(*w.mdl.Meshes, *w.mdl.Materials, mat)
			w.volume.End()
			for _, cmd := range cmds {
				cmd()
			}
			raylib.EndTextureMode()

			//draw underwater texture
			err = w.underw.SetUniform("uWaterColor", w.waterColor)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = w.underw.SetUniform("uStrength", 1.0)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			w.underw.Begin()
			raylib.DrawTexturePro(w.rt.Texture, raylib.Rectangle{0, 0, float32(w.rt.Texture.Width), -float32(w.rt.Texture.Height)}, raylib.Rectangle{0, 0, float32(w.rt.Texture.Width), float32(w.rt.Texture.Height)}, raylib.Vector2{0, 0}, 0.0, raylib.RayWhite)
			w.underw.End()
		},
	}
}

func (w *Water) Postrender(cam *camera.Cam) []func() {
	if w == nil {
		return []func(){}
	}

	return []func(){}
}

func (w *Water) Update(dt float32) {
	if w == nil {
		return
	}

	size := w.GetScale()
	w.waveFreq = float32(1 / (size.X + 0.1))
	w.waveAmp = float32(rng.SandBand) / 2 * size.Y
	raylib.UnloadModel(w.mdl)

	top := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, 0.0, 0.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	bottom := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, -w.depth, 0.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	north := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, 0.0, 0.0}, raylib.Vector3{0.0, -w.depth, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	south := raylib.GenMeshPlaneExData(raylib.Vector3{1.0, -w.depth, 0.0}, raylib.Vector3{0.0, w.depth, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	east := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, -w.depth, 1.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, w.depth, 0.0}, int(WaterDetail), int(WaterDetail))
	west := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, 0.0, 0.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, -w.depth, 0.0}, int(WaterDetail), int(WaterDetail))

	mesh := raylib.MergeMeshes(raylib.MergeMeshes(
		raylib.MergeMeshes(top, bottom),
		raylib.MergeMeshes(north, south)),
		raylib.MergeMeshes(east, west))
	raylib.UnloadMesh(&top)
	raylib.UnloadMesh(&bottom)
	raylib.UnloadMesh(&north)
	raylib.UnloadMesh(&south)
	raylib.UnloadMesh(&east)
	raylib.UnloadMesh(&west)
	raylib.UploadMesh(&mesh, false)
	w.mdl = raylib.LoadModelFromMesh(mesh)

	hm := w.parent.(*Terrain).GetHeightMap()
	img := raylib.NewImageFromImage(hm)
	tex := raylib.LoadTextureFromImage(img)
	raylib.SetMaterialTexture(w.mdl.Materials, raylib.MapDiffuse, tex)
}

func (w *Water) OnResize(width int32, height int32) {
	if w == nil {
		return
	}
	w.resize = true
}

func (w *Water) GetModelMatrix() raylib.Matrix {
	if w == nil {
		return raylib.Matrix{}
	}

	return w.parent.GetModelMatrix()
}

func (w *Water) GetModel() *raylib.Model {
	if w == nil {
		return nil
	}

	return &w.mdl
}

func (w *Water) SetColor(col color.Color) {
	if w == nil {
		return
	}
	w.parent.SetColor(col)
}

func (w *Water) GetColor() color.Color {
	if w == nil {
		return nil
	}

	return w.parent.GetColor()
}

func (w *Water) SetScale(sc raylib.Vector3) {
	if w == nil {
		return
	}
	w.parent.SetScale(sc)
}

func (w *Water) GetScale() raylib.Vector3 {
	if w == nil {
		return raylib.Vector3{}
	}

	return w.parent.GetScale()
}

func (w *Water) SetPos(p raylib.Vector3) {
	if w == nil {
		return
	}
	w.parent.SetPos(p)
}

func (w *Water) GetPos() raylib.Vector3 {
	if w == nil {
		return raylib.Vector3{}
	}
	return w.parent.GetPos()
}

func (w *Water) GetPitch() float32 {
	if w == nil {
		return 0
	}
	return w.parent.GetPitch()
}

func (w *Water) SetPitch(p float32) {
	if w == nil {
		return
	}
	w.parent.SetPitch(p)
}

func (w *Water) GetYaw() float32 {
	if w == nil {
		return 0
	}
	return w.parent.GetYaw()
}

func (w *Water) SetYaw(y float32) {
	if w == nil {
		return
	}
	w.parent.SetYaw(y)
}

func (w *Water) GetRoll() float32 {
	if w == nil {
		return 0
	}
	return w.parent.GetRoll()
}

func (w *Water) SetRoll(r float32) {
	if w == nil {
		return
	}
	w.parent.SetRoll(r)
}

func (w *Water) GetVertices() []raylib.Vector3 {
	if w == nil {
		return []raylib.Vector3{}
	}

	verts := []raylib.Vector3{}
	length := w.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(w.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (w *Water) GetUVs() []raylib.Vector2 {
	if w == nil {
		return []raylib.Vector2{}
	}

	uvs := []raylib.Vector2{}
	length := w.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(w.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (w *Water) SetUVs(uvs []raylib.Vector2) {
	if w == nil {
		return
	}

	length := int(w.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(w.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	pub_object.UpdateModelUVs(&w.mdl)
}

func (w *Water) GetMaterials() *raylib.Material {
	if w == nil {
		return nil
	}

	return w.mdl.Materials
}

func (w *Water) SetTexture(tex raylib.Texture2D) {
	raylib.SetMaterialTexture(w.mdl.Materials, raylib.MapDiffuse, tex)
}

func (w *Water) GetTexture() *raylib.Texture2D {
	if w == nil {
		return nil
	}

	return &w.mdl.Materials.Maps.Texture
}

func (w *Water) GetCollider() pub_object.Collider {
	if w == nil {
		return nil
	}

	return nil
}

func (w *Water) OnAdd() {
	if w == nil {
		return
	}
}

func (w *Water) OnRemove() {
	if w == nil {
		return
	}
}

func (w *Water) AddChild(obj pub_object.Object) {
	if w == nil {
		return
	}
}

func (w *Water) RemChild(obj pub_object.Object) {
	if w == nil {
		return
	}
}

func (w *Water) GetChilds() []pub_object.Object {
	if w == nil {
		return []pub_object.Object{}
	}

	return []pub_object.Object{}
}
