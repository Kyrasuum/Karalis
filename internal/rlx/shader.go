package rlx

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func CompileShader(shaderCode string, type_ int32) uint32 {
	return Call(func() uint32 {
		return rl.CompileShader(shaderCode, type_)
	})
}

func LoadComputeShaderProgram(shaderID uint32) uint32 {
	return Call(func() uint32 {
		return rl.LoadComputeShaderProgram(shaderID)
	})
}

func GetLocationUniform(shaderId uint32, uniformName string) int32 {
	return Call(func() int32 {
		return rl.GetLocationUniform(shaderId, uniformName)
	})
}

func SetUniformMatrix(locIndex int32, mat rl.Matrix) {
	Do(func() {
		rl.SetUniformMatrix(locIndex, mat)
	})
}

func SetUniformSampler(locIndex int32, texture uint32) {
	Do(func() {
		rl.SetUniformSampler(locIndex, texture)
	})
}

func ComputeShaderDispatch(groupsX, groupsY, groupsZ uint32) {
	Do(func() {
		rl.ComputeShaderDispatch(groupsX, groupsY, groupsZ)
	})
}

func UnloadShaderProgram(id uint32) {
	Do(func() {
		rl.UnloadShaderProgram(id)
	})
}

func EnableShader(shaderID uint32) {
	Do(func() {
		rl.EnableShader(shaderID)
	})
}

func DisableShader() {
	Do(func() {
		rl.DisableShader()
	})
}

func LoadShader(vsFileName string, csFileName string, esFileName string, gsFileName string, fsFileName string) rl.Shader {
	return Call(func() rl.Shader {
		return rl.LoadShader(vsFileName, csFileName, esFileName, gsFileName, fsFileName)
	})
}

func LoadShaderFromMemory(vsCode string, csCode string, esCode string, gsCode string, fsCode string) rl.Shader {
	return Call(func() rl.Shader {
		return rl.LoadShaderFromMemory(vsCode, csCode, esCode, gsCode, fsCode)
	})
}

func IsShaderValid(shader rl.Shader) bool {
	return Call(func() bool {
		return rl.IsShaderValid(shader)
	})
}

func GetShaderLocation(shader rl.Shader, uniformName string) int32 {
	return Call(func() int32 {
		return rl.GetShaderLocation(shader, uniformName)
	})
}

func GetShaderLocationAttrib(shader rl.Shader, attribName string) int32 {
	return Call(func() int32 {
		return rl.GetShaderLocationAttrib(shader, attribName)
	})
}

func SetShaderValue(shader rl.Shader, locIndex int32, value []float32, uniformType rl.ShaderUniformDataType) {
	buf := append([]float32(nil), value...)
	Do(func() {
		rl.SetShaderValue(shader, locIndex, buf, uniformType)
	})
}

func SetShaderValueV(shader rl.Shader, locIndex int32, value []float32, uniformType rl.ShaderUniformDataType, count int32) {
	buf := append([]float32(nil), value...)
	Do(func() {
		rl.SetShaderValueV(shader, locIndex, buf, uniformType, count)
	})
}

func SetShaderValueMatrix(shader rl.Shader, locIndex int32, mat rl.Matrix) {
	Do(func() {
		rl.SetShaderValueMatrix(shader, locIndex, mat)
	})
}

func SetShaderValueTexture(shader rl.Shader, locIndex int32, texture rl.Texture2D) {
	Do(func() {
		rl.SetShaderValueTexture(shader, locIndex, texture)
	})
}

func UnloadShader(shader rl.Shader) {
	Do(func() {
		rl.UnloadShader(shader)
	})
}

func SetUniform(locIndex int32, value any, uniformType, count int32) {
	Do(func() {
		rl.SetUniform(locIndex, value, uniformType, count)
	})
}

func LoadShaderBuffer(size uint32, data unsafe.Pointer, usageHint int32) uint32 {
	return Call(func() uint32 {
		return rl.LoadShaderBuffer(size, data, usageHint)
	})
}

func UpdateShaderBuffer(id uint32, data unsafe.Pointer, dataSize uint32, offset uint32) {
	Do(func() {
		rl.UpdateShaderBuffer(id, data, dataSize, offset)
	})
}

func BindShaderBuffer(id uint32, index uint32) {
	Do(func() {
		rl.BindShaderBuffer(id, index)
	})
}

func ReadShaderBuffer(id uint32, dest unsafe.Pointer, count uint32, offset uint32) {
	Do(func() {
		rl.ReadShaderBuffer(id, dest, count, offset)
	})
}

func UnloadShaderBuffer(id uint32) {
	Do(func() {
		rl.UnloadShaderBuffer(id)
	})
}
