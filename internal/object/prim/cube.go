package prim

import (
	"fmt"

	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

func NewCube() (p *Prim, err error) {
	p = &Prim{}
	p.init()

	mdl, err := res.GetRes("mdl/cube.obj")
	if err != nil {
		return nil, err
	}

	switch model := mdl.(type) {
	case raylib.Model:
		p.mdl = model
	default:
		return nil, fmt.Errorf("Invalid model object\n")
	}

	return p, nil
}
