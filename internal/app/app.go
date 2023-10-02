package app

import (
	"image/color"
	"time"

	"karalis/internal/stage"
	App "karalis/pkg/app"
	"karalis/pkg/config"
	"karalis/pkg/input"
	pub_stage "karalis/pkg/stage"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type app struct {
	curStage pub_stage.Stage
	console  interface{}

	logicInterval int64
	drawInterval  int64

	width  int32
	height int32
}

// initialize app
func (a *app) init() {
	a.curStage = nil
	a.console = nil

	a.width = 800
	a.height = 512

	a.logicInterval = 16
	a.drawInterval = 16

	App.CurApp = a
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
func (a *app) run() {
	raylib.SetConfigFlags(raylib.FlagWindowResizable)
	raylib.InitWindow(a.width, a.height, config.AppName)
	raylib.SetTargetFPS(int32(time.Second / (time.Duration(a.drawInterval) * time.Millisecond)))

	defer a.Exit()
	err := input.InitBindings()
	if err != nil {
		return
	}

	menu := stage.Game{}
	menu.Init()
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

// Exit the application
func (a *app) Exit() {
	if a.curStage != nil {
		a.curStage.OnRemove()
	}
	raylib.CloseWindow()
}

// start the application
func (a *app) Start() {
	a.run()
}

// create new application object
func NewApp() app {
	a := app{}
	a.init()
	return a
}
