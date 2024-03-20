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

func (s *Cam) Init() error {
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

	return nil
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
	ql := lmath.Quat{}
	ql = *ql.FromEuler(0-float64(s.pitch), float64(s.yaw), float64(s.roll))

	view := ql.RotateVec3(lmath.Vec3{0, 0, float64(s.dist)})
	view.X += float64(s.camera.Target.X)
	view.Y -= float64(s.camera.Target.Y)
	view.Z += float64(s.camera.Target.Z)

	camMat := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(ql.X), float32(ql.Y), float32(ql.Z), float32(ql.W)))
	camMat = raylib.MatrixMultiply(camMat, raylib.MatrixTranslate(float32(view.X), float32(view.Y), float32(view.Z)))

	return camMat
}

func (s *Cam) GetCameraMatrix() raylib.Matrix {
	return raylib.GetCameraMatrix(s.camera)
}

func (s *Cam) GetWorldToScreen(pos raylib.Vector3) raylib.Vector2 {
	return raylib.GetWorldToScreen(pos, s.camera)
}

func (s *Cam) GetProjMatrix(width int32, height int32) raylib.Matrix {
	return raylib.GetCameraProjectionMatrix(&s.camera, float32(width)/float32(height))
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
