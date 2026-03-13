package prim

import (
	"runtime"

	"karalis/internal/collider"
	"karalis/internal/rlx"
	"karalis/pkg/app"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

func NewCylinder(r, h float32, n int) (p *Prim, err error) {
	p = &Prim{}
	p.init()

	mesh := rlx.GenMeshCylinder(r, h, n)
	p.mdl = rlx.LoadModelFromMesh(mesh)
	if p.cleaner != nil {
		p.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(p, func(mdl rl.Model) {
		rlx.UnloadModel(mdl)
	}, p.mdl)
	p.cleaner = &cleaner

	sh := app.CurApp.GetShader()
	p.mdl.Materials.Shader = *sh.GetShader()

	col, err := collider.NewCollider(p)
	if err != nil {
		return nil, err
	}
	p.col = col

	return p, nil
}
