package world

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"runtime"
	"unsafe"

	"karalis/internal/shader"
	"karalis/pkg/app"
	"karalis/pkg/lmath"
	pub_object "karalis/pkg/object"
	"karalis/pkg/rng"
	pub_shader "karalis/pkg/shader"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var (
	WaterDetail  = float32(2.0)
	UnderWaterRT = raylib.LoadRenderTexture(int32(raylib.GetRenderWidth()), int32(raylib.GetRenderHeight()))
	UnderWaterSh = &shader.Shader{}
)

type Water struct {
	parent  pub_object.Object
	cleaner *runtime.Cleanup

	depth  float32
	shader pub_shader.Shader
	volume pub_shader.Shader

	resize bool

	waterColor raylib.Vector4
	waveAmp    float32
	waveFreq   float32
	waveSpeed  float32
	fresPow    float32
	specPow    float32
	detailStr  float32

	bot   *raylib.Model
	sfc   *raylib.Model
	siden *raylib.Model
	sides *raylib.Model
	sidee *raylib.Model
	sidew *raylib.Model
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

	var err error
	w.shader, err = shader.NewShader("water")
	if err != nil {
		return err
	}
	w.volume = w.shader.Extend("watervolume")

	w.Update(0.0)
	return nil
}

func (w *Water) DrawEdges(mat raylib.Matrix, shader pub_shader.Shader) []func() {
	if w == nil {
		return []func(){}
	}

	ply := app.CurApp.GetStage().GetPlayer()
	plypos := raylib.Vector2{float32(lmath.Round(ply.GetPos().X / CellScale.X)), float32(lmath.Round(ply.GetPos().Z / CellScale.Z))}

	cpos := raylib.Vector2{float32(lmath.Round((w.GetPos().X + CellScale.X/2) / CellScale.X)), float32(lmath.Round((w.GetPos().Z + CellScale.Z/2) / CellScale.Z))}
	diff := raylib.Vector2{plypos.X - cpos.X, plypos.Y - cpos.Y}

	if diff.X == float32(CellRender) {
		w.siden.Materials.Shader = *shader.GetShader()
		raylib.DrawMesh(*w.siden.Meshes, *w.siden.Materials, mat)
	}
	if diff.X == -float32(CellRender) {
		w.sides.Materials.Shader = *shader.GetShader()
		raylib.DrawMesh(*w.sides.Meshes, *w.sides.Materials, mat)
	}
	if diff.Y == float32(CellRender) {
		w.sidew.Materials.Shader = *shader.GetShader()
		raylib.DrawMesh(*w.sidew.Meshes, *w.sidew.Materials, mat)
	}
	if diff.Y == -float32(CellRender) {
		w.sidee.Materials.Shader = *shader.GetShader()
		raylib.DrawMesh(*w.sidee.Meshes, *w.sidee.Materials, mat)
	}
	return []func(){}
}

func (w *Water) DrawWater(cam pub_object.Camera, time float32) []func() {
	if w == nil {
		return []func(){}
	}
	//enables cutting into water
	raylib.DisableDepthMask()

	//update water shader uniforms
	err := w.shader.SetUniform("uCameraPos", cam.GetPos())
	if err != nil {
		log.Printf("%+v\n", err)
	}
	err = w.shader.SetUniform("uTime", time)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	err = w.shader.SetUniform("uTexSize", float32(TerrainDetail))
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
	w.sfc.Materials.Shader = *w.shader.GetShader()
	raylib.DrawMesh(*w.sfc.Meshes, *w.sfc.Materials, mat)
	w.bot.Materials.Shader = *w.shader.GetShader()
	raylib.DrawMesh(*w.bot.Meshes, *w.bot.Materials, mat)
	w.DrawEdges(mat, w.shader)

	//enables cutting into water
	raylib.EnableDepthMask()
	w.shader.End()

	return []func(){}
}

