package shader

import (
	"fmt"

	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Compute struct {
	name string
	id   uint32
}

func (c *Compute) Init(shader string) error {
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

	shaderID := raylib.CompileShader(fmt.Sprintf("%s", code), raylib.ComputeShader)
	if shaderID == 0 {
		return fmt.Errorf("compute shader compile failed")
	}

	prog := raylib.LoadComputeShaderProgram(shaderID)
	if prog == 0 {
		return fmt.Errorf("compute program link failed")
	}

	c.id = prog
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

	raylib.EnableShader(c.id)
	return c.setUniform(uniform, val)
}

func (c *Compute) setUniform(uniform string, val interface{}) error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	loc := raylib.GetLocationUniform(c.id, uniform)
	if loc == -1 {
		return fmt.Errorf("Invalid uniform")
	}
	switch tval := val.(type) {
	case []float32:
		raylib.SetUniform(loc, tval, int32(raylib.ShaderUniformFloat), 1)
	case []int32:
		raylib.SetUniform(loc, tval, int32(raylib.ShaderUniformInt), 1)
	case float64:
		raylib.SetUniform(loc, []float32{float32(tval)}, int32(raylib.ShaderUniformFloat), 1)
	case float32:
		raylib.SetUniform(loc, []float32{tval}, int32(raylib.ShaderUniformFloat), 1)
	case int32:
		raylib.SetUniform(loc, []int32{tval}, int32(raylib.ShaderUniformInt), 1)
	case raylib.Vector2:
		raylib.SetUniform(loc, []float32{tval.X, tval.Y}, int32(raylib.ShaderUniformVec2), 1)
	case *raylib.Vector2:
		raylib.SetUniform(loc, []float32{tval.X, tval.Y}, int32(raylib.ShaderUniformVec2), 1)
	case raylib.Vector3:
		raylib.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z}, int32(raylib.ShaderUniformVec3), 1)
	case *raylib.Vector3:
		raylib.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z}, int32(raylib.ShaderUniformVec3), 1)
	case raylib.Vector4:
		raylib.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, int32(raylib.ShaderUniformVec4), 1)
	case *raylib.Vector4:
		raylib.SetUniform(loc, []float32{tval.X, tval.Y, tval.Z, tval.W}, int32(raylib.ShaderUniformVec4), 1)
	case raylib.Matrix:
		raylib.SetUniformMatrix(loc, tval)
	case raylib.Texture2D:
		raylib.SetUniformSampler(loc, tval.ID)
	default:
		return fmt.Errorf("Invalid uniform type %t", val)
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

	raylib.EnableShader(c.id)
	return nil
}

func (c *Compute) End() error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	if c.id == 0 {
		return fmt.Errorf("Invalid shader")
	}

	raylib.DisableShader()
	return nil
}

func (c *Compute) Dispatch(gx, gy, gz uint32) {
	raylib.ComputeShaderDispatch(gx, gy, gz)
}

func (c *Compute) OnRemove() error {
	if c == nil {
		return fmt.Errorf("Invalid shader")
	}

	if c.id != 0 {
		raylib.UnloadShaderProgram(c.id)
		c.id = 0
	}
	return nil
}
