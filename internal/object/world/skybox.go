package world

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"runtime"
	"strings"

	"karalis/internal/rlx"
	"karalis/internal/shader"
	"karalis/res"

	pub_object "karalis/pkg/object"
	pub_shader "karalis/pkg/shader"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Skybox struct {
	parent  pub_object.Object
	cleaner *runtime.Cleanup

	tex *rl.Texture2D
	shd pub_shader.Shader

	uvs   [][]rl.Vector2
	verts [][]rl.Vector3
}

// constructor for skybox
func NewSkybox(i interface{}) (s *Skybox, err error) {
	s = &Skybox{}
	err = s.Init()
	if err == nil {
		s.LoadImage(i)
	}

	return s, err
}

func (s *Skybox) Init() error {
	if s == nil {
		return fmt.Errorf("Invalid skybox")
	}
	s.parent = nil
	s.LoadImage(nil)

	var err error
	s.shd, err = shader.NewShader("skybox")
	if err != nil {
		return err
	}

	points := [][]float32{
		[]float32{-1, -1, -1},
		[]float32{1, -1, -1},
		[]float32{1, 1, -1},
		[]float32{-1, 1, -1},
		[]float32{-1, -1, 1},
		[]float32{1, -1, 1},
		[]float32{1, 1, 1},
		[]float32{-1, 1, 1},
	}

	texCoords := [][]float32{
		[]float32{0, 0},
		[]float32{1, 0},
		[]float32{1, 1},
		[]float32{0, 1},
	}

	indices := [][]int{
		[]int{0, 1, 2, 3},
		[]int{1, 5, 6, 2},
		[]int{5, 4, 7, 6},
		[]int{4, 0, 3, 7},
		[]int{3, 2, 6, 7},
		[]int{4, 5, 1, 0},
	}

	for _, face := range indices {
		verts := []rl.Vector3{}
		uvs := []rl.Vector2{}
		for i, indice := range face {
			point := points[indice]
			vert := rl.NewVector3(point[0], point[1], point[2])
			verts = append(verts, vert)

			coord := texCoords[i]
			uv := rl.NewVector2(coord[0], coord[1])
			uvs = append(uvs, uv)
		}
		s.verts = append(s.verts, verts)
		s.uvs = append(s.uvs, uvs)
	}

	return nil
}

func (s *Skybox) LoadImage(i interface{}) {
	if s == nil {
		return
	}

	var img *rl.Image
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
		return
	case image.Image:
		img = rlx.NewImageFromImage(data)
	case rl.Color:
		img = rlx.GenImageColor(1536, 256, data)
	default:
		width := 1536
		height := 256
		colors := []color.RGBA{
			color.RGBA{uint8(255), uint8(0), uint8(0), uint8(255)},
			color.RGBA{uint8(0), uint8(255), uint8(0), uint8(255)},
			color.RGBA{uint8(0), uint8(0), uint8(255), uint8(255)},
			color.RGBA{uint8(255), uint8(255), uint8(0), uint8(255)},
			color.RGBA{uint8(0), uint8(255), uint8(255), uint8(255)},
			color.RGBA{uint8(255), uint8(0), uint8(255), uint8(255)},
		}
		cube := image.NewRGBA(image.Rect(0, 0, width, height))
		for i := range 6 {
			for j := range width / 6 {
				for k := range height {
					cube.Set(i*width/6+j, k, colors[i])
				}
			}
		}
		img = rlx.NewImageFromImage(cube)
	}

	tex := rlx.LoadTextureCubemap(img, rl.CubemapLayoutAutoDetect)
	s.tex = &tex
	if s.cleaner != nil {
		s.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(s, func(tex rl.Texture2D) {
		rlx.UnloadTexture(tex)
	}, tex)
	s.cleaner = &cleaner
}

func (s *Skybox) GetModelMatrix() rl.Matrix {
	if s == nil {
		return rl.Matrix{}
	}

	return rl.MatrixIdentity()
}