func (w *Water) DrawUnderwater(cam pub_object.Camera, time float32) []func() {
	if w == nil {
		return []func(){}
	}
	//render to texture so we can sample for underwater effect
	raylib.BeginTextureMode(UnderWaterRT)
	cmds := cam.Render()

	//update water shader uniforms
	err := w.volume.SetUniform("uTime", time)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	err = w.volume.SetUniform("uTexSize", float32(TerrainDetail))
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
	w.sfc.Materials.Shader = *w.volume.GetShader()
	raylib.DrawMesh(*w.sfc.Meshes, *w.sfc.Materials, mat)
	w.bot.Materials.Shader = *w.volume.GetShader()
	raylib.DrawMesh(*w.bot.Meshes, *w.bot.Materials, mat)
	w.DrawEdges(mat, w.volume)
	w.volume.End()
	for _, cmd := range cmds {
		cmd()
	}
	raylib.EndTextureMode()

	return []func(){}
}

func DrawUnderwater(world *World) []func() {
	//check if shader is ready
	if UnderWaterSh.GetShader() == nil {
		shd, err := shader.NewShader("underwater")
		if err != nil {
			log.Printf("%+v\n", err)
		}
		UnderWaterSh = shd.(*shader.Shader)
	}

	cmds := []func(){}
	if world != nil {
		cmds = append(cmds, func() {
			//find current water cell
			cpos := raylib.Vector2{float32(lmath.Round(world.GetPos().X / CellScale.X)), float32(lmath.Round(world.GetPos().Z / CellScale.Z))}
			spos := fmt.Sprintf("%d %d", cpos.X, cpos.Y)
			var water *Water
			if cell, ok := world.cells[spos]; ok {
				water = cell.terrain.(*Terrain).wtr
			} else {
				return
			}

			//draw underwater texture
			err := UnderWaterSh.SetUniform("uWaterColor", water.waterColor)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			err = UnderWaterSh.SetUniform("uStrength", 1.0)
			if err != nil {
				log.Printf("%+v\n", err)
			}
			UnderWaterSh.Begin()
			raylib.DrawTexturePro(
				UnderWaterRT.Texture,
				raylib.Rectangle{0, 0, float32(UnderWaterRT.Texture.Width), -float32(UnderWaterRT.Texture.Height)},
				raylib.Rectangle{0, 0, float32(UnderWaterRT.Texture.Width), float32(UnderWaterRT.Texture.Height)},
				raylib.Vector2{0, 0}, 0.0, raylib.RayWhite)
			UnderWaterSh.End()
		})
	}

	//update render texture
	width := int32(raylib.GetRenderWidth())
	height := int32(raylib.GetRenderHeight())
	if UnderWaterRT.Texture.Width != width || UnderWaterRT.Texture.Height != height {
		raylib.UnloadRenderTexture(UnderWaterRT)
		UnderWaterRT = raylib.LoadRenderTexture(width, height)
	}
	return cmds
}

func (w *Water) Prerender(cam pub_object.Camera) []func() {
	if w == nil {
		return []func(){}
	}

	raylib.BeginTextureMode(UnderWaterRT)
	raylib.ClearBackground(raylib.Black)
	raylib.EndTextureMode()

	return []func(){}
}

func (w *Water) Render(cam pub_object.Camera) []func() {
	if w == nil {
		return []func(){}
	}

	time := float32(raylib.GetTime())

	return []func(){
		func() {
			w.DrawWater(cam, time)
		},
	}
}

func (w *Water) Postrender(cam pub_object.Camera) []func() {
	if w == nil {
		return []func(){}
	}
	time := float32(raylib.GetTime())
	return w.DrawUnderwater(cam, time)
	// return []func(){}
}

