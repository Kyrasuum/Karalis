package world

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"karalis/internal/rlx"
	"karalis/pkg/app"
	"karalis/pkg/lmath"
	"karalis/pkg/rng"
	"karalis/res"

	pub_object "karalis/pkg/object"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Terrain struct {
	tex *rl.Texture2D
	mdl *rl.Model

	hm   *image.RGBA
	grs  *Grass
	wtr  *Water
	seed int64

	parent  pub_object.Object
	cleaner *runtime.Cleanup
	pos     rl.Vector3
	rot     rl.Vector3
	scale   rl.Vector3
}

func RandTerrain(offx, offy float64, width, height int, seed int64) (t *Terrain, err error) {
	t = &Terrain{}
	if seed == int64(0) {
		t.seed = rand.Int63()
	} else {
		t.seed = seed
	}
	err = t.Init()
	if err != nil {
		return nil, err
	}
	err = t.RandMap(offx, offy, width, height)

	return t, err
}

func NewTerrain(m string, i interface{}, seed int64) (t *Terrain, err error) {
	t = &Terrain{}
	err = t.Init()
	if err != nil {
		return nil, err
	}
	if seed == int64(0) {
		t.seed = rand.Int63()
	} else {
		t.seed = seed
	}
	err = t.LoadImage(i)
	if err != nil {
		return nil, fmt.Errorf("Error loading image: %+v", err)
	}
	err = t.LoadMap(m)
	return t, err
}

func (t *Terrain) Init() error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}
	t.parent = nil
	t.pos = rl.NewVector3(0, 0, 0)
	t.rot = rl.NewVector3(0, 0, 0)
	t.scale = rl.NewVector3(1, 1, 1)

	err := t.LoadImage(nil)
	if err != nil {
		return fmt.Errorf("Error loading image: %+v", err)
	}
	err = t.LoadMap("")
	if err != nil {
		return fmt.Errorf("Error loading map: %+v", err)
	}

	t.grs, err = NewGrass(t, uint32(t.seed))
	if err != nil {
		return fmt.Errorf("Error creating grass: %+v", err)
	}

	t.wtr, err = NewWater(t, float32(rng.SeaLevel+rng.SandBand/2))
	if err != nil {
		return fmt.Errorf("Error creating water: %+v", err)
	}

	return nil
}

func (t *Terrain) LoadImage(i interface{}) error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}

	var img *rl.Image
	switch data := i.(type) {
	case string:
		tex, err := res.GetRes(data)
		if err != nil {
			t.LoadImage(nil)
			return err
		}
		var pic image.Image

		pos := strings.Index(data, ".") + 1
		ext := data[pos:]
		switch ext {
		case "png":
			pic, err = png.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				t.LoadImage(nil)
				return err
			}
		case "jpeg":
			pic, err = jpeg.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				t.LoadImage(nil)
				return err
			}
		default:
			pic, _, err = image.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				t.LoadImage(nil)
				return err
			}
		}
		if err != nil {
			return err
		}
		return t.LoadImage(pic)
	case image.Image:
		img = rlx.NewImageFromImage(data)
	case rl.Color:
		img = rlx.GenImageColor(1536, 256, data)
	default:
		if t.tex != nil {
			return nil
		}
		width := 256
		height := 256
		color := color.RGBA{uint8(255), uint8(255), uint8(255), uint8(255)}
		cube := image.NewRGBA(image.Rect(0, 0, width, height))
		for i := range width {
			for j := range height {
				cube.Set(i, j, color)
			}
		}
		img = rlx.NewImageFromImage(cube)
	}
	tex := rlx.LoadTextureFromImage(img)
	t.tex = &tex
	if t.cleaner != nil {
		t.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(t, func(in []interface{}) {
		rlx.UnloadTexture(in[0].(rl.Texture2D))
	}, []interface{}{*t.tex})
	t.cleaner = &cleaner

	if t.mdl != nil {
		rlx.SetMaterialTexture(t.mdl.Materials, rl.MapDiffuse, *t.tex)
		t.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
		if t.cleaner != nil {
			t.cleaner.Stop()
		}
		cleaner := runtime.AddCleanup(t, func(in []interface{}) {
			rlx.UnloadTexture(in[0].(rl.Texture2D))
			rlx.UnloadModel(in[1].(rl.Model))
		}, []interface{}{*t.tex, *t.mdl})
		t.cleaner = &cleaner
	}
	t.wtr.Update(0.0)
	return nil
}