func (s *Skybox) GetModel() *rl.Model {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Skybox) SetColor(col color.Color) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetColor() color.Color {
	if s == nil {
		return nil
	}

	return rl.White
}

func (s *Skybox) GetScale() rl.Vector3 {
	if s == nil {
		return rl.Vector3{}
	}

	return rl.Vector3{1, 1, 1}
}

func (s *Skybox) SetScale(sc rl.Vector3) {
	if s == nil {
		return
	}
}

func (s *Skybox) SetPos(pos rl.Vector3) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetPos() rl.Vector3 {
	if s == nil {
		return rl.Vector3{}
	}

	return rl.Vector3{0, 0, 0}
}

func (s *Skybox) GetPitch() float32 {
	if s == nil {
		return 0
	}

	return 0
}

func (s *Skybox) SetPitch(pitch float32) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetYaw() float32 {
	if s == nil {
		return 0
	}

	return 0
}

func (s *Skybox) SetYaw(yaw float32) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetRoll() float32 {
	if s == nil {
		return 0
	}

	return 0
}

func (s *Skybox) SetRoll(roll float32) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetVertices() []rl.Vector3 {
	if s == nil {
		return []rl.Vector3{}
	}

	verts := []rl.Vector3{}
	return verts
}

func (s *Skybox) GetUVs() []rl.Vector2 {
	if s == nil {
		return []rl.Vector2{}
	}

	uvs := []rl.Vector2{}
	return uvs
}

func (s *Skybox) SetUVs(uvs []rl.Vector2) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetMaterials() *rl.Material {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Skybox) SetTexture(tex rl.Texture2D) {
	if s == nil {
		return
	}

	*s.tex = tex
}

func (s *Skybox) GetTexture() *rl.Texture2D {
	if s == nil {
		return nil
	}

	return s.tex
}

func (s *Skybox) Prerender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if s == nil {
		return cmds
	}

	return cmds
}

func (s *Skybox) Render(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if s == nil {
		return cmds
	}

	rlx.DisableDepthMask()
	rlx.DisableDepthTest()
	rlx.PushMatrix()
	rlx.Begin(rl.Quads)
	rlx.EnableTextureCubemap(s.tex.ID)
	s.shd.Begin()
	s.shd.SetUniform("matView", rlx.GetMatrixModelview())
	s.shd.SetUniform("matProjection", rlx.GetMatrixProjection())

	rlx.Color4ub(255, 255, 255, 255)
	for i, quad := range s.verts {
		for j, vert := range quad {
			rlx.TexCoord2f(s.uvs[i][j].X, s.uvs[i][j].Y)
			rlx.Vertex3f(vert.X, vert.Y, vert.Z)
		}
	}

	s.shd.End()
	rlx.DisableTextureCubemap()
	rlx.End()
	rlx.PopMatrix()
	rlx.EnableDepthTest()
	rlx.EnableDepthMask()

	return cmds
}

func (s *Skybox) Postrender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if s == nil {
		return cmds
	}

	return cmds
}

func (s *Skybox) OnResize(w int32, h int32) {
	if s == nil {
		return
	}
}

func (s *Skybox) Update(dt float32) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetCollider() pub_object.Collider {
	if s == nil {
		return nil
	}

	return nil
}

func (s *Skybox) OnAdd(obj pub_object.Object) {
	if s == nil {
		return
	}
	s.parent = obj
}

func (s *Skybox) OnRemove() {
	if s == nil {
		return
	}
	s.parent = nil
}

func (s *Skybox) AddChild(obj pub_object.Object) {
	if s == nil {
		return
	}
}

func (s *Skybox) RemChild(obj pub_object.Object) {
	if s == nil {
		return
	}
}

func (s *Skybox) GetChilds() []pub_object.Object {
	if s == nil {
		return []pub_object.Object{}
	}

	return []pub_object.Object{}
}

func (s *Skybox) GetParent() pub_object.Object {
	if s == nil {
		return nil
	}
	return s.parent
}
