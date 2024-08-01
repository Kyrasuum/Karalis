package object

import (
	"bytes"
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

type Skybox struct {
	tex *raylib.Texture2D
	shd pub_shader.Shader

	uvs   [][]raylib.Vector2
	verts [][]raylib.Vector3
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
	s.LoadImage(nil)

	s.shd = &shader.Shader{}
	s.shd.Init("skybox")

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
		verts := []raylib.Vector3{}
		uvs := []raylib.Vector2{}
		for i, indice := range face {
			point := points[indice]
			vert := raylib.NewVector3(point[0], point[1], point[2])
			verts = append(verts, vert)

			coord := texCoords[i]
			uv := raylib.NewVector2(coord[0], coord[1])
			uvs = append(uvs, uv)
		}
		s.verts = append(s.verts, verts)
		s.uvs = append(s.uvs, uvs)
	}

	return nil
}

func (s *Skybox) LoadImage(i interface{}) {
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
		return
	case image.Image:
		img = raylib.NewImageFromImage(data)
	case raylib.Color:
		img = raylib.GenImageColor(1536, 256, data)
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
		img = raylib.NewImageFromImage(cube)
	}

	tex := raylib.LoadTextureCubemap(img, raylib.CubemapLayoutAutoDetect)
	s.tex = &tex
}

func (s *Skybox) GetModelMatrix() raylib.Matrix {
	return raylib.MatrixIdentity()
}

func (s *Skybox) SetColor(col color.Color) {
}

func (s *Skybox) GetColor() color.Color {
	return raylib.White
}

func (s *Skybox) GetScale() raylib.Vector3 {
	return raylib.Vector3{1, 1, 1}
}

func (s *Skybox) SetScale(sc raylib.Vector3) {
}

func (s *Skybox) SetPos(pos raylib.Vector3) {
}

func (s *Skybox) GetPos() raylib.Vector3 {
	return raylib.Vector3{0, 0, 0}
}

func (s *Skybox) GetPitch() float32 {
	return 0
}

func (s *Skybox) SetPitch(pitch float32) {
}

func (s *Skybox) GetYaw() float32 {
	return 0
}

func (s *Skybox) SetYaw(yaw float32) {
}

func (s *Skybox) GetRoll() float32 {
	return 0
}

func (s *Skybox) SetRoll(roll float32) {
}

func (s *Skybox) GetVertices() []raylib.Vector3 {
	verts := []raylib.Vector3{}
	return verts
}

func (s *Skybox) GetUVs() []raylib.Vector2 {
	uvs := []raylib.Vector2{}
	return uvs
}

func (s *Skybox) SetUVs(uvs []raylib.Vector2) {
}

func (s *Skybox) GetMaterials() *raylib.Material {
	return nil
}

func (s *Skybox) SetTexture(tex raylib.Texture2D) {
	*s.tex = tex
}

func (s *Skybox) GetTexture() raylib.Texture2D {
	return *s.tex
}

func (s *Skybox) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}
	return cmds
}

func (s *Skybox) Render(cam *camera.Cam) []func() {
	cmds := []func(){}

	raylib.DisableDepthMask()
	raylib.DisableDepthTest()
	raylib.PushMatrix()
	raylib.Begin(raylib.Quads)
	raylib.EnableTextureCubemap(s.tex.ID)
	s.shd.Begin()
	s.shd.SetUniform("matView", raylib.GetMatrixModelview())
	s.shd.SetUniform("matProjection", raylib.GetMatrixProjection())

	raylib.Color4ub(255, 255, 255, 255)
	for i, quad := range s.verts {
		for j, vert := range quad {
			raylib.TexCoord2f(s.uvs[i][j].X, s.uvs[i][j].Y)
			raylib.Vertex3f(vert.X, vert.Y, vert.Z)
		}
	}

	s.shd.End()
	raylib.DisableTextureCubemap()
	raylib.End()
	raylib.PopMatrix()
	raylib.EnableDepthTest()
	raylib.EnableDepthMask()

	return cmds
}

func (s *Skybox) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	return cmds
}

func (s *Skybox) Update(dt float32) {
}

func (s *Skybox) Collide(data pub_object.CollisionData) {
}

func (s *Skybox) RegCollideHandler(handler func(pub_object.CollisionData) bool) {
}

func (s *Skybox) CanCollide() bool {
	return false
}

func (s *Skybox) GetCollider() pub_object.Collider {
	col := pub_object.Collider{}
	return col
}

func (s *Skybox) OnAdd() {
}

func (s *Skybox) OnRemove() {
}

func (s *Skybox) AddChild(obj pub_object.Object) {
}

func (s *Skybox) RemChild(obj pub_object.Object) {
}

func (s *Skybox) GetChilds() []pub_object.Object {
	return []pub_object.Object{}
}
