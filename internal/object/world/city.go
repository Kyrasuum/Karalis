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
	"strings"
	"unsafe"

	"karalis/internal/camera"
	"karalis/pkg/app"
	pub_object "karalis/pkg/object"
	"karalis/pkg/rng"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

type City struct {
	tex *raylib.Texture2D
	mdl *raylib.Model

	hm   *image.RGBA
	seed int64

	pos   raylib.Vector3
	rot   raylib.Vector3
	scale raylib.Vector3
}

func RandCity(offx, offy, width, height int, seed int64) (c *City, err error) {
	c = &City{}
	if seed == int64(0) {
		c.seed = rand.Int63()
	} else {
		c.seed = seed
	}
	err = c.Init()
	if err != nil {
		return nil, err
	}
	err = c.RandMap(offx, offy, width, height)

	return c, err
}

func NewCity(m string, i interface{}, seed int64) (c *City, err error) {
	c = &City{}
	err = c.Init()
	if err != nil {
		return nil, err
	}
	if seed == int64(0) {
		c.seed = rand.Int63()
	} else {
		c.seed = seed
	}
	err = c.LoadImage(i)
	if err != nil {
		return nil, err
	}
	err = c.LoadMap(m)
	return c, err
}

func (c *City) Init() error {
	if c == nil {
		return fmt.Errorf("Invalid City")
	}

	c.pos = raylib.NewVector3(0, 0, 0)
	c.rot = raylib.NewVector3(0, 0, 0)
	c.scale = raylib.NewVector3(1, 1, 1)

	err := c.LoadImage(nil)
	if err != nil {
		return err
	}
	err = c.LoadMap("")
	if err != nil {
		return err
	}

	return nil
}

func (c *City) LoadImage(i interface{}) error {
	if c == nil {
		return fmt.Errorf("Invalid City")
	}

	var img *raylib.Image
	switch data := i.(type) {
	case string:
		tex, err := res.GetRes(data)
		if err != nil {
			c.LoadImage(nil)
			return err
		}
		var pic image.Image

		pos := strings.Index(data, ".") + 1
		ext := data[pos:]
		switch ext {
		case "png":
			pic, err = png.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				c.LoadImage(nil)
				return err
			}
		case "jpeg":
			pic, err = jpeg.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				c.LoadImage(nil)
				return err
			}
		default:
			pic, _, err = image.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				c.LoadImage(nil)
				return err
			}
		}
		if err != nil {
			return err
		}
		return c.LoadImage(pic)
	case image.Image:
		img = raylib.NewImageFromImage(data)
	case raylib.Color:
		img = raylib.GenImageColor(1536, 256, data)
	default:
		if c.tex != nil {
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
	c.tex = &tex
	if c.mdl != nil {
		raylib.SetMaterialTexture(c.mdl.Materials, raylib.MapDiffuse, *c.tex)
		c.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
	}
	return nil
}

func (c *City) RandMap(offx, offy, width, height int) error {
	if c == nil {
		return fmt.Errorf("Invalid City")
	}

	h := rng.GenerateCityTile(width, height, c.seed, float64(offx), float64(offy), 10)
	hm := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			hm.Set(x, y, color.RGBA{h[y*width+x], h[y*width+x], h[y*width+x], uint8(255)})
		}
	}
	c.LoadImage(hm)
	c.GenCity(hm)
	return nil
}

func (c *City) LoadMap(m string) error {
	if c == nil {
		return fmt.Errorf("Invalid City")
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
	c.GenCity(cube)
	return nil
}

func (c *City) GenCity(img *image.RGBA) {
	c.hm = img
	hm := raylib.NewImageFromImage(img)
	mesh := raylib.GenMeshHeightmap(*hm, raylib.NewVector3(1, 1, 1))
	mdl := raylib.LoadModelFromMesh(mesh)

	c.mdl = &mdl
	if c.tex == nil {
		c.LoadImage(nil)
	}

	raylib.SetMaterialTexture(c.mdl.Materials, raylib.MapDiffuse, *c.tex)
	c.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
}

func (c *City) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if c == nil {
		return cmds
	}

	return cmds
}

func (c *City) Render(cam *camera.Cam) []func() {
	cmds := []func(){}
	if c == nil {
		return cmds
	}

	raylib.Color4ub(255, 255, 255, 255)
	matTransform := c.GetModelMatrix()
	raylib.DrawMesh(*c.mdl.Meshes, *c.mdl.Materials, matTransform)

	return cmds
}

