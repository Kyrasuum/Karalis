package object

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"

	"karalis/internal/camera"
	"karalis/internal/shader"
	pub_object "karalis/pkg/object"
	pub_shader "karalis/pkg/shader"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Terrain struct {
	tex *raylib.Texture2D
	shd pub_shader.Shader
	mdl raylib.Model
}

func (s *Terrain) Init() error {
	s.LoadImage(nil)
	s.LoadMap("")

	s.shd = &shader.Shader{}
	s.shd.Init("shader")

	return nil
}

func (s *Terrain) LoadImage(i interface{}) {
	var img *raylib.Image
	switch data := i.(type) {
	case string:
		tex, err := res.GetRes(data)
		if err != nil {
			s.LoadImage(nil)
			return
		}
		var img image.Image

		pos := strings.Index(data, ".") + 1
		ext := data[pos:]
		switch ext {
		case "png":
			img, err = png.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				s.LoadImage(nil)
				return
			}
		case "jpeg":
			img, err = jpeg.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				s.LoadImage(nil)
				return
			}
		default:
			img, _, err = image.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				s.LoadImage(nil)
				return
			}
		}
		if err != nil {
			s.LoadImage(nil)
			return
		}
		s.LoadImage(img)
	case image.Image:
		img = raylib.NewImageFromImage(data)
	case raylib.Color:
		img = raylib.GenImageColor(1536, 256, data)
	default:
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
	s.tex = &tex
}

func (s *Terrain) LoadMap(m string) {
	tex, err := res.GetRes(m)
	if err != nil {
		fmt.Printf("Error retrieving image: %+v\n", err)
		m = ""
	}
	var goimg image.Image

	if m != "" {
		pos := strings.Index(m, ".") + 1
		ext := m[pos:]
		switch ext {
		case "png":
			goimg, err = png.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				fmt.Printf("Error decoding image: %+v\n", err)
				return
			}
		case "jpeg":
			goimg, err = jpeg.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				fmt.Printf("Error decoding image: %+v\n", err)
				return
			}
		default:
			goimg, _, err = image.Decode(bytes.NewReader(tex.([]byte)))
			if err != nil {
				fmt.Printf("Error decoding image: %+v\n", err)
				return
			}
		}
	} else {
		width := 128
		height := 128
		color := color.RGBA{uint8(0), uint8(0), uint8(0), uint8(255)}
		cube := image.NewRGBA(image.Rect(0, 0, width, height))
		for i := range width {
			for j := range height {
				cube.Set(i, j, color)
			}
		}
		goimg = cube
	}

	img := raylib.NewImageFromImage(goimg)
	mesh := raylib.GenMeshHeightmap(*img, raylib.NewVector3(2, 2, 2))
	s.mdl = raylib.LoadModelFromMesh(mesh)
	raylib.SetMaterialTexture(s.mdl.Materials, raylib.MapDiffuse, *s.tex)
}

func (s *Terrain) GetModelMatrix() raylib.Matrix {
	return raylib.MatrixIdentity()
}

func (s *Terrain) GetPos() raylib.Vector3 {
	return raylib.Vector3{0, 0, 0}
}

func (s *Terrain) GetPitch() float32 {
	return 0
}

func (s *Terrain) SetPitch(pitch float32) {
}

func (s *Terrain) GetYaw() float32 {
	return 0
}

func (s *Terrain) SetYaw(yaw float32) {
}

func (s *Terrain) GetRoll() float32 {
	return 0
}

func (s *Terrain) SetRoll(roll float32) {
}

func (s *Terrain) GetVertices() []raylib.Vector3 {
	verts := []raylib.Vector3{}
	return verts
}

func (s *Terrain) GetUVs() []raylib.Vector2 {
	uvs := []raylib.Vector2{}
	return uvs
}

func (s *Terrain) SetUVs(uvs []raylib.Vector2) {
}

func (s *Terrain) GetMaterials() *raylib.Material {
	return nil
}

func (s *Terrain) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	*s.tex = tex
}

func (s *Terrain) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return *s.tex
}

func (s *Terrain) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}
	return cmds
}

func (s *Terrain) Render(cam *camera.Cam) []func() {
	cmds := []func(){}

	raylib.Color4ub(255, 255, 255, 255)
	matTransform := s.GetModelMatrix()
	s.mdl.Materials.Shader = *s.shd.GetShader()
	raylib.DrawMesh(*s.mdl.Meshes, *s.mdl.Materials, matTransform)

	return cmds
}

func (s *Terrain) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	return cmds
}

func (s *Terrain) Update(dt float32) {
}

func (s *Terrain) OnAdd() {
}

func (s *Terrain) OnRemove() {
}

func (s *Terrain) AddChild(obj pub_object.Object) {
}

func (s *Terrain) RemChild(obj pub_object.Object) {
}

func (s *Terrain) SetPos(pos raylib.Vector3) {
}
