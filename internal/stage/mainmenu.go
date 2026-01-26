package stage

import (
	"fmt"

	"karalis/internal/cell"
	"karalis/pkg/object"
)

var ()

type MainMenu struct{}

func (s *MainMenu) Init() error {
	if s == nil {
		return fmt.Errorf("Invalid stage")
	}

	return nil
}

func (s *MainMenu) OnResize(w int32, h int32) {
	if s == nil {
		return
	}
}

func (s *MainMenu) Prerender() []func() {
	if s == nil {
		return []func(){}
	}

	return []func(){}
}

func (s *MainMenu) Render() []func() {
	if s == nil {
		return []func(){}
	}

	return []func(){}
}

func (s *MainMenu) Postender() []func() {
	if s == nil {
		return []func(){}
	}

	return []func(){}
}

func (s *MainMenu) Update(dt float32) {
	if s == nil {
		return
	}
}

func (s *MainMenu) OnInput(dt float32) {
	if s == nil {
		return
	}
}

func (s *MainMenu) OnAdd() {
	if s == nil {
		return
	}
}

func (s *MainMenu) OnRemove() {
	if s == nil {
		return
	}
}

func (s *MainMenu) GetPlayer() object.Object {
	if s == nil {
		return nil
	}

	return nil
}
func (s *MainMenu) GetCurrentCell() *cell.Cell {
	if s == nil {
		return nil
	}

	return nil
}
