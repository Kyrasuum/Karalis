package object

import (
	"godev/internal/camera"
	"godev/internal/cell"
	"godev/pkg/app"
	pub_object "godev/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Portal struct {
	target *raylib.RenderTexture2D
	scene  cell.Cell
	obj    pub_object.Object
	cam    *camera.Cam
}

func (p *Portal) Init() {
	p.scene = cell.Cell{}
	p.scene.Init()
	p.obj = nil
	p.target = nil
}

func (p *Portal) GetMaterials() *raylib.Material {
	return &raylib.Material{}
}

func (p *Portal) SetTexture(mat *raylib.Material, tex raylib.Texture2D) {
	p.target.Texture = tex
}

func (p *Portal) GetTexture(mat *raylib.Material) raylib.Texture2D {
	return p.target.Texture
}

func (p *Portal) SetPortal(obj pub_object.Object) {
	if p.obj != nil {
		p.obj.OnRemove()
	}
	p.obj = obj
	p.obj.OnAdd()
	p.obj.SetTexture(p.obj.GetMaterials(), p.GetTexture(nil))
}

func (p *Portal) GetPortal() pub_object.Object {
	return p.obj
}

func (p *Portal) SetCam(obj *camera.Cam) {
	if p.cam != nil {
		p.scene.RemCam(p.cam)
	}
	p.scene.AddCam(obj)
	p.cam = obj
}

func (p *Portal) GetCam() *camera.Cam {
	return p.cam
}

func (p *Portal) Prerender() []func() {
	cmds := []func(){}
	if p.target != nil {
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
	}
	return cmds
}

func (p *Portal) Render() []func() {
	cmds := []func(){}
	if p.target != nil && p.obj != nil {
		cmds = p.obj.Render()
	}
	return cmds
}

func (p *Portal) Postrender() []func() {
	cmds := []func(){}
	if p.target != nil {
		if p.obj == nil {
			raylib.DrawTexture(p.target.Texture, 0, 0, raylib.White)
		}
	}
	return cmds
}

func (p *Portal) Update(dt float32) {
	if p.obj != nil {
		p.obj.Update(dt)
	}
	p.scene.Update(dt)
}

func (p *Portal) OnAdd() {
	if p.target == nil {
		text := raylib.LoadRenderTexture(app.CurApp.GetWidth(), app.CurApp.GetHeight())
		p.target = &text
	}
	p.scene.OnAdd()
}

func (p *Portal) OnRemove() {
	if p.target != nil {
		raylib.UnloadRenderTexture(*p.target)
		p.target = nil
	}
	p.scene.OnRemove()
}

func (p *Portal) OnResize(w int32, h int32) {
	p.scene.OnResize(w, h)
}

func (p *Portal) AddChild(obj pub_object.Object) {
	if obj != nil {
		p.scene.AddChild(obj.(pub_object.Object))
	}
}

func (p *Portal) RemChild(obj pub_object.Object) {
	if obj != nil {
		p.scene.RemChild(obj.(pub_object.Object))
	}
}
