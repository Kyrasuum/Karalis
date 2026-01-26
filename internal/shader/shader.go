package shader

import (
	"fmt"
	"slices"
	"strings"

	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Shader struct {
	shader   *raylib.Shader
	name     string
	shaders  map[string]*raylib.Shader
	defines  map[string]bool
	uniforms map[string]interface{}
}

func (s *Shader) Init(shader string) error {
	if s == nil {
		return nil
	}

	s.shaders = map[string]*raylib.Shader{}
	s.defines = map[string]bool{}
	s.uniforms = map[string]interface{}{}
	s.name = shader
	err := s.genShader()
	if err != nil {
		return err
	}

	return nil
}

func (s *Shader) shaderKey() string {
	if s == nil {
		return ""
	}

	keys := []string{}
	for define, val := range s.defines {
		if val {
			keys = append(keys, define)
		}
	}
	slices.Sort(keys)
	return strings.Join(keys, "")
}

func (s *Shader) loadShader(name string) string {
	if s == nil {
		return ""
	}

	strss := ""

	ss, err := res.GetRes(name)
	if err != nil {
		return strss
	}

	strss = string(ss.([]byte))
	posss := strings.Index(strss, "\n") + 1
	for define, _ := range s.defines {
		strss = strss[:posss] + "#define " + define + "\n" + strss[posss:]
	}

	return strss
}

func (s *Shader) genShader() error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	key := s.shaderKey()
	if shader, ok := s.shaders[key]; ok {
		s.shader = shader
	} else {
		strgs := s.loadShader("shader/" + s.name + ".geom")
		strcs := s.loadShader("shader/" + s.name + ".ctrl")
		stres := s.loadShader("shader/" + s.name + ".eval")
		strvs := s.loadShader("shader/" + s.name + ".vert")
		strfs := s.loadShader("shader/" + s.name + ".frag")

		shader := raylib.LoadShaderFromMemory(strvs, strcs, stres, strgs, strfs)
		s.shader = &shader
		s.shaders[key] = s.shader
	}

	for uniform, val := range s.uniforms {
		s.setUniform(uniform, val)
	}

	return nil
}

func (s *Shader) GetShader() *raylib.Shader {
	if s == nil {
		return nil
	}

	return s.shader
}

func (s *Shader) GetID() uint32 {
	if s == nil {
		return 0
	}

	return s.shader.ID
}

func (s *Shader) GetLocs() *int32 {
	if s == nil {
		return nil
	}

	return s.shader.Locs
}

func (s *Shader) GetLoc(uniform string) (loc int32, err error) {
	if s == nil {
		return 0, fmt.Errorf("Invalid shader")
	}

	loc = raylib.GetShaderLocation(*s.shader, uniform)
	if loc == -1 {
		return loc, fmt.Errorf("Invalid uniform")
	}
	return loc, nil
}

func (s *Shader) SetDefine(define string, val bool) error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	_, ok := s.defines[define]
	if ok && !val {
		delete(s.defines, define)
	} else {
		s.defines[define] = val
	}
	err := s.genShader()
	if err != nil {
		return err
	}
	return nil
}

func (s *Shader) GetDefine(define string) bool {
	if s == nil {
		return false
	}

	if _, ok := s.defines[define]; ok {
		return s.defines[define]
	}
	return false
}

func (s *Shader) SetUniform(uniform string, val interface{}) error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	switch val.(type) {
	case []float32, float64, float32, raylib.Vector2, raylib.Vector3, raylib.Vector4, *raylib.Vector2, *raylib.Vector3, *raylib.Vector4, raylib.Matrix, raylib.Texture2D:
	default:
		return fmt.Errorf("Invalid uniform type %t", val)
	}
	s.uniforms[uniform] = val
	s.setUniform(uniform, val)
	return nil
}

func (s *Shader) setUniform(uniform string, val interface{}) error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	loc := raylib.GetShaderLocation(*s.shader, uniform)
	if loc == -1 {
		return fmt.Errorf("Invalid uniform")
	}
	switch tval := val.(type) {
	case []float32:
		raylib.SetShaderValue(*s.shader, loc, tval, raylib.ShaderUniformFloat)
	case float64:
		raylib.SetShaderValue(*s.shader, loc, []float32{float32(tval)}, raylib.ShaderUniformFloat)
	case float32:
		raylib.SetShaderValue(*s.shader, loc, []float32{tval}, raylib.ShaderUniformFloat)
	case raylib.Vector2:
		raylib.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y}, raylib.ShaderUniformVec2)
	case *raylib.Vector2:
		raylib.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y}, raylib.ShaderUniformVec2)
	case raylib.Vector3:
		raylib.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z}, raylib.ShaderUniformVec3)
	case *raylib.Vector3:
		raylib.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z}, raylib.ShaderUniformVec3)
	case raylib.Vector4:
		raylib.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, raylib.ShaderUniformVec4)
	case *raylib.Vector4:
		raylib.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, raylib.ShaderUniformVec4)
	case raylib.Matrix:
		raylib.SetShaderValueMatrix(*s.shader, loc, tval)
	case raylib.Texture2D:
		raylib.SetShaderValueTexture(*s.shader, loc, tval)
	default:
		return fmt.Errorf("Invalid uniform type %t", val)
	}

	return nil
}

func (s *Shader) Begin() error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	if s.shader == nil {
		return fmt.Errorf("Invalid shader")
	}
	raylib.BeginShaderMode(*s.shader)
	return nil
}

func (s *Shader) End() error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	raylib.EndShaderMode()
	return nil
}

func (s *Shader) OnRemove() error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	if s.shader != nil {
		raylib.UnloadShader(*s.shader)
		s.shader = nil
	}
	return nil
}
