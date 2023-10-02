package goray

/*
#include "raylib.h"
#include "rlgl.h"
#include "raymath.h"

// Update buffers
void UpdateModelUVs(Model* mdl) {
	UpdateMeshBuffer(mdl->meshes[0], 1, &mdl->meshes->texcoords[0], mdl->meshes->vertexCount*2*sizeof(float), 0);
}
*/
import "C"
import (
	"unsafe"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

func UpdateModelUVs(mdl *raylib.Model) {
	C.UpdateModelUVs((*C.Model)(unsafe.Pointer(mdl)))
}

func DrawRenderBatchActive() {
	C.rlDrawRenderBatchActive()
}

// import "github.com/go-gl/gl/v3.3-core/gl"
// gl.Enable(gl.STENCIL_TEST)
// gl.Disable(gl.STENCIL_TEST)
// gl.StencilFunc(gl.ALWAYS, 1, 0xff)
// gl.StencilFunc(gl.EQUAL, 1, 0xff)
// gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)
// gl.StencilOp(gl.KEEP, gl.KEEP, gl.KEEP)
// gl.ColorMask(false, false, false, false)
// gl.ColorMask(true, true, true, true)
