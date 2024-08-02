package character

import (
	"image/color"

	"karalis/internal/camera"
	"karalis/pkg/app"
	"karalis/pkg/input"
	"karalis/pkg/object"

	lmath "karalis/pkg/lmath"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Player struct {
	cam  *camera.Cam
	char *Character

	pos   raylib.Vector3
	rot   raylib.Vector3
	scale raylib.Vector3

	rchan chan (int)

	mode    int
	capture bool
}

func NewPlayer() (p *Player, err error) {
	p = &Player{}
	err = p.Init()

	return p, err
}

func (p *Player) Init() (err error) {
	p.pos = raylib.NewVector3(0, 0, 0)
	p.rot = raylib.NewVector3(0, 0, 0)
	p.scale = raylib.NewVector3(1, 1, 1)

	p.rchan = make(chan int)

	p.char, err = NewCharacter()
	if err != nil {
		return err
	}

	p.cam, err = camera.NewCam()
	if err != nil {
		return err
	}

	capt := p.ToggleCapture
	err = input.RegisterAction("ToggleMouseCapture", &capt, nil, true)
	if err != nil {
		return err
	}

	view := p.ToggleView
	err = input.RegisterAction("ToggleViewMode", &view, nil, true)
	if err != nil {
		return err
	}

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

	//update cant directly set mouse position
	select {
	case <-p.rchan:
		raylib.SetMousePosition(int(float32(app.CurApp.GetWidth())/2), int(float32(app.CurApp.GetHeight())/2))
	default:
	}

	return cmds
}

func (p *Player) Postrender(cam *camera.Cam) []func() {
	cmds := p.char.Postrender(cam)
	return cmds
}

func (p *Player) Update(dt float32) {
	p.char.Update(dt)
}

func (p *Player) Collide(data object.CollisionData) {
}

func (p *Player) RegCollideHandler(handler func(object.CollisionData) bool) {
}

func (p *Player) GetCollidable() []object.Object {
	return p.char.GetCollidable()
}

func (p *Player) GetCollider() object.Collider {
	return p.char.GetCollider()
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

func (p *Player) GetChilds() []object.Object {
	return p.char.GetChilds()
}

func (p *Player) GetModelMatrix() raylib.Matrix {
	matScale := raylib.MatrixScale(p.scale.X, p.scale.Y, p.scale.Z)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(p.GetPitch()), float64(p.GetYaw()), float64(p.GetRoll()))
	matRotation := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	matTranslation := raylib.MatrixTranslate(p.pos.X, p.pos.Y, p.pos.Z)
	matTransform := raylib.MatrixMultiply(raylib.MatrixMultiply(matScale, matRotation), matTranslation)
	return matTransform
}

func (p *Player) SetColor(col color.Color) {
}

func (p *Player) GetColor() color.Color {
	return raylib.White
}

func (p *Player) SetScale(sc raylib.Vector3) {
}

func (p *Player) GetScale() raylib.Vector3 {
	return raylib.NewVector3(1, 1, 1)
}

func (p *Player) SetPos(pos raylib.Vector3) {
	p.pos = pos
}

func (p *Player) GetPos() raylib.Vector3 {
	return p.pos
}

func (p *Player) GetPitch() float32 {
	return p.rot.X
}

func (p *Player) SetPitch(pi float32) {
	p.rot.X = pi
}

func (p *Player) GetYaw() float32 {
	return p.rot.Y
}

func (p *Player) SetYaw(y float32) {
	p.rot.Y = y
}

func (p *Player) GetRoll() float32 {
	return p.rot.Z
}

func (p *Player) SetRoll(r float32) {
	p.rot.Z = r
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

func (p *Player) GetMaterials() *raylib.Material {
	return p.char.GetMaterials()
}

func (p *Player) SetTexture(tex raylib.Texture2D) {
	p.char.SetTexture(tex)
}

func (p *Player) GetTexture() raylib.Texture2D {
	return p.char.GetTexture()
}

func (p *Player) GetCam() *camera.Cam {
	return p.cam
}

func (p *Player) CaptureMouse() {
	//can't directly call from update
	// raylib.SetMousePosition(int(float32(app.CurApp.GetWidth())/2), int(float32(app.CurApp.GetHeight())/2))
	p.rchan <- 1

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
		mpos := raylib.GetMousePosition()
		raylib.DisableCursor()
		raylib.HideCursor()

		zoom = float32(raylib.GetMouseWheelMove()) * dt * 20

		dx = dt * 20 * raylib.Deg2rad * (float32(app.CurApp.GetWidth())/2 - mpos.X)
		dy = dt * 20 * raylib.Deg2rad * (float32(app.CurApp.GetHeight())/2 - mpos.Y)

		//cant directly call due to bug using this in a subroutine on windows
		// raylib.SetMousePosition(int(float32(app.CurApp.GetWidth())/2), int(float32(app.CurApp.GetHeight())/2))
		p.rchan <- 1
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
