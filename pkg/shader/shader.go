package shader

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Shader interface {
	Init() error
	GetShader() *raylib.Shader
	GetID() uint32
	GetLocs() *int32
	GetLoc(uniform string) (loc int32, err error)
	SetDefine(string, bool) error
	SetUniform(string, interface{}) error
	Begin() error
	End() error
	OnRemove() error
}
