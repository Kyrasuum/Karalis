package prim

import (
	"fmt"
	"runtime"

	"karalis/internal/collider"
	"karalis/internal/rlx"
	"karalis/pkg/app"
	"karalis/res"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

func NewSquareMdl() (p *Prim, err error) {
	p = &Prim{}
	p.init()

	mdl, err := res.GetRes("mdl/square.obj")
	if err != nil {
		return nil, err
	}

	switch model := mdl.(type) {
	case rl.Model:
		p.mdl = model
	default:
		return nil, fmt.Errorf("Invalid model object\n")
	}
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

func NewSquare(l, w float32, dl, dw int) (p *Prim, err error) {
	p = &Prim{}
	p.init()

	mesh := rlx.GenMeshPlane(l, w, dl, dw)
	p.mdl = rlx.LoadModelFromMesh(mesh)
	if p.cleaner != nil {
		p.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(p, func(mdl rl.Model) {
		rl.UnloadModel(mdl)
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
