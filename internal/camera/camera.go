package camera

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

var ()

type Cam struct {
	camera raylib.Camera

	dist  float32
	yaw   float32
	pitch float32
	roll  float32
}

func (s *Cam) Init() {
	s.camera = raylib.Camera3D{}
	s.camera.Position = raylib.NewVector3(0.0, 0.0, 0.0)
	s.camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	s.camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	s.camera.Fovy = 45.0
	s.camera.Projection = raylib.CameraPerspective

	s.dist = 4
	s.roll = 0
	s.pitch = raylib.Deg2rad * -30
	s.yaw = raylib.Deg2rad * 45
	s.UpdateCam()
}

func (s *Cam) OnResize(w int32, h int32) {
}

func (s *Cam) Prerender() []func() {
	return []func(){}
}

func (s *Cam) Render() []func() {
	raylib.BeginMode3D(s.camera)
	return []func(){raylib.EndMode3D}
}

func (s *Cam) Postrender() []func() {
	return []func(){}
}

func (s *Cam) Update(dt float32) {
}

func (s *Cam) GetPos() raylib.Vector3 {
	return s.camera.Position
}

func (s *Cam) SetPos(pos raylib.Vector3) {
	s.camera.Position = pos
}

func (s *Cam) GetTar() raylib.Vector3 {
	return s.camera.Target
}

func (s *Cam) SetTar(tar raylib.Vector3) {
	s.camera.Target = tar
}

func (s *Cam) GetModelMatrix() raylib.Matrix {
	camQuat := lmath.Quat{}
	camQuat = *camQuat.FromEuler(float64(s.GetPitch()), float64(s.GetYaw()), float64(s.GetRoll()))
	camMat := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(camQuat.X), float32(camQuat.Y), float32(camQuat.Z), float32(camQuat.W)))
	camMat = raylib.MatrixMultiply(camMat, raylib.MatrixTranslate(s.camera.Position.X, s.camera.Position.Y, s.camera.Position.Z))
	return camMat
}

func (s *Cam) GetWorldToScreen(pos raylib.Vector3, aspect float32) raylib.Vector2 {
	view := raylib.GetCameraViewMatrix(&s.camera)
	proj := raylib.GetCameraProjectionMatrix(&s.camera, aspect)
	pos = raylib.Vector3Transform(raylib.Vector3Transform(pos, view), proj)
	return raylib.NewVector2(pos.X/pos.Z, pos.Y/pos.Z)
}

func (s *Cam) GetDist() float32 {
	return s.dist
}

func (s *Cam) SetDist(d float32) {
	s.dist = d
}

func (s *Cam) GetPitch() float32 {
	return s.pitch
}

func (s *Cam) SetPitch(p float32) {
	s.pitch = p
}

func (s *Cam) GetYaw() float32 {
	return s.yaw
}

func (s *Cam) SetYaw(y float32) {
	s.yaw = y
}

func (s *Cam) GetRoll() float32 {
	return s.roll
}

func (s *Cam) SetRoll(r float32) {
	s.roll = r
}

func (s *Cam) ZoomCam(zoom float32) {
	s.dist -= zoom
}

func (s *Cam) RotateCam(pitch float32, yaw float32) {
	s.pitch += pitch
	s.yaw += yaw
}

func (s *Cam) MoveCam(move lmath.Vec3) {
	s.camera.Target.X += float32(move.X)
	s.camera.Target.Y += float32(move.Y)
	s.camera.Target.Z += float32(move.Z)
}

func (s *Cam) UpdateCam() {
	ql := lmath.Quat{}
	ql = *ql.FromEuler(float64(s.pitch), float64(s.yaw), float64(s.roll))

	view := ql.RotateVec3(lmath.Vec3{0, 0, float64(s.dist)})

	s.camera.Position.X = float32(view.X) + s.camera.Target.X
	s.camera.Position.Y = float32(view.Y) + s.camera.Target.Y
	s.camera.Position.Z = float32(view.Z) + s.camera.Target.Z
}

func (s *Cam) OnAdd() {
}

func (s *Cam) OnRemove() {
}
