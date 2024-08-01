package app

import (
	"fmt"
	"image/color"
	"time"

	"karalis/internal/shader"
	"karalis/internal/stage"
	App "karalis/pkg/app"
	"karalis/pkg/config"
	"karalis/pkg/input"
	pub_shader "karalis/pkg/shader"
	pub_stage "karalis/pkg/stage"
	"karalis/res"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type app struct {
	curStage  pub_stage.Stage
	curShader pub_shader.Shader
	console   interface{}

	logicInterval int64
	drawInterval  int64

	width  int32
	height int32
}

// initialize app
func (a *app) init() error {
	a.curStage = nil
	a.console = nil

	a.width = 800
	a.height = 512

	a.logicInterval = 16
	a.drawInterval = 16

	App.CurApp = a

	err := res.Init()
	if err != nil {
		return err
	}

	a.curShader = &shader.Shader{}

	return nil
}

// handle input
func (a *app) handleInput(dt float32) {
	if a.curStage != nil {
		a.curStage.OnInput(dt)
	}
}

// render cycle
func (a *app) render() {
	if a.curStage != nil {
		cmds := a.curStage.Prerender()
		for _, cmd := range cmds {
			cmd()
		}
	}

	raylib.BeginDrawing()
	raylib.ClearBackground(color.RGBA{0, 0, 0, 255})
	if a.curStage != nil {
		cmds := a.curStage.Render()
		for _, cmd := range cmds {
			cmd()
		}
	}

	if a.curStage != nil {
		cmds := a.curStage.Postrender()
		for _, cmd := range cmds {
			cmd()
		}
	}
	raylib.EndDrawing()
}

// handle resizing
func (a *app) onResize() {
	w := int32(raylib.GetScreenWidth())
	h := int32(raylib.GetScreenHeight())

	//check for resize event
	if a.width != w || a.height != h {
		if a.curStage != nil {
			a.curStage.OnResize(w, h)
		}
		a.width = w
		a.height = h
	}
}

// update cycle
func (a *app) update(dt float32) {
	if a.curStage != nil {
		a.curStage.Update(dt)
	}
}

// get window width
func (a *app) GetWidth() int32 {
	return a.width
}

// get window height
func (a *app) GetHeight() int32 {
	return a.height
}

// detect if app should continue running
func (a *app) Running() bool {
	return !raylib.WindowShouldClose()
}

// main run loop for the app while running
func (a *app) run() error {
	raylib.SetConfigFlags(raylib.FlagWindowResizable)
	raylib.InitWindow(a.width, a.height, config.AppName)
	raylib.SetTargetFPS(int32(time.Second / (time.Duration(a.drawInterval) * time.Millisecond)))

	defer a.Exit()
	err := input.InitBindings()
	if err != nil {
		return err
	}

	err = a.curShader.Init("shader")
	if err != nil {
		return err
	}

	menu := stage.Game{}
	err = menu.Init()
	if err != nil {
		return err
	}
	a.SetStage(&menu)

	//logic loop
	go func() {
		for a.Running() {
			dt := raylib.GetFrameTime()
			if raylib.IsCursorOnScreen() {
				a.handleInput(dt)
			}
			a.update(dt)
			a.onResize()
			time.Sleep(time.Duration(a.logicInterval) * time.Millisecond)
		}
	}()

	//render loop
	for a.Running() {
		a.render()
		time.Sleep(time.Duration(a.drawInterval) * time.Millisecond)
	}

	return nil
}

// set  the currently active stage
func (a *app) SetStage(nextStage pub_stage.Stage) {
	if a.curStage != nil {
		a.curStage.OnRemove()
	}
	if nextStage != nil {
		a.curStage = nextStage
		a.curStage.OnAdd()
	}
}

// get the currently active stage in the app
func (a *app) GetStage() pub_stage.Stage {
	return a.curStage
}

// get the currently active shader in the app
func (a *app) GetShader() pub_shader.Shader {
	return a.curShader
}

// Exit the application
func (a *app) Exit() {
	if a.curStage != nil {
		a.curStage.OnRemove()
	}
	if a.curShader != nil {
		a.curShader.OnRemove()
	}
	raylib.CloseWindow()
}

// start the application
func (a *app) Start() error {
	return a.run()
}

// create a new app
func NewApp() *app {
	a := &app{}
	err := a.init()
	if err != nil {
		fmt.Printf("ERR: %+v\n", err)
		return nil
	}
	return a
}
