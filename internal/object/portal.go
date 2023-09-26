package object

import (
	"karalis/internal/camera"
	"karalis/internal/cell"
	"karalis/pkg/app"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Portal struct {
	scene *cell.Cell

	exit *Portal

	target *raylib.RenderTexture2D
	obj    pub_object.Object
	cam    *camera.Cam

	rendering bool
	visible   bool
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
		p.scene.AddCam(p.cam)
	}

	if obj != nil {
		p.obj = obj
	} else {
		p.obj = nil
	}

	text := raylib.LoadRenderTexture(app.CurApp.GetWidth(), app.CurApp.GetHeight())
	p.target = &text

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

// get portal render material
func (p *Portal) GetMaterials() *raylib.Material {
	return &raylib.Material{}
}

// set portal render texture
func (p *Portal) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	p.target.Texture = tex
}

// get portal render texture
func (p *Portal) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return p.target.Texture
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
	if p.cam != nil {
		p.scene.RemCam(p.cam)
	}
	p.scene.AddCam(obj)
	p.cam = obj
}

// return camera for portal
func (p *Portal) GetCam() *camera.Cam {
	return p.cam
}

// prerender hook
func (p *Portal) Prerender() []func() {
	cmds := []func(){}

	if p.target != nil {
		if !p.rendering && p.visible {
			p.rendering = true
			raylib.BeginTextureMode(*p.target)
			raylib.ClearBackground(raylib.RayWhite)

			cmds = p.scene.Prerender()
			for _, cmd := range cmds {
				cmd()
			}

			cmds = p.scene.Render()
			for _, cmd := range cmds {
				cmd()
			}

			cmds = p.scene.Postrender()
			for _, cmd := range cmds {
				cmd()
			}

			raylib.EndTextureMode()
			p.rendering = false
		}
	}

	return cmds
}

// render hook
func (p *Portal) Render() []func() {
	cmds := []func(){}

	if p.visible {
		if p.target != nil && p.obj != nil {
			cmds = p.obj.Render()
		}
	}

	return cmds
}

// postrender hook
func (p *Portal) Postrender() []func() {
	cmds := []func(){}

	if p.visible {
		if p.target != nil {
			if p.obj == nil {
				raylib.DrawTexture(p.target.Texture, 0, 0, raylib.White)
			}
		}
	}

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
	if p.target != nil {
		raylib.UnloadRenderTexture(*p.target)
		p.target = nil
	}
	p.scene.OnRemove()
}

// handle resize event
func (p *Portal) OnResize(w int32, h int32) {
	p.scene.OnResize(w, h)
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
