package shader

import (
	"fmt"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type ComputeProgram struct {
	name      string
	ProgramID uint32
}

func (c *ComputeProgram) InitFromSource(name string) error {
	if c == nil {
		return fmt.Errorf("nil ComputeProgram")
	}
	shaderID := raylib.CompileShader(name, raylib.ComputeShader) // name may differ
	if shaderID == 0 {
		return fmt.Errorf("compute shader compile failed")
	}

	prog := raylib.LoadComputeShaderProgram(shaderID)
	if prog == 0 {
		return fmt.Errorf("compute program link failed")
	}

	c.name = name
	c.ProgramID = prog
	return nil
}

func (c *ComputeProgram) Begin() {
	raylib.EnableShader(c.ProgramID)
}

func (c *ComputeProgram) End() {
	raylib.DisableShader()
}

func (c *ComputeProgram) Dispatch(gx, gy, gz uint32) {
	raylib.ComputeShaderDispatch(gx, gy, gz)
}

func (c *ComputeProgram) Unload() {
	if c.ProgramID != 0 {
		raylib.UnloadShaderProgram(c.ProgramID)
		c.ProgramID = 0
	}
}
