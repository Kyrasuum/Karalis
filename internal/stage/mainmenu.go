package stage

import (
	"godev/internal/cell"
	"godev/pkg/object"
)

var ()

type MainMenu struct{}

func (s *MainMenu) Init() {
}

func (s *MainMenu) OnResize(w int32, h int32) {
}

func (s *MainMenu) Prerender() []func() {
	return []func(){}
}

func (s *MainMenu) Render() []func() {
	return []func(){}
}

func (s *MainMenu) Postender() []func() {
	return []func(){}
}

func (s *MainMenu) Update(dt float32) {
}

func (s *MainMenu) OnInput(dt float32) {
}

func (s *MainMenu) OnAdd() {
}

func (s *MainMenu) OnRemove() {
}

func (s *MainMenu) GetPlayer() object.Object {
	return nil
}
func (s *MainMenu) GetCurrentCell() *cell.Cell {
	return nil
}
