package prim

import (
	"karalis/internal/collider"
	"karalis/pkg/app"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

func NewHalfSphere(r float32, n, s int) (p *Prim, err error) {
	p = &Prim{}
	p.init()

	mesh := raylib.GenMeshHemiSphere(r, n, s)
	p.mdl = raylib.LoadModelFromMesh(mesh)

	sh := app.CurApp.GetShader()
	p.mdl.Materials.Shader = *sh.GetShader()

	col, err := collider.NewCollider(p)
	if err != nil {
		return nil, err
	}
	p.col = col

	return p, nil
}