func (c *City) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if c == nil {
		return cmds
	}

	return cmds
}

func (c *City) Update(dt float32) {
	if c == nil {
		return
	}
}

func (c *City) GetCollider() pub_object.Collider {
	if c == nil {
		return nil
	}

	return nil
}

func (c *City) OnAdd() {
	if c == nil {
		return
	}
}

func (c *City) OnRemove() {
	if c == nil {
		return
	}
}

func (c *City) AddChild(obj pub_object.Object) {
	if c == nil {
		return
	}
}

func (c *City) RemChild(obj pub_object.Object) {
	if c == nil {
		return
	}
}

func (c *City) GetChilds() []pub_object.Object {
	if c == nil {
		return []pub_object.Object{}
	}

	return []pub_object.Object{}
}

func (c *City) GetModelMatrix() raylib.Matrix {
	if c == nil {
		return raylib.Matrix{}
	}

	matScale := raylib.MatrixScale(c.scale.X, c.scale.Y, c.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(c.GetPitch()), float64(c.GetYaw()), float64(c.GetRoll()))
	matRotation := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := raylib.MatrixTranslate(c.pos.X, c.pos.Y, c.pos.Z)
	matTransform := raylib.MatrixMultiply(raylib.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = raylib.MatrixMultiply(c.mdl.Transform, matTransform)
	return matTransform
}

func (c *City) GetModel() *raylib.Model {
	if c == nil {
		return nil
	}

	return c.mdl
}

func (c *City) SetColor(col color.Color) {
	if c == nil {
		return
	}
}

func (c *City) GetColor() color.Color {
	if c == nil {
		return nil
	}

	return raylib.White
}

func (c *City) GetScale() raylib.Vector3 {
	if c == nil {
		return raylib.Vector3{}
	}

	return c.scale
}

func (c *City) SetScale(sc raylib.Vector3) {
	if c == nil {
		return
	}

	c.scale = sc
}

func (c *City) SetPos(pos raylib.Vector3) {
	if c == nil {
		return
	}

	c.pos = pos
}

func (c *City) GetPos() raylib.Vector3 {
	if c == nil {
		return raylib.Vector3{}
	}

	return c.pos
}

func (c *City) GetPitch() float32 {
	if c == nil {
		return 0
	}

	return c.rot.X
}

func (c *City) SetPitch(pitch float32) {
	if c == nil {
		return
	}

	c.rot.X = pitch
}

func (c *City) GetYaw() float32 {
	if c == nil {
		return 0
	}

	return c.rot.Y
}

func (c *City) SetYaw(yaw float32) {
	if c == nil {
		return
	}

	c.rot.Y = yaw
}

func (c *City) GetRoll() float32 {
	if c == nil {
		return 0
	}

	return c.rot.Z
}

func (c *City) SetRoll(roll float32) {
	if c == nil {
		return
	}

	c.rot.Z = roll
}

func (c *City) GetVertices() []raylib.Vector3 {
	if c == nil {
		return []raylib.Vector3{}
	}

	verts := []raylib.Vector3{}
	length := c.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, raylib.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (c *City) GetUVs() []raylib.Vector2 {
	if c == nil {
		return []raylib.Vector2{}
	}

	uvs := []raylib.Vector2{}
	length := c.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, raylib.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (c *City) SetUVs(uvs []raylib.Vector2) {
	if c == nil {
		return
	}

	length := int(c.mdl.Meshes.VertexCount)
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Texcoords))
	header.Len = length * 2
	header.Cap = length * 2

	for i := 0; i < len(uvs); i++ {
		mdluvs[i*2] = uvs[i].X
		mdluvs[i*2+1] = uvs[i].Y
	}
	pub_object.UpdateModelUVs(c.mdl)
}

func (c *City) GetMaterials() *raylib.Material {
	if c == nil {
		return nil
	}

	return c.mdl.Materials
}

func (c *City) SetTexture(tex raylib.Texture2D) {
	if c == nil {
		return
	}

	*c.tex = tex
}

func (c *City) GetTexture() *raylib.Texture2D {
	if c == nil {
		return nil
	}

	return c.tex
}

func (c *City) GetHeightMap() *image.RGBA {
	if c == nil {
		return nil
	}

	return c.hm
}
