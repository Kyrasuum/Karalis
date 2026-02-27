package world

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

type Dungeon struct {
	tex *raylib.Texture2D
	mdl *raylib.Model
	shd pub_shader.Shader

	hm   *image.RGBA
	seed int64

	pos   raylib.Vector3
	rot   raylib.Vector3
	scale raylib.Vector3
}

func RandDungeon(offx, offy, width, height int, seed int64) (d *Dungeon, err error) {
	d = &Dungeon{}
	if seed == int64(0) {
		d.seed = rand.Int63()
	} else {
		d.seed = seed
	}
	err = d.Init()
	if err != nil {
		return nil, err
	}
	err = d.RandMap(offx, offy, width, height)

	return d, err
}

func NewDungeon(m string, i interface{}, seed int64) (d *Dungeon, err error) {
	d = &Dungeon{}
	err = d.Init()
	if err != nil {
		return nil, err
	}
	if seed == int64(0) {
		d.seed = rand.Int63()
	} else {
		d.seed = seed
	}
	err = d.LoadImage(i)
	if err != nil {
		return nil, err
	}
	err = d.LoadMap(m)
	return d, err
}

func (d *Dungeon) Init() error {
	if d == nil {
		return fmt.Errorf("Invalid Dungeon")
	}

	d.pos = raylib.NewVector3(0, 0, 0)
	d.rot = raylib.NewVector3(0, 0, 0)
	d.scale = raylib.NewVector3(1, 1, 1)

	err := d.LoadImage(nil)
	if err != nil {
		return err
	}
	err = d.LoadMap("")
	if err != nil {
		return err
	}

	d.shd = &shader.Shader{}
	err = d.shd.Init("shader")
	if err != nil {
		return err
	}

	return nil
}

func (d *Dungeon) LoadImage(i interface{}) error {
	if d == nil {
		return fmt.Errorf("Invalid Dungeon")
	}

	var img *raylib.Image
	switch data := i.(type) {
	case string:
		tex, err := res.GetRes(data)
		if err != nil {
			d.LoadImage(nil)
			return err
		}
		var pic image.Image

		pos := strings.Index(data, ".") + 1
		ext := data[pos:]
		switch ext {
		case "png":
			pic, err = png.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				d.LoadImage(nil)
				return err
			}
		case "jpeg":
			pic, err = jpeg.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				d.LoadImage(nil)
				return err
			}
		default:
			pic, _, err = image.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				d.LoadImage(nil)
				return err
			}
		}
		if err != nil {
			return err
		}
		return d.LoadImage(pic)
	case image.Image:
		img = raylib.NewImageFromImage(data)
	case raylib.Color:
		img = raylib.GenImageColor(1536, 256, data)
	default:
		if d.tex != nil {
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
	d.tex = &tex
	if d.mdl != nil {
		raylib.SetMaterialTexture(d.mdl.Materials, raylib.MapDiffuse, *d.tex)
	}
	return nil
}

func (d *Dungeon) RandMap(offx, offy, width, height int) error {
	if d == nil {
		return fmt.Errorf("Invalid Dungeon")
	}

	h := rng.GenerateDungeonTile(width, height, d.seed, offx, offy)
	hm := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			hm.Set(x, y, color.RGBA{h[y*width+x], h[y*width+x], h[y*width+x], uint8(255)})
		}
	}
	d.LoadImage(hm)
	d.GenDungeon(hm)

	alb := rng.ToGrayscaleColors(h, width, height)
	col := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			col.Set(x, y, alb[y*width+x])
		}
	}
	d.LoadImage(col)
	return nil
}

func (d *Dungeon) LoadMap(m string) error {
	if d == nil {
		return fmt.Errorf("Invalid Dungeon")
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
	d.GenDungeon(cube)
	return nil
}

func (d *Dungeon) GenDungeon(img *image.RGBA) {
	d.hm = img
	hm := raylib.NewImageFromImage(img)
	mesh := raylib.GenMeshHeightmap(*hm, raylib.NewVector3(1, 1, 1))
	mdl := raylib.LoadModelFromMesh(mesh)

	d.mdl = &mdl
	if d.tex == nil {
		d.LoadImage(nil)
	}

	raylib.SetMaterialTexture(d.mdl.Materials, raylib.MapDiffuse, *d.tex)
}

func (d *Dungeon) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if d == nil {
		return cmds
	}

	return cmds
}

