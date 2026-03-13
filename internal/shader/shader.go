package shader

import (
	"fmt"
	"runtime"
	"slices"
	"strings"

	"karalis/internal/rlx"
	"karalis/pkg/shader"
	"karalis/res"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	shaders = map[string]shader.Shader{}
)

type Shader struct {
	cleaner  *runtime.Cleanup
	shader   *rl.Shader
	name     string
	shaders  map[string]*rl.Shader
	defines  map[string]bool
	uniforms map[string]interface{}

	gsname string
	csname string
	esname string
	vsname string
	fsname string
}

func NewShader(shader string) (shader.Shader, error) {
	if s, ok := shaders[shader]; ok {
		return s, nil
	}

	s := &Shader{}
	err := s.init(shader)
	shaders[shader] = s

	return s, err
}

func (s *Shader) init(shader string) error {
	if s == nil {
		return nil
	}

	s.gsname = ""
	s.csname = ""
	s.esname = ""
	s.vsname = ""
	s.fsname = ""
	s.shaders = map[string]*rl.Shader{}
	s.defines = map[string]bool{}
	s.uniforms = map[string]interface{}{}
	s.name = shader
	err := s.genShader()
	if err != nil {
		return err
	}

	return nil
}

func (s *Shader) Extend(shader string) shader.Shader {
	if sh, ok := shaders[shader]; ok {
		return sh
	}
	if s == nil {
		return nil
	}

	ns := &Shader{}
	ns.init(s.name)
	ns.shaders = map[string]*rl.Shader{}
	ns.name = shader
	err := ns.genShader()
	if err != nil {
		return nil
	}
	shaders[shader] = ns

	return ns
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
		if len(strgs) > 0 {
			s.gsname = "shader/" + s.name + ".geom"
		} else if len(s.gsname) > 0 {
			strgs = s.loadShader(s.gsname)
		}

		strcs := s.loadShader("shader/" + s.name + ".ctrl")
		if len(strgs) > 0 {
			s.csname = "shader/" + s.name + ".ctrl"
		} else if len(s.csname) > 0 {
			strcs = s.loadShader(s.csname)
		}

		stres := s.loadShader("shader/" + s.name + ".eval")
		if len(stres) > 0 {
			s.esname = "shader/" + s.name + ".eval"
		} else if len(s.esname) > 0 {
			stres = s.loadShader(s.esname)
		}

		strvs := s.loadShader("shader/" + s.name + ".vert")
		if len(strvs) > 0 {
			s.vsname = "shader/" + s.name + ".vert"
		} else if len(s.vsname) > 0 {
			strvs = s.loadShader(s.vsname)
		}

		strfs := s.loadShader("shader/" + s.name + ".frag")
		if len(strfs) > 0 {
			s.fsname = "shader/" + s.name + ".frag"
		} else if len(s.fsname) > 0 {
			strfs = s.loadShader(s.fsname)
		}

		if len(strgs) > 0 && !s.GetDefine("GEOM_SHADER") {
			return s.SetDefine("GEOM_SHADER", true)
		}
		if len(strcs) > 0 && !s.GetDefine("CTRL_SHADER") {
			return s.SetDefine("CTRL_SHADER", true)
		}
		if len(stres) > 0 && !s.GetDefine("EVAL_SHADER") {
			return s.SetDefine("EVAL_SHADER", true)
		}

		shader := rlx.LoadShaderFromMemory(strvs, strcs, stres, strgs, strfs)
		s.shader = &shader
		s.shaders[key] = s.shader

		if s.cleaner != nil {
			s.cleaner.Stop()
		}
		cleaner := runtime.AddCleanup(s, func(shaders map[string]*rl.Shader) {
			for _, shader := range shaders {
				if shader != nil {
					rlx.UnloadShader(*shader)
				}
			}
		}, s.shaders)
		s.cleaner = &cleaner
	}

	for uniform, val := range s.uniforms {
		s.setUniform(uniform, val)
	}

	return nil
}

func (s *Shader) GetShader() *rl.Shader {
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

	loc = rlx.GetShaderLocation(*s.shader, uniform)
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

	rlx.BeginShaderMode(*s.shader)
	return s.setUniform(uniform, val)
}

func (s *Shader) setUniform(uniform string, val interface{}) error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	loc := rlx.GetShaderLocation(*s.shader, uniform)
	if loc == -1 {
		return fmt.Errorf("(%s)Invalid uniform: %s", s.name, uniform)
	}
	switch tval := val.(type) {
	case float64:
		rlx.SetShaderValue(*s.shader, loc, []float32{float32(tval)}, rl.ShaderUniformFloat)
	case float32:
		rlx.SetShaderValue(*s.shader, loc, []float32{float32(tval)}, rl.ShaderUniformFloat)
	case int32:
		rlx.SetShaderValue(*s.shader, loc, []float32{float32(tval)}, rl.ShaderUniformInt)
	case int:
		rlx.SetShaderValue(*s.shader, loc, []float32{float32(tval)}, rl.ShaderUniformInt)
	case uint32:
		rlx.SetShaderValue(*s.shader, loc, []float32{float32(tval)}, rl.ShaderUniformUint)
	case uint:
		rlx.SetShaderValue(*s.shader, loc, []float32{float32(tval)}, rl.ShaderUniformUint)
	case rl.Vector2:
		rlx.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y}, rl.ShaderUniformVec2)
	case *rl.Vector2:
		rlx.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y}, rl.ShaderUniformVec2)
	case rl.Vector3:
		rlx.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z}, rl.ShaderUniformVec3)
	case *rl.Vector3:
		rlx.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z}, rl.ShaderUniformVec3)
	case rl.Vector4:
		rlx.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, rl.ShaderUniformVec4)
	case *rl.Vector4:
		rlx.SetShaderValue(*s.shader, loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, rl.ShaderUniformVec4)
	case rl.Matrix:
		rlx.SetShaderValueMatrix(*s.shader, loc, tval)
	case rl.Texture2D:
		rlx.SetShaderValueTexture(*s.shader, loc, tval)
	default:
		return fmt.Errorf("(%s)Invalid uniform type %t", s.name, val)
	}
	s.uniforms[uniform] = val

	return nil
}

func (s *Shader) Begin() error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	if s.shader == nil {
		return fmt.Errorf("Invalid shader")
	}
	rlx.BeginShaderMode(*s.shader)
	return nil
}

func (s *Shader) End() error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	rlx.EndShaderMode()
	return nil
}

func (s *Shader) OnRemove() error {
	if s == nil {
		return fmt.Errorf("Invalid shader")
	}

	if s.shader != nil {
		rlx.UnloadShader(*s.shader)
		s.shader = nil
	}
	return nil
}
