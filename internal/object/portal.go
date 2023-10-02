package object

import (
	"fmt"

	"karalis/internal/camera"
	"karalis/internal/cell"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Portal struct {
	scene *cell.Cell

	exit *Portal

	obj pub_object.Object
	cam *camera.Cam

	rendering bool
	visible   bool
	showexit  bool
}

// initialize portal object
func (p *Portal) Init(scene *cell.Cell, exit *Portal, cam *camera.Cam, obj pub_object.Object) {
	if scene != nil {
		p.scene = scene
	} else {
		p.scene = &cell.Cell{}
		p.scene.Init()
	}

	if exit != nil {
		p.exit = exit
	} else {
		p.exit = nil
	}

	if cam != nil {
		p.cam = cam
	} else {
		p.cam = &camera.Cam{}
		p.cam.Init()
	}

	if obj != nil {
		p.obj = obj
	} else {
		p.obj = nil
	}

	p.rendering = false
	p.visible = true
}

// get exit portal pair
func (p *Portal) GetPair() *Portal {
	return p.exit
}

// set exit portal pair
func (p *Portal) Pair(e *Portal) {
	if e.exit != nil {
		e.exit.exit = nil
	}
	if p.exit != nil {
		p.exit.exit = nil
	}
	e.exit = p
	p.exit = e
}

// retrieve the portal display objects vertices
func (p *Portal) GetVertices() []raylib.Vector3 {
	if p.obj != nil {
		return p.obj.GetVertices()
	} else {
		return []raylib.Vector3{}
	}
}

// retrieve the portal texture uvs for the display object
func (p *Portal) GetUVs() []raylib.Vector2 {
	if p.obj != nil {
		return p.obj.GetUVs()
	} else {
		return []raylib.Vector2{}
	}
}

// set the texture uvs for the portal display object
func (p *Portal) SetUVs(uvs []raylib.Vector2) {
	if p.obj != nil {
		p.obj.SetUVs(uvs)
	}
}

// get the portal render objects model matrix
func (p *Portal) GetModelMatrix() raylib.Matrix {
	if p.obj != nil {
		return p.obj.GetModelMatrix()
	} else {
		return raylib.MatrixTranslate(0, 0, 0)
	}
}

// get portal render material
func (p *Portal) GetMaterials() *raylib.Material {
	return &raylib.Material{}
}

// set portal render texture
func (p *Portal) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
}

// get portal render texture
func (p *Portal) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return raylib.Texture2D{}
}

// set portal render object
func (p *Portal) SetPortal(obj pub_object.Object) {
	if p.obj != nil {
		p.obj.OnRemove()
	}
	p.obj = obj
	p.obj.OnAdd()
	p.obj.SetTexture(p.obj.GetMaterials(), p.GetTexture(nil))
}

// get portal render object
func (p *Portal) GetPortal() pub_object.Object {
	return p.obj
}

// set camera for portal
func (p *Portal) SetCam(obj *camera.Cam) {
	p.cam = obj
}

// return camera for portal
func (p *Portal) GetCam() *camera.Cam {
	return p.cam
}

// prerender hook
func (p *Portal) Prerender(cam *camera.Cam) []func() {
	cmds := []func(){}

	return cmds
}

// render hook
func (p *Portal) Render(cam *camera.Cam) []func() {
	cmds := []func(){}

	//guards to ensure we only render when we should be
	if !p.rendering && p.visible && p.obj != nil {
		//prevent rerendering a portal a second time
		p.rendering = true
		if p.exit != nil {
			p.exit.visible = false
		}

		//calculate portal camera position based on calling render camera
		camMdl := cam.GetModelMatrix()
		portalmdl := p.GetModelMatrix()
		wldToLcl := raylib.MatrixInvert(portalmdl)
		lclToWld := raylib.MatrixRotateX(180)
		if p.exit != nil {
			lclToWld = p.exit.GetModelMatrix()
		}
		transform := raylib.MatrixMultiply(lclToWld, raylib.MatrixMultiply(wldToLcl, camMdl))
		tar := cam.GetTar()
		p.cam.SetPos(raylib.Vector3Transform(raylib.NewVector3(0, 0, 0), transform))
		p.cam.SetTar(raylib.NewVector3(float32(tar.X), 0-float32(tar.Y), float32(tar.Z)))

		//render portal object for stencil
		raylib.DisableDepthMask()
		cmds = p.obj.Render(cam)

		//render from portals perspective
		cmds = p.cam.Prerender()
		cmds = append(cmds, p.scene.Prerender(p.cam)...)
		for _, cmd := range cmds {
			cmd()
		}

		cmds = p.cam.Render()
		cmds = append(cmds, p.scene.Render(p.cam)...)
		for _, cmd := range cmds {
			cmd()
		}

		cmds = p.cam.Postrender()
		cmds = append(cmds, p.scene.Postrender(p.cam)...)
		for _, cmd := range cmds {
			cmd()
		}

		p.rendering = false
		if p.exit != nil {
			p.exit.visible = true
		}

		//render portal object for depth
		raylib.EnableDepthMask()
		cmds = p.obj.Render(cam)
	}

	return cmds
}

// postrender hook
func (p *Portal) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}

	return cmds
}

// handle update cycle
func (p *Portal) Update(dt float32) {
	if p.obj != nil {
		p.obj.Update(dt)
	}
	p.scene.Update(dt)
}

// handle add event
func (p *Portal) OnAdd() {
	p.scene.OnAdd()
}

// handle remove event
func (p *Portal) OnRemove() {
	p.scene.OnRemove()
}

// handle resize event
func (p *Portal) OnResize(w int32, h int32) {
	p.cam.OnResize(w, h)
	fmt.Printf("\n")
}

// add child to object
func (p *Portal) AddChild(obj pub_object.Object) {
	if obj != nil {
		p.scene.AddChild(obj.(pub_object.Object))
	}
}

// removes child from object
func (p *Portal) RemChild(obj pub_object.Object) {
	if obj != nil {
		p.scene.RemChild(obj.(pub_object.Object))
	}
}