func (t *Terrain) RandMap(offx, offy float64, width, height int) error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}

	h := rng.GenerateHeightmapTiledWorldSize(width, height, t.seed, offx, offy, 10)
	hm := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			hm.Set(x, y, h[y*width+x])
		}
	}
	t.LoadImage(hm)
	t.GenTerrain(hm)

	alb := rng.ColorizeHeightmapTiled(h, width, height, rand.Int63(), 0, 0, 10, width*8, height*8, false)
	col := image.NewRGBA(image.Rect(0, 0, width*8, height*8))
	for y := range height * 8 {
		for x := range width * 8 {
			col.Set(x, y, alb[y*width*8+x])
		}
	}
	t.LoadImage(col)
	return nil
}

func (t *Terrain) LoadMap(m string) error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}

	tex, err := res.GetRes(m)
	if err != nil && strings.Compare(m, "") != 0 {
		log.Printf("Error retrieving image (%s): %+v\n", m, err)
		m = ""
	}
	var goimg image.Image

	if tex != nil {
		pos := strings.Index(m, ".") + 1
		ext := m[pos:]
		switch ext {
		case "png":
			goimg, err = png.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				return fmt.Errorf("Error decoding image: %+v\n", err)
			}
		case "jpeg":
			goimg, err = jpeg.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				return fmt.Errorf("Error decoding image: %+v\n", err)

			}
		default:
			goimg, _, err = image.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				return fmt.Errorf("Error decoding image: %+v\n", err)

			}
		}
	} else {
		goimg = image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{128, 128}})
	}
	width := goimg.Bounds().Dx()
	height := goimg.Bounds().Dy()
	cube := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := range width {
		for j := range height {
			cube.Set(i, j, goimg.At(i, j))
		}
	}
	t.GenTerrain(cube)
	return nil
}

func (t *Terrain) GenTerrain(img *image.RGBA) {
	t.hm = img
	hm := rlx.NewImageFromImage(img)
	mesh := rlx.GenMeshHeightmap(*hm, rl.NewVector3(1, 1, 1))
	mdl := rlx.LoadModelFromMesh(mesh)

	t.mdl = &mdl
	if t.tex == nil {
		t.LoadImage(nil)
	}

	if t.grs != nil {
		t.grs.Update(0.0)
	}

	rlx.SetMaterialTexture(t.mdl.Materials, rl.MapDiffuse, *t.tex)
	t.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
	if t.cleaner != nil {
		t.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(t, func(in []interface{}) {
		rlx.UnloadTexture(in[0].(rl.Texture2D))
		rlx.UnloadModel(in[1].(rl.Model))
	}, []interface{}{*t.tex, *t.mdl})
	t.cleaner = &cleaner
}

func (t *Terrain) Prerender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if t == nil {
		return cmds
	}
	cmds = append(cmds, t.grs.Prerender(cam)...)
	cmds = append(cmds, t.wtr.Prerender(cam)...)

	return cmds
}

func (t *Terrain) Render(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if t == nil {
		return cmds
	}

	rlx.Color4ub(255, 255, 255, 255)
	matTransform := t.GetModelMatrix()

	rlx.DrawMesh(*t.mdl.Meshes, *t.mdl.Materials, matTransform)
	cmds = append(cmds, t.grs.Render(cam)...)
	cmds = append(cmds, t.wtr.Render(cam)...)

	return cmds
}

func (t *Terrain) Postrender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if t == nil {
		return cmds
	}
	cmds = append(cmds, t.grs.Postrender(cam)...)
	cmds = append(cmds, t.wtr.Postrender(cam)...)

	return cmds
}

func (t *Terrain) OnResize(w int32, h int32) {
	if t == nil {
		return
	}
	t.grs.OnResize(w, h)
	t.wtr.OnResize(w, h)
}

func (t *Terrain) Update(dt float32) {
	if t == nil {
		return
	}
}

func (t *Terrain) GetCollider() pub_object.Collider {
	if t == nil {
		return nil
	}

	return nil
}

func (t *Terrain) OnAdd(obj pub_object.Object) {
	if t == nil {
		return
	}
	t.parent = obj
}

func (t *Terrain) OnRemove() {
	if t == nil {
		return
	}
	t.parent = nil
}

func (t *Terrain) AddChild(obj pub_object.Object) {
	if t == nil {
		return
	}
}

func (t *Terrain) RemChild(obj pub_object.Object) {
	if t == nil {
		return
	}
}

func (t *Terrain) GetChilds() []pub_object.Object {
	if t == nil {
		return []pub_object.Object{}
	}

	return []pub_object.Object{}
}