func (d *Dungeon) Render(cam *camera.Cam) []func() {
	cmds := []func(){}
	if d == nil {
		return cmds
	}

	raylib.Color4ub(255, 255, 255, 255)
	matTransform := d.GetModelMatrix()

	d.mdl.Materials.Shader = *d.shd.GetShader()
	raylib.DrawMesh(*d.mdl.Meshes, *d.mdl.Materials, matTransform)

	return cmds
}

func (d *Dungeon) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if d == nil {
		return cmds
	}

	return cmds
}

func (d *Dungeon) Update(dt float32) {
	if d == nil {
		return
	}
}

func (d *Dungeon) GetCollider() pub_object.Collider {
	if d == nil {
		return nil
	}

	return nil
}

func (d *Dungeon) OnAdd() {
	if d == nil {
		return
	}
}

func (d *Dungeon) OnRemove() {
	if d == nil {
		return
	}
}

func (d *Dungeon) AddChild(obj pub_object.Object) {
	if d == nil {
		return
	}
}

func (d *Dungeon) RemChild(obj pub_object.Object) {
	if d == nil {
		return
	}
}

func (d *Dungeon) GetChilds() []pub_object.Object {
	if d == nil {
		return []pub_object.Object{}
	}

	return []pub_object.Object{}
}

func (d *Dungeon) GetModelMatrix() raylib.Matrix {
	if d == nil {
		return raylib.Matrix{}
	}

	matScale := raylib.MatrixScale(d.scale.X, d.scale.Y, d.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(d.GetPitch()), float64(d.GetYaw()), float64(d.GetRoll()))
	matRotation := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := raylib.MatrixTranslate(d.pos.X, d.pos.Y, d.pos.Z)
	matTransform := raylib.MatrixMultiply(raylib.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = raylib.MatrixMultiply(d.mdl.Transform, matTransform)
	return matTransform
}

func (d *Dungeon) GetModel() *raylib.Model {
	if d == nil {
		return nil
	}

	return d.mdl
}

func (d *Dungeon) SetColor(col color.Color) {
	if d == nil {
		return
	}
}

func (d *Dungeon) GetColor() color.Color {
	if d == nil {
		return nil
	}

	return raylib.White
}

func (d *Dungeon) GetScale() raylib.Vector3 {
	if d == nil {
		return raylib.Vector3{}
	}

	return d.scale
}

func (d *Dungeon) SetScale(sc raylib.Vector3) {
	if d == nil {
		return
	}

	d.scale = sc
}

func (d *Dungeon) SetPos(pos raylib.Vector3) {
	if d == nil {
		return
	}

	d.pos = pos
}

func (d *Dungeon) GetPos() raylib.Vector3 {
	if d == nil {
		return raylib.Vector3{}
	}

	return d.pos
}

func (d *Dungeon) GetPitch() float32 {
	if d == nil {
		return 0
	}

	return d.rot.X
}

func (d *Dungeon) SetPitch(pitch float32) {
	if d == nil {
		return
	}

	d.rot.X = pitch
}

func (d *Dungeon) GetYaw() float32 {
	if d == nil {
		return 0
	}

	return d.rot.Y
}

func (d *Dungeon) SetYaw(yaw float32) {
	if d == nil {
		return
	}

	d.rot.Y = yaw
}

func (d *Dungeon) GetRoll() float32 {
	if d == nil {
		return 0
	}

	return d.rot.Z
}

func (d *Dungeon) SetRoll(roll float32) {
	if d == nil {
		return
	}

	d.rot.Z = roll
}

func (d *Dungeon) GetVertices() []raylib.Vector3 {
	if d == nil {
		return []raylib.Vector3{}
	}

	verts := []raylib.Vector3{}
	length := d.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(d.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (d *Dungeon) GetUVs() []raylib.Vector2 {
	if d == nil {
		return []raylib.Vector2{}
	}

	uvs := []raylib.Vector2{}
	length := d.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(d.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (d *Dungeon) SetUVs(uvs []raylib.Vector2) {
	if d == nil {
		return
	}

	length := int(d.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(d.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	pub_object.UpdateModelUVs(d.mdl)
}

func (d *Dungeon) GetMaterials() *raylib.Material {
	if d == nil {
		return nil
	}

	return d.mdl.Materials
}

func (d *Dungeon) SetTexture(tex raylib.Texture2D) {
	if d == nil {
		return
	}

	*d.tex = tex
}

func (d *Dungeon) GetTexture() *raylib.Texture2D {
	if d == nil {
		return nil
	}

	return d.tex
}

func (d *Dungeon) GetHeightMap() *image.RGBA {
	if d == nil {
		return nil
	}

	return d.hm
}
