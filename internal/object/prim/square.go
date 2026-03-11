package prim

import (
	"fmt"

	"karalis/internal/collider"
	"karalis/pkg/app"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
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
	case raylib.Model:
		p.mdl = model
	default:
		return nil, fmt.Errorf("Invalid model object\n")
	}

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

	mesh := raylib.GenMeshPlane(l, w, dl, dw)
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
