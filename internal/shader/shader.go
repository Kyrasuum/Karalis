package shader

import (
	"fmt"
	"unsafe"

	"karalis/pkg/app"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

/*
#include "shader.h"
*/
import "C"

type Shader struct {
	shader *raylib.Shader
}

func (s *Shader) Init() error {
	fs, err := res.GetRes("shader/shader.frag")
	if err != nil {
		return err
	}
	vx, err := res.GetRes("shader/shader.vert")
	if err != nil {
		return err
	}
	shader := raylib.LoadShaderFromMemory(string(vx.([]byte)), string(fs.([]byte)))
	s.shader = &shader
	C.SetLocs((*C.int)(unsafe.Pointer(shader.Locs)))

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

//export getloc
func getloc(uni *C.char) C.int {
	sh := app.CurApp.GetShader()
	uniform := C.GoString(uni)

	loc, err := sh.GetLoc(uniform)
	if loc == -1 || err != nil {
		return C.int(-1)
	}

	return C.int(loc)
}

func (s *Shader) GetLoc(uniform string) (loc int32, err error) {
	loc = raylib.GetShaderLocation(*s.shader, uniform)
	if loc == -1 {
		return loc, fmt.Errorf("Invalid uniform")
	}
	return loc, nil
}

func (s *Shader) SetUniform(uniform string, val interface{}) error {
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
