package prim

import (
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

func NewSquare() (p *Prim, err error) {
	p = &Prim{}
	p.init()

	mdl, err := res.GetRes("mdl/square.obj")
	if err != nil {
		return nil, err
	}
	p.mdl = mdl.(raylib.Model)

	return p, nil
}
