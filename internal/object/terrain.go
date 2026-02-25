package object

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"reflect"
	"strings"
	"unsafe"

	"karalis/internal/camera"
	"karalis/internal/shader"
	pub_object "karalis/pkg/object"
	"karalis/pkg/rng"
	pub_shader "karalis/pkg/shader"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

type Terrain struct {
	tex *raylib.Texture2D
	hgt *raylib.Texture2D
	mdl *raylib.Model
	shd pub_shader.Shader

	hm  raylib.Texture2D
	grs pub_shader.Shader
	grd *raylib.Texture2D

	pos   raylib.Vector3
	rot   raylib.Vector3
	scale raylib.Vector3
}

func RandTerrain() (t *Terrain, err error) {
	t = &Terrain{}
	err = t.Init()
	if err != nil {
		return nil, err
	}
	err = t.RandMap()

	return t, err
}

func NewTerrain(m string, i interface{}) (t *Terrain, err error) {
	t = &Terrain{}
	err = t.Init()
	if err != nil {
		return nil, err
	}
	err = t.LoadImage(i)
	if err != nil {
		return nil, err
	}
	err = t.LoadMap(m)
	return t, err
}

func (t *Terrain) Init() error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}

	t.pos = raylib.NewVector3(0, 0, 0)
	t.rot = raylib.NewVector3(0, 0, 0)
	t.scale = raylib.NewVector3(1, 1, 1)

	err := t.LoadImage(nil)
	if err != nil {
		return err
	}
	err = t.LoadMap("")
	if err != nil {
		return err
	}

	t.shd = &shader.Shader{}
	t.shd.Init("shader")
	t.grs = t.shd.Extend("grass")
	t.LoadGrass()

	return nil
}

func (t *Terrain) LoadGrass() error {
	data, err := res.GetRes("tex/grass.png")
	if err != nil {
		return err
	}

	pic, err := png.Decode(bytes.NewReader(data.([]byte)))
	if err != nil {
		return err
	}
	img := raylib.NewImageFromImage(pic)
	tex := raylib.LoadTextureFromImage(img)
	t.grd = &tex
	return nil
}

func (t *Terrain) LoadImage(i interface{}) error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}

	var img *raylib.Image
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
			t.LoadImage(nil)
			return err
		}
		t.LoadImage(pic)
		return nil
	case image.Image:
		img = raylib.NewImageFromImage(data)
	case raylib.Color:
		img = raylib.GenImageColor(1536, 256, data)
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
		img = raylib.NewImageFromImage(cube)
	}
	tex := raylib.LoadTextureFromImage(img)
	t.tex = &tex
	if t.mdl != nil {
		raylib.SetMaterialTexture(t.mdl.Materials, raylib.MapDiffuse, *t.tex)
	}
	return nil
}

func (t *Terrain) RandMap() error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}

	width := 256
	height := 256
	h := rng.GenerateHeightmapTiledWorldSize(width, height, int64(12345), 0, 0, 10)
	hm := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			hm.Set(x, y, h[y*width+x])
		}
	}
	himg := raylib.NewImageFromImage(hm)
	t.LoadImage(hm)
	t.GenTerrain(himg)

	alb := rng.ColorizeHeightmapTiled(h, width, height, rand.Int63(), 0, 0, 10, width*8, height*8)
	col := image.NewRGBA(image.Rect(0, 0, width*8, height*8))
	for y := range height * 8 {
		for x := range width * 8 {
			col.Set(x, y, alb[y*width*8+x])
		}
	}
	t.LoadImage(col)
	t.hm = raylib.LoadTextureFromImage(himg)
	return nil
}

