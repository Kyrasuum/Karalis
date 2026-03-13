package shader

import (
	"fmt"
	"runtime"

	"karalis/internal/rlx"
	"karalis/res"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	computes = map[string]Compute{}
)

type Compute struct {
	name string
	id   uint32

	cleaner *runtime.Cleanup
}

func NewCompute(shader string) (Compute, error) {
	if c, ok := computes[shader]; ok {
		return c, nil
	}

	c := Compute{}
	err := c.init(shader)
	computes[shader] = c

	return c, err
}

func (c *Compute) init(shader string) error {
	if c == nil {
		return nil
	}

	c.name = shader
	c.id = 0
	err := c.genShader()
	if err != nil {
		return err
	}

	return nil
}

func (c *Compute) genShader() error {
	if c == nil {
		return fmt.Errorf("nil ComputeProgram")
	}

	code, err := res.GetRes("shader/" + c.name + ".comp")
	if err != nil {
		return fmt.Errorf("invalid source code")
	}

	shaderID := rlx.CompileShader(fmt.Sprintf("%s", code), rl.ComputeShader)
	if shaderID == 0 {
		return fmt.Errorf("compute shader compile failed")
	}

	prog := rlx.LoadComputeShaderProgram(shaderID)
	if prog == 0 {
		return fmt.Errorf("compute program link failed")
	}

	c.id = prog
	if c.cleaner != nil {
		c.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(c, func(shader rl.Shader) {
		rlx.UnloadShader(shader)
	}, rl.Shader{ID: c.id})
	c.cleaner = &cleaner
	return nil
}

func (c *Compute) GetID() uint32 {
	if c == nil {
		return 0
	}

	return c.id
}

func (c *Compute) SetUniform(uniform string, val interface{}) error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	rlx.EnableShader(c.id)
	return c.setUniform(uniform, val)
}

func (c *Compute) setUniform(uniform string, val interface{}) error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	loc := rlx.GetLocationUniform(c.id, uniform)
	if loc == -1 {
		return fmt.Errorf("(%s)Invalid uniform: %s", c.name, uniform)
	}
	switch tval := val.(type) {
	case []float32:
		rlx.SetUniform(loc, tval, int32(rl.ShaderUniformFloat), int32(len(tval)))
	case []int32:
		rlx.SetUniform(loc, tval, int32(rl.ShaderUniformInt), int32(len(tval)))
	case []int:
		rlx.SetUniform(loc, tval, int32(rl.ShaderUniformInt), int32(len(tval)))
	case []uint32:
		rlx.SetUniform(loc, tval, int32(rl.ShaderUniformUint), int32(len(tval)))
	case []uint:
		rlx.SetUniform(loc, tval, int32(rl.ShaderUniformUint), int32(len(tval)))
	case float64:
		rlx.SetUniform(loc, []float32{float32(tval)}, int32(rl.ShaderUniformFloat), 1)
	case float32:
		rlx.SetUniform(loc, []float32{tval}, int32(rl.ShaderUniformFloat), 1)
	case int32:
		rlx.SetUniform(loc, []int32{tval}, int32(rl.ShaderUniformInt), 1)
	case int:
		rlx.SetUniform(loc, []int32{int32(tval)}, int32(rl.ShaderUniformInt), 1)
	case uint32:
		rlx.SetUniform(loc, []uint32{tval}, int32(rl.ShaderUniformUint), 1)
	case uint:
		rlx.SetUniform(loc, []uint32{uint32(tval)}, int32(rl.ShaderUniformUint), 1)
	case rl.Vector2:
		rlx.SetUniform(loc, []float32{tval.X, tval.Y}, int32(rl.ShaderUniformVec2), 1)
	case *rl.Vector2:
		rlx.SetUniform(loc, []float32{tval.X, tval.Y}, int32(rl.ShaderUniformVec2), 1)
	case rl.Vector3:
		rlx.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z}, int32(rl.ShaderUniformVec3), 1)
	case *rl.Vector3:
		rlx.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z}, int32(rl.ShaderUniformVec3), 1)
	case rl.Vector4:
		rlx.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, int32(rl.ShaderUniformVec4), 1)
	case *rl.Vector4:
		rlx.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, int32(rl.ShaderUniformVec4), 1)
	case rl.Matrix:
		rlx.SetUniformMatrix(loc, tval)
	case rl.Texture2D:
		rlx.SetUniformSampler(loc, tval.ID)
	default:
		return fmt.Errorf("(%s)Invalid uniform type: %t", c.name, val)
	}

	return nil
}

func (c *Compute) Begin() error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	if c.id == 0 {
		return fmt.Errorf("Invalid shader")
	}

	rlx.EnableShader(c.id)
	return nil
}

func (c *Compute) End() error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	if c.id == 0 {
		return fmt.Errorf("Invalid shader")
	}

	rlx.DisableShader()
	return nil
}

func (c *Compute) Dispatch(gx, gy, gz uint32) {
	rlx.ComputeShaderDispatch(gx, gy, gz)
}

func (c *Compute) OnRemove() error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	if c.id != 0 {
		rlx.UnloadShaderProgram(c.id)
		c.id = 0
	}
	return nil
}
