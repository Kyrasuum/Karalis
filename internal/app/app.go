package app

import (
	"image/color"
	"time"

	"godev/internal/stage"
	App "godev/pkg/app"
	"godev/pkg/config"
	"godev/pkg/input"
	pub_stage "godev/pkg/stage"

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

func (a *app) defaults() {
	a.curStage = nil
	a.console = nil

	a.width = 800
	a.height = 512

	a.logicInterval = 16
	a.drawInterval = 16
}

func (a *app) init() {
	a.defaults()
	App.CurApp = a
}

func (a *app) handleInput(dt float32) {
	if a.curStage != nil {
		a.curStage.OnInput(dt)
	}
}

func (a *app) render() {
	if a.curStage != nil {
		cmds := a.curStage.Prerender()
		for _, cmd := range cmds {
			cmd()
		}
	}

	raylib.BeginDrawing()
	raylib.ClearBackground(color.RGBA{0, 0, 0, 1})
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

func (a *app) update(dt float32) {
	if a.curStage != nil {
		a.curStage.Update(dt)
	}
}

func (a *app) GetWidth() int32 {
	return a.width
}

func (a *app) GetHeight() int32 {
	return a.height
}

func (a *app) Running() bool {
	return !raylib.WindowShouldClose()
}

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

func (a *app) SetStage(nextStage pub_stage.Stage) {
	if a.curStage != nil {
		a.curStage.OnRemove()
	}
	if nextStage != nil {
		a.curStage = nextStage
		a.curStage.OnAdd()
	}
}

func (a *app) Exit() {
	if a.curStage != nil {
		a.curStage.OnRemove()
	}
	raylib.CloseWindow()
}

func (a *app) Start() {
	a.run()
}

func NewApp() app {
	a := app{}
	a.init()
	return a
}