func (t *Terrain) GetModelMatrix() rl.Matrix {
	if t == nil {
		return rl.Matrix{}
	}

	matScale := rl.MatrixScale(t.scale.X, t.scale.Y, t.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(t.GetPitch()), float64(t.GetYaw()), float64(t.GetRoll()))
	matRotation := rl.QuaternionToMatrix(rl.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := rl.MatrixTranslate(t.pos.X, t.pos.Y, t.pos.Z)
	matTransform := rl.MatrixMultiply(rl.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = rl.MatrixMultiply(t.mdl.Transform, matTransform)
	return matTransform
}

func (t *Terrain) GetModel() *rl.Model {
	if t == nil {
		return nil
	}

	return t.mdl
}

func (t *Terrain) SetColor(col color.Color) {
	if t == nil {
		return
	}
}

func (t *Terrain) GetColor() color.Color {
	if t == nil {
		return nil
	}

	return rl.White
}

func (t *Terrain) GetScale() rl.Vector3 {
	if t == nil {
		return rl.Vector3{}
	}

	return t.scale
}

func (t *Terrain) SetScale(sc rl.Vector3) {
	if t == nil {
		return
	}

	t.scale = sc
	if t.grs != nil {
		t.grs.Update(0.0)
	}
	if t.wtr != nil {
		t.wtr.Update(0.0)
	}
}

func (t *Terrain) SetPos(pos rl.Vector3) {
	if t == nil {
		return
	}

	t.pos = pos
	if t.grs != nil {
		t.grs.Update(0.0)
	}
}

func (t *Terrain) GetPos() rl.Vector3 {
	if t == nil {
		return rl.Vector3{}
	}

	return t.pos
}

func (t *Terrain) GetPitch() float32 {
	if t == nil {
		return 0
	}

	return t.rot.X
}

func (t *Terrain) SetPitch(pitch float32) {
	if t == nil {
		return
	}

	t.rot.X = pitch
	if t.grs != nil {
		t.grs.Update(0.0)
	}
}

func (t *Terrain) GetYaw() float32 {
	if t == nil {
		return 0
	}

	return t.rot.Y
}

func (t *Terrain) SetYaw(yaw float32) {
	if t == nil {
		return
	}

	t.rot.Y = yaw
	if t.grs != nil {
		t.grs.Update(0.0)
	}
}

func (t *Terrain) GetRoll() float32 {
	if t == nil {
		return 0
	}

	return t.rot.Z
}

func (t *Terrain) SetRoll(roll float32) {
	if t == nil {
		return
	}

	t.rot.Z = roll
	if t.grs != nil {
		t.grs.Update(0.0)
	}
}

func (t *Terrain) GetVertices() []rl.Vector3 {
	if t == nil {
		return []rl.Vector3{}
	}

	verts := []rl.Vector3{}
	length := t.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(t.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, rl.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (t *Terrain) GetUVs() []rl.Vector2 {
	if t == nil {
		return []rl.Vector2{}
	}

	uvs := []rl.Vector2{}
	length := t.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(t.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, rl.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (t *Terrain) SetUVs(uvs []rl.Vector2) {
	if t == nil {
		return
	}

	length := int(t.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(t.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	pub_object.UpdateModelUVs(t.mdl)
}

func (t *Terrain) GetMaterials() *rl.Material {
	if t == nil {
		return nil
	}

	return t.mdl.Materials
}

func (t *Terrain) SetTexture(tex rl.Texture2D) {
	if t == nil {
		return
	}

	*t.tex = tex
	if t.cleaner != nil {
		t.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(t, func(in []interface{}) {
		rlx.UnloadTexture(in[0].(rl.Texture2D))
	}, []interface{}{*t.tex})
	t.cleaner = &cleaner

	if t.mdl != nil {
		rlx.SetMaterialTexture(t.mdl.Materials, rl.MapDiffuse, *t.tex)
		t.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
		if t.cleaner != nil {
			t.cleaner.Stop()
		}
		cleaner := runtime.AddCleanup(t, func(in []interface{}) {
			rlx.UnloadTexture(in[0].(rl.Texture2D))
			rlx.UnloadModel(in[1].(rl.Model))
		}, []interface{}{*t.tex, *t.mdl})
		t.cleaner = &cleaner
	}
}

func (t *Terrain) GetTexture() *rl.Texture2D {
	if t == nil {
		return nil
	}

	return t.tex
}

func (t *Terrain) GetHeightMap() *image.RGBA {
	if t == nil {
		return nil
	}

	return t.hm
}

func (t *Terrain) GetParent() pub_object.Object {
	if t == nil {
		return nil
	}
	return t.parent
}