func (w *Water) Update(dt float32) {
	if w == nil {
		return
	}

	size := w.GetScale()
	w.waveFreq = float32(1 / (size.X + 0.1))
	w.waveAmp = float32(rng.SandBand) / 2 * size.Y
	if w.sfc != nil {
		raylib.UnloadModel(*w.sfc)
	}
	if w.bot != nil {
		raylib.UnloadModel(*w.bot)
	}
	if w.siden != nil {
		raylib.UnloadModel(*w.siden)
	}
	if w.sides != nil {
		raylib.UnloadModel(*w.sides)
	}
	if w.sidee != nil {
		raylib.UnloadModel(*w.sidee)
	}
	if w.sidew != nil {
		raylib.UnloadModel(*w.sidew)
	}

	top := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, 0.0, 0.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	bottom := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, -w.depth, 0.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	north := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, 0.0, 0.0}, raylib.Vector3{0.0, -w.depth, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	south := raylib.GenMeshPlaneExData(raylib.Vector3{1.0, -w.depth, 0.0}, raylib.Vector3{0.0, w.depth, 0.0}, raylib.Vector3{0.0, 0.0, 1.0}, int(WaterDetail), int(WaterDetail))
	east := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, -w.depth, 1.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, w.depth, 0.0}, int(WaterDetail), int(WaterDetail))
	west := raylib.GenMeshPlaneExData(raylib.Vector3{0.0, 0.0, 0.0}, raylib.Vector3{1.0, 0.0, 0.0}, raylib.Vector3{0.0, -w.depth, 0.0}, int(WaterDetail), int(WaterDetail))

	raylib.UploadMesh(&north, false)
	raylib.UploadMesh(&south, false)
	raylib.UploadMesh(&east, false)
	raylib.UploadMesh(&west, false)
	raylib.UploadMesh(&top, false)
	raylib.UploadMesh(&bottom, false)
	bot := raylib.LoadModelFromMesh(bottom)
	sfc := raylib.LoadModelFromMesh(top)
	siden := raylib.LoadModelFromMesh(north)
	sides := raylib.LoadModelFromMesh(south)
	sidee := raylib.LoadModelFromMesh(east)
	sidew := raylib.LoadModelFromMesh(west)
	w.bot = &bot
	w.sfc = &sfc
	w.siden = &siden
	w.sides = &sides
	w.sidee = &sidee
	w.sidew = &sidew

	hm := w.parent.(*Terrain).GetHeightMap()
	img := raylib.NewImageFromImage(hm)
	tex := raylib.LoadTextureFromImage(img)
	raylib.SetMaterialTexture(w.sfc.Materials, raylib.MapDiffuse, tex)
	raylib.SetMaterialTexture(w.bot.Materials, raylib.MapDiffuse, tex)
	raylib.SetMaterialTexture(w.siden.Materials, raylib.MapDiffuse, tex)
	raylib.SetMaterialTexture(w.sides.Materials, raylib.MapDiffuse, tex)
	raylib.SetMaterialTexture(w.sidee.Materials, raylib.MapDiffuse, tex)
	raylib.SetMaterialTexture(w.sidew.Materials, raylib.MapDiffuse, tex)

	if w.cleaner != nil {
		w.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(w, func(mdls []raylib.Model) {
		for _, mdl := range mdls {
			raylib.UnloadModel(mdl)
		}
	}, []raylib.Model{*w.bot, *w.sfc, *w.siden, *w.sides, *w.sidee, *w.sidew})
	w.cleaner = &cleaner
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

	return w.sfc
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
	length := w.sfc.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(w.sfc.Meshes.Vertices))
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
	length := w.sfc.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(w.sfc.Meshes.Texcoords))
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

	length := int(w.sfc.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(w.sfc.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	pub_object.UpdateModelUVs(w.sfc)
}

func (w *Water) GetMaterials() *raylib.Material {
	if w == nil {
		return nil
	}

	return w.sfc.Materials
}

func (w *Water) SetTexture(tex raylib.Texture2D) {
	raylib.SetMaterialTexture(w.sfc.Materials, raylib.MapDiffuse, tex)
}

func (w *Water) GetTexture() *raylib.Texture2D {
	if w == nil {
		return nil
	}

	return &w.sfc.Materials.Maps.Texture
}

func (w *Water) GetCollider() pub_object.Collider {
	if w == nil {
		return nil
	}

	return nil
}

func (w *Water) OnAdd(obj pub_object.Object) {
	if w == nil {
		return
	}
	w.parent = obj
}

func (w *Water) OnRemove() {
	if w == nil {
		return
	}
	w.parent = nil
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

func (w *Water) GetParent() pub_object.Object {
	if w == nil {
		return nil
	}
	return w.parent
}
