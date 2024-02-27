package shader

import (
	"fmt"
	"karalis/res"
	"slices"
	"strings"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Shader struct {
	shader   *raylib.Shader
	shaders  map[string]*raylib.Shader
	defines  map[string]bool
	uniforms map[string]interface{}
}

func (s *Shader) Init() error {
	s.shaders = map[string]*raylib.Shader{}
	s.defines = map[string]bool{}
	s.uniforms = map[string]interface{}{}
	err := s.genShader()
	if err != nil {
		return err
	}

	return nil
}

func (s *Shader) shaderKey() string {
	keys := []string{}
	for define, val := range s.defines {
		if val {
			keys = append(keys, define)
		}
	}
	slices.Sort(keys)
	return strings.Join(keys, "")
}

func (s *Shader) genShader() error {
	key := s.shaderKey()
	if shader, ok := s.shaders[key]; ok {
		s.shader = shader
	} else {
		fs, err := res.GetRes("shader/shader.frag")
		if err != nil {
			return err
		}
		vx, err := res.GetRes("shader/shader.vert")
		if err != nil {
			return err
		}
		strvx := string(vx.([]byte))
		strfs := string(fs.([]byte))

		for define, _ := range s.defines {
			strvx = "#define " + define + "\n" + strvx
			strfs = "#define " + define + "\n" + strfs
		}

		strvx = "#version 330\n" + strvx
		strfs = "#version 330\n" + strfs

		shader := raylib.LoadShaderFromMemory(strvx, strfs)
		s.shader = &shader
		s.shaders[key] = &shader
	}

	for uniform, val := range s.uniforms {
		s.setUniform(uniform, val)
	}

	return nil
}

func (s *Shader) GetShader() *raylib.Shader {
	return s.shader
}

func (s *Shader) GetID() uint32 {
	return s.shader.ID
}

func (s *Shader) GetLocs() *int32 {
	return s.shader.Locs
}

func (s *Shader) GetLoc(uniform string) (loc int32, err error) {
	loc = raylib.GetShaderLocation(*s.shader, uniform)
	if loc == -1 {
		return loc, fmt.Errorf("Invalid uniform")
	}
	return loc, nil
}

func (s *Shader) SetDefine(define string, val bool) error {
	if def, ok := s.defines[define]; !ok || !def {
		s.defines[define] = val
	} else {
		if !val {
			delete(s.defines, define)
		} else {
			return nil
		}
	}
	err := s.genShader()
	if err != nil {
		return err
	}
	return nil
}

func (s *Shader) SetUniform(uniform string, val interface{}) error {
	switch val.(type) {
	case []float32, float64, float32, raylib.Matrix, raylib.Texture2D:
	default:
		return fmt.Errorf("Invalid uniform type %t", val)
	}
	s.uniforms[uniform] = val
	s.setUniform(uniform, val)
	return nil
}

func (s *Shader) setUniform(uniform string, val interface{}) error {
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
	if s.shader == nil {
		return fmt.Errorf("Invalid shader")
	}
	raylib.BeginShaderMode(*s.shader)
	return nil
}

func (s *Shader) End() error {
	raylib.EndShaderMode()
	return nil
}

func (s *Shader) OnRemove() error {
	if s.shader != nil {
		raylib.UnloadShader(*s.shader)
		s.shader = nil
	}
	return nil
}
