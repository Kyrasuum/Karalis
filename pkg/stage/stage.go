package stage

import (
	"godev/internal/cell"
	"godev/pkg/object"
)

var ()

type Stage interface {
	Init()
	OnResize(w int32, h int32)
	Prerender() []func()
	Render() []func()
	Postrender() []func()
	Update(dt float32)
	OnInput(dt float32)
	OnAdd()
	OnRemove()
	GetPlayer() object.Object
	GetCurrentCell() *cell.Cell
}