func (t *Terrain) LoadMap(m string) error {
	if t == nil {
		return fmt.Errorf("Invalid terrain")
	}

	tex, err := res.GetRes(m)
	if err != nil && strings.Compare(m, "") != 0 {
		fmt.Printf("Error retrieving image (%s): %+v\n", m, err)
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

	img := raylib.NewImageFromImage(cube)
	t.GenTerrain(img)
	return nil
}

func (t *Terrain) GenTerrain(img *raylib.Image) {
	mesh := raylib.GenMeshHeightmap(*img, raylib.NewVector3(1, 1, 1))
	mdl := raylib.LoadModelFromMesh(mesh)
	t.mdl = &mdl
	if t.tex == nil {
		t.LoadImage(nil)
	}
	hmap := raylib.LoadTextureFromImage(img)
	t.hgt = &hmap
	raylib.SetMaterialTexture(t.mdl.Materials, raylib.MapDiffuse, *t.tex)
}

func (t *Terrain) GetModelMatrix() raylib.Matrix {
	if t == nil {
		return raylib.Matrix{}
	}

	matScale := raylib.MatrixScale(t.scale.X, t.scale.Y, t.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(t.GetPitch()), float64(t.GetYaw()), float64(t.GetRoll()))
	matRotation := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := raylib.MatrixTranslate(t.pos.X, t.pos.Y, t.pos.Z)
	matTransform := raylib.MatrixMultiply(raylib.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = raylib.MatrixMultiply(t.mdl.Transform, matTransform)
	return matTransform
}

func (t *Terrain) GetModel() *raylib.Model {
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

	return raylib.White
}

func (t *Terrain) GetScale() raylib.Vector3 {
	if t == nil {
		return raylib.Vector3{}
	}

	return t.scale
}

func (t *Terrain) SetScale(sc raylib.Vector3) {
	if t == nil {
		return
	}

	t.scale = sc
}

func (t *Terrain) SetPos(pos raylib.Vector3) {
	if t == nil {
		return
	}

	t.pos = pos
}

func (t *Terrain) GetPos() raylib.Vector3 {
	if t == nil {
		return raylib.Vector3{}
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
}

func (t *Terrain) GetVertices() []raylib.Vector3 {
	if t == nil {
		return []raylib.Vector3{}
	}

	verts := []raylib.Vector3{}
	length := t.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(t.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (t *Terrain) GetUVs() []raylib.Vector2 {
	if t == nil {
		return []raylib.Vector2{}
	}

	uvs := []raylib.Vector2{}
	length := t.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(t.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (t *Terrain) SetUVs(uvs []raylib.Vector2) {
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

func (t *Terrain) GetMaterials() *raylib.Material {
	if t == nil {
		return nil
	}

	return t.mdl.Materials
}

func (t *Terrain) SetTexture(tex raylib.Texture2D) {
	if t == nil {
		return
	}

	*t.tex = tex
}

func (t *Terrain) GetTexture() raylib.Texture2D {
	if t == nil {
		return raylib.Texture2D{}
	}

	return *t.tex
}

func (t *Terrain) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if t == nil {
		return cmds
	}

	return cmds
}

func (t *Terrain) Render(cam *camera.Cam) []func() {
	cmds := []func(){}
	if t == nil {
		return cmds
	}

	raylib.Color4ub(255, 255, 255, 255)
	matTransform := t.GetModelMatrix()

	t.mdl.Materials.Shader = *t.shd.GetShader()
	raylib.DrawMesh(*t.mdl.Meshes, *t.mdl.Materials, matTransform)

	raylib.SetMaterialTexture(t.mdl.Materials, raylib.MapDiffuse, *t.grd)
	raylib.SetMaterialTexture(t.mdl.Materials, raylib.MapDiffuse, t.hm)
	err := t.grs.SetUniform("texture1", t.hm)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uHeightmapTexelScale", raylib.Vector2{8, 8})
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uGrassMinHeight", rng.MinGrassHeight())
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uGrassMaxHeight", rng.MaxGrassHeight())
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uGrassDensity", 0.35*255)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uMaxSlope", 0.99)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uBladeHeight", 0.05)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uBladeHalfWidth", 0.01)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = t.grs.SetUniform("uSeed", float64(rand.Int63()))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	t.mdl.Materials.Shader = *t.grs.GetShader()
	raylib.DrawMesh(*t.mdl.Meshes, *t.mdl.Materials, matTransform)
	raylib.SetMaterialTexture(t.mdl.Materials, raylib.MapDiffuse, *t.tex)

	return cmds
}

func (t *Terrain) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if t == nil {
		return cmds
	}

	return cmds
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

func (t *Terrain) OnAdd() {
	if t == nil {
		return
	}
}

func (t *Terrain) OnRemove() {
	if t == nil {
		return
	}
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
