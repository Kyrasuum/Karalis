package character

import (
	"karalis/internal/camera"
	"karalis/pkg/app"
	"karalis/pkg/input"
	"karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

var ()

type Player struct {
	cam  camera.Cam
	char Character

	lastmpos raylib.Vector2
	nextmpos raylib.Vector2

	mode    int
	capture bool
}

func (p *Player) Init() error {
	p.char.Init()

	p.cam = camera.Cam{}
	p.cam.Init()

	capt := p.ToggleCapture
	input.RegisterAction("ToggleMouseCapture", &capt, nil, true)
	view := p.ToggleView
	input.RegisterAction("ToggleViewMode", &view, nil, true)

	p.mode = 1
	p.capture = false

	return nil
}

func (p *Player) Prerender(cam *camera.Cam) []func() {
	cmds := p.char.Prerender(cam)
	return cmds
}

func (p *Player) Render(cam *camera.Cam) []func() {
	cmds := p.char.Render(cam)
	return cmds
}

func (p *Player) Postrender(cam *camera.Cam) []func() {
	cmds := p.char.Postrender(cam)
	return cmds
}

func (p *Player) Update(dt float32) {
	p.char.Update(dt)
}

func (p *Player) OnAdd() {
	p.char.OnAdd()
}

func (p *Player) OnRemove() {
	p.char.OnRemove()
}

func (p *Player) AddChild(obj object.Object) {
	p.char.AddChild(obj)
}

func (p *Player) RemChild(obj object.Object) {
	p.char.AddChild(obj)
}

// retrieve the portal display objects vertices
func (p *Player) GetVertices() []raylib.Vector3 {
	return p.char.GetVertices()
}

// retrieve the portal texture uvs for the display object
func (p *Player) GetUVs() []raylib.Vector2 {
	return p.char.GetUVs()
}

// set the texture uvs for the portal display object
func (p *Player) SetUVs(uvs []raylib.Vector2) {
	p.char.SetUVs(uvs)
}

func (p *Player) GetModelMatrix() raylib.Matrix {
	return p.char.GetModelMatrix()
}

func (p *Player) GetMaterials() *raylib.Material {
	return p.char.GetMaterials()
}

func (p *Player) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	p.char.SetTexture(mat, tex)
}

func (p *Player) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return p.char.GetTexture(mat)
}

func (p *Player) GetCam() *camera.Cam {
	return &p.cam
}

func (p *Player) CaptureMouse() {
	p.nextmpos = raylib.NewVector2(float32(app.CurApp.GetWidth())/2, float32(app.CurApp.GetHeight())/2)
	raylib.SetMousePosition(int(p.nextmpos.X), int(p.nextmpos.Y))
	p.lastmpos = p.nextmpos
	p.capture = true
}

func (p *Player) ReleaseMouse() {
	p.capture = false
	raylib.EnableCursor()
	raylib.ShowCursor()
}

func (p *Player) MouseCaptured() bool {
	return p.capture
}

func (p *Player) ToggleCapture() {
	if p.capture {
		p.ReleaseMouse()
	} else {
		p.CaptureMouse()
	}
}

func (p *Player) ToggleView() {
	if p.mode == 0 {
		p.mode = 1
		p.cam.SetDist(3)
	} else {
		p.mode = 0
		p.cam.SetDist(0.01)
	}
}

func (p *Player) OnInput(dt float32) {
	input.HandleInput("Character")

	var zoom float32 = 0
	var dx float32 = 0
	var dy float32 = 0

	move := lmath.Vec3{0, 0, 0}
	if input.Actions["MoveRight"].Pressed {
		move.X += float64(dt)
	}
	if input.Actions["MoveLeft"].Pressed {
		move.X -= float64(dt)
	}
	if input.Actions["MoveForward"].Pressed {
		move.Z -= float64(dt)
	}
	if input.Actions["MoveBackward"].Pressed {
		move.Z += float64(dt)
	}
	if input.Actions["MoveUp"].Pressed {
		move.Y += float64(dt)
	}
	if input.Actions["MoveDown"].Pressed {
		move.Y -= float64(dt)
	}
	if input.Actions["MoveFast"].Pressed {
		move.X *= 3
		move.Y *= 3
		move.Z *= 3
	}

	if p.MouseCaptured() {
		p.nextmpos = raylib.GetMousePosition()
		raylib.DisableCursor()
		raylib.HideCursor()

		zoom = float32(raylib.GetMouseWheelMove()) * dt * 20

		dx = dt * 20 * raylib.Deg2rad * (p.lastmpos.X - p.nextmpos.X)
		dy = dt * 20 * raylib.Deg2rad * (p.lastmpos.Y - p.nextmpos.Y)

		p.lastmpos = raylib.NewVector2(float32(app.CurApp.GetWidth())/2, float32(app.CurApp.GetHeight())/2)
		raylib.SetMousePosition(int(p.lastmpos.X), int(p.lastmpos.Y))
	}
	p.updateCam(move, zoom, dx, dy)
}

func (p *Player) updateCam(move lmath.Vec3, zoom, dx, dy float32) {
	dist := p.cam.GetDist() - zoom
	if dist <= 0.01 || (dist < 1 && dist > 0.75) {
		dist = 0.01
		p.mode = 0
	} else if dist > 0.01 && dist <= 0.75 {
		dist = 1
		p.mode = 1
	} else {
		//orbit view
	}

	if p.mode == 0 {
		//first person
	}
	if p.mode == 1 {
		//orbit view
		dx *= -1
		dy *= -1
	}

	ql := lmath.Quat{}
	if p.mode == 0 {
		move = ql.FromEuler(float64(p.cam.GetPitch()), float64(p.cam.GetYaw()), float64(p.cam.GetRoll())).RotateVec3(lmath.Vec3{float64(move.X), float64(move.Y), float64(move.Z)})
	}
	if p.mode == 1 {
		move = ql.FromEuler(0, float64(p.cam.GetYaw()), 0).RotateVec3(lmath.Vec3{float64(move.X), float64(move.Y), float64(move.Z)})
	}

	p.cam.SetDist(dist)
	p.cam.RotateCam(dy, dx)
	p.cam.MoveCam(move)
	p.cam.UpdateCam()
}
