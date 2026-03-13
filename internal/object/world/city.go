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

type City struct {
	tex *rl.Texture2D
	mdl *rl.Model

	hm   *image.RGBA
	seed int64

	parent  pub_object.Object
	cleaner *runtime.Cleanup
	pos     rl.Vector3
	rot     rl.Vector3
	scale   rl.Vector3
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
	c.parent = nil
	c.pos = rl.NewVector3(0, 0, 0)
	c.rot = rl.NewVector3(0, 0, 0)
	c.scale = rl.NewVector3(1, 1, 1)

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

	var img *rl.Image
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
		img = rlx.NewImageFromImage(data)
	case rl.Color:
		img = rlx.GenImageColor(1536, 256, data)
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
		img = rlx.NewImageFromImage(cube)
	}
	tex := rlx.LoadTextureFromImage(img)
	c.tex = &tex
	if c.cleaner != nil {
		c.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(c, func(in []interface{}) {
		rlx.UnloadTexture(in[0].(rl.Texture2D))
	}, []interface{}{*c.tex})
	c.cleaner = &cleaner
	if c.mdl != nil {
		rlx.SetMaterialTexture(c.mdl.Materials, rl.MapDiffuse, *c.tex)
		c.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
		if c.cleaner != nil {
			c.cleaner.Stop()
		}
		cleaner := runtime.AddCleanup(c, func(in []interface{}) {
			rlx.UnloadTexture(in[0].(rl.Texture2D))
			rlx.UnloadModel(in[1].(rl.Model))
		}, []interface{}{*c.tex, *c.mdl})
		c.cleaner = &cleaner
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
	hm := rlx.NewImageFromImage(img)
	mesh := rlx.GenMeshHeightmap(*hm, rl.NewVector3(1, 1, 1))
	mdl := rlx.LoadModelFromMesh(mesh)

	c.mdl = &mdl
	if c.tex == nil {
		c.LoadImage(nil)
	}
	if c.cleaner != nil {
		c.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(c, func(in []interface{}) {
		rlx.UnloadTexture(in[0].(rl.Texture2D))
		rlx.UnloadModel(in[1].(rl.Model))
	}, []interface{}{*c.tex, *c.mdl})
	c.cleaner = &cleaner

	rlx.SetMaterialTexture(c.mdl.Materials, rl.MapDiffuse, *c.tex)
	c.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
}

func (c *City) Prerender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if c == nil {
		return cmds
	}

	return cmds
}

func (c *City) Render(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if c == nil {
		return cmds
	}

	rl.Color4ub(255, 255, 255, 255)
	matTransform := c.GetModelMatrix()
	rlx.DrawMesh(*c.mdl.Meshes, *c.mdl.Materials, matTransform)

	return cmds
}

func (c *City) Postrender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if c == nil {
		return cmds
	}

	return cmds
}

func (c *City) OnResize(w int32, h int32) {
	if c == nil {
		return
	}
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

func (c *City) OnAdd(obj pub_object.Object) {
	if c == nil {
		return
	}
	c.parent = obj
}

func (c *City) OnRemove() {
	if c == nil {
		return
	}
	c.parent = nil
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

func (c *City) GetModelMatrix() rl.Matrix {
	if c == nil {
		return rl.Matrix{}
	}

	matScale := rl.MatrixScale(c.scale.X, c.scale.Y, c.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(c.GetPitch()), float64(c.GetYaw()), float64(c.GetRoll()))
	matRotation := rl.QuaternionToMatrix(rl.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := rl.MatrixTranslate(c.pos.X, c.pos.Y, c.pos.Z)
	matTransform := rl.MatrixMultiply(rl.MatrixMultiply(matScale, matRotation), matTranslation)
	matTransform = rl.MatrixMultiply(c.mdl.Transform, matTransform)
	return matTransform
}

func (c *City) GetModel() *rl.Model {
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

	return rl.White
}

func (c *City) GetScale() rl.Vector3 {
	if c == nil {
		return rl.Vector3{}
	}

	return c.scale
}

func (c *City) SetScale(sc rl.Vector3) {
	if c == nil {
		return
	}

	c.scale = sc
}

func (c *City) SetPos(pos rl.Vector3) {
	if c == nil {
		return
	}

	c.pos = pos
}

func (c *City) GetPos() rl.Vector3 {
	if c == nil {
		return rl.Vector3{}
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

func (c *City) GetVertices() []rl.Vector3 {
	if c == nil {
		return []rl.Vector3{}
	}

	verts := []rl.Vector3{}
	length := c.mdl.Meshes.VertexCount

	var mdlverts []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdlverts))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Vertices))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdlverts); i++ {
		verts = append(verts, rl.NewVector3(mdlverts[3*i], mdlverts[3*i+1], mdlverts[3*i+2]))
	}
	return verts
}

func (c *City) GetUVs() []rl.Vector2 {
	if c == nil {
		return []rl.Vector2{}
	}

	uvs := []rl.Vector2{}
	length := c.mdl.Meshes.VertexCount
	var mdluvs []float32

	header := (*reflect.SliceHeader)(unsafe.Pointer(&mdluvs))
	header.Data = uintptr(unsafe.Pointer(c.mdl.Meshes.Texcoords))
	header.Len = int(length)
	header.Cap = int(length)

	for i := 0; i < len(mdluvs); i++ {
		uvs = append(uvs, rl.NewVector2(mdluvs[2*i], mdluvs[2*i+1]))
	}
	return uvs
}

func (c *City) SetUVs(uvs []rl.Vector2) {
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

func (c *City) GetMaterials() *rl.Material {
	if c == nil {
		return nil
	}

	return c.mdl.Materials
}

func (c *City) SetTexture(tex rl.Texture2D) {
	if c == nil {
		return
	}

	*c.tex = tex

	if c.cleaner != nil {
		c.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(c, func(in []interface{}) {
		rlx.UnloadTexture(in[0].(rl.Texture2D))
	}, []interface{}{*c.tex})
	c.cleaner = &cleaner
	if c.mdl != nil {
		rlx.SetMaterialTexture(c.mdl.Materials, rl.MapDiffuse, *c.tex)
		c.mdl.Materials.Shader = *app.CurApp.GetShader().GetShader()
		if c.cleaner != nil {
			c.cleaner.Stop()
		}
		cleaner := runtime.AddCleanup(c, func(in []interface{}) {
			rlx.UnloadTexture(in[0].(rl.Texture2D))
			rlx.UnloadModel(in[1].(rl.Model))
		}, []interface{}{*c.tex, *c.mdl})
		c.cleaner = &cleaner
	}
}

func (c *City) GetTexture() *rl.Texture2D {
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

func (c *City) GetParent() pub_object.Object {
	if c == nil {
		return nil
	}
	return c.parent
}
