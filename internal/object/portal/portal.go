package portal

import (
	"image/color"
	"log"
	"runtime"
	"slices"

	"karalis/internal/object/prim"
	"karalis/internal/rlx"
	"karalis/internal/scene"
	"karalis/pkg/app"
	"karalis/pkg/lmath"

	pub_object "karalis/pkg/object"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Portal struct {
	parent pub_object.Object
	scene  *scene.Scene

	exit *Portal

	cleaner *runtime.Cleanup
	target  *rl.RenderTexture2D
	obj     pub_object.Object
	cam     pub_object.Camera

	touching []pub_object.Object

	rendering bool
	visible   bool
}

// Constuctor
func NewPortal(scene pub_object.Cell, exit *Portal, cam pub_object.Camera, obj pub_object.Object) (p *Portal, err error) {
	p = &Portal{}
	err = p.Init(scene, exit, cam, obj)

	return p, err
}

// initialize portal object
func (p *Portal) Init(scene pub_object.Cell, exit *Portal, cam pub_object.Camera, obj pub_object.Object) error {
	if p == nil {
		return nil
	}
	p.parent = nil

	if scene != nil {
		p.scene = scene
	} else {
		p.scene = &scene.Scene{}
		err := p.scene.Init()
		if err != nil {
			return err
		}
	}
	p.scene.OnAdd(p)

	if exit != nil {
		exit.Pair(p)
	}

	if cam != nil {
		p.cam = cam
	} else {
		p.cam = &pub_object.Camera{}
		err := p.cam.Init()
		if err != nil {
			return err
		}
	}

	text := rlx.LoadRenderTexture(app.CurApp.GetWidth(), app.CurApp.GetHeight())
	rlx.SetTextureFilter(text.Texture, rl.FilterBilinear)
	rlx.SetTextureWrap(text.Texture, rl.WrapRepeat)
	p.target = &text
	if p.cleaner != nil {
		p.cleaner.Stop()
	}
	cleaner := runtime.AddCleanup(p, func(text rl.RenderTexture2D) {
		rlx.UnloadRenderTexture(text)
	}, text)
	p.cleaner = &cleaner

	p.touching = []pub_object.Object{}

	p.rendering = false
	p.visible = true

	if obj != nil {
		p.SetPortal(obj)
	} else {
		sqr, err := prim.NewSquare()
		if err != nil {
			return err
		}
		p.SetPortal(sqr)
	}

	return nil
}

// get exit scene
func (p *Portal) GetScene() *pub_object.Object {
	if p == nil {
		return nil
	}

	return p.scene
}

// get exit portal pair
func (p *Portal) GetPair() pub_object.Portal {
	if p == nil {
		return nil
	}

	return p.exit
}

// set exit portal pair
func (p *Portal) Pair(e pub_object.Portal) {
	if p == nil {
		return
	}

	if e.exit != nil {
		e.exit = nil
	}
	if p.exit != nil {
		p.exit = nil
	}
	e.exit = p
	p.exit = e
}

// retrieve the portal display objects vertices
func (p *Portal) GetVertices() []rl.Vector3 {
	if p == nil {
		return []rl.Vector3{}
	}

	if p.obj != nil {
		return p.obj.GetVertices()
	} else {
		return []rl.Vector3{}
	}
}

func (p *Portal) SetColor(col color.Color) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetColor(col)
}

func (p *Portal) GetColor() color.Color {
	if p == nil || p.obj == nil {
		return nil
	}

	return p.obj.GetColor()
}

func (p *Portal) SetScale(sc rl.Vector3) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetScale(sc)
}

func (p *Portal) GetScale() rl.Vector3 {
	if p == nil || p.obj == nil {
		return rl.Vector3{}
	}

	return p.obj.GetScale()
}

func (p *Portal) SetPos(pos rl.Vector3) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetPos(pos)
}

func (p *Portal) GetPos() rl.Vector3 {
	if p == nil || p.obj == nil {
		return rl.Vector3{}
	}

	return p.obj.GetPos()
}

// retrieve the portal texture uvs for the display object
func (p *Portal) GetUVs() []rl.Vector2 {
	if p == nil {
		return []rl.Vector2{}
	}

	if p.obj != nil {
		return p.obj.GetUVs()
	} else {
		return []rl.Vector2{}
	}
}

// set the texture uvs for the portal display object
func (p *Portal) SetUVs(uvs []rl.Vector2) {
	if p == nil {
		return
	}

	if p.obj != nil {
		p.obj.SetUVs(uvs)
	}
}

// get the portal render objects model matrix
func (p *Portal) GetModelMatrix() rl.Matrix {
	if p == nil {
		return rl.Matrix{}
	}

	if p.obj != nil {
		return p.obj.GetModelMatrix()
	} else {
		return rl.MatrixTranslate(0, 0, 0)
	}
}

func (p *Portal) GetModel() *rl.Model {
	if p == nil || p.obj == nil {
		return nil
	}

	return p.obj.GetModel()
}

func (p *Portal) GetPitch() float32 {
	if p == nil || p.obj == nil {
		return 0
	}

	return p.obj.GetPitch()
}

func (p *Portal) SetPitch(pi float32) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetPitch(pi)
}

func (p *Portal) GetYaw() float32 {
	if p == nil || p.obj == nil {
		return 0
	}

	return p.obj.GetYaw()
}

func (p *Portal) SetYaw(y float32) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetYaw(y)
}

func (p *Portal) GetRoll() float32 {
	if p == nil || p.obj == nil {
		return 0
	}

	return p.obj.GetRoll()
}

func (p *Portal) SetRoll(r float32) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetRoll(r)
}

// get portal render material
func (p *Portal) GetMaterials() *rl.Material {
	if p == nil {
		return nil
	}

	if p.obj != nil {
		return p.obj.GetMaterials()
	}
	return &rl.Material{}
}

// set portal render texture
func (p *Portal) SetTexture(tex rl.Texture2D) {
	if p == nil {
		return
	}

	p.target.Texture = tex
}

// get portal render texture
func (p *Portal) GetTexture() *rl.Texture2D {
	if p == nil {
		return nil
	}

	return &p.target.Texture
}

// set portal render object
func (p *Portal) SetPortal(obj pub_object.Portal) {
	if p == nil {
		return
	}

	if p.obj != nil {
		p.obj.OnRemove()
	}
	p.obj = obj
	p.obj.OnAdd(p)
	p.obj.SetTexture(p.GetTexture())
	rlx.SetTextureFilter(p.target.Texture, rl.FilterBilinear)
	rlx.SetTextureWrap(p.target.Texture, rl.WrapRepeat)
}

// get portal render object
func (p *Portal) GetPortal() pub_object.Portal {
	if p == nil {
		return nil
	}

	return p.obj
}

// set camera for portal
func (p *Portal) SetCam(obj pub_object.Camera) {
	if p == nil {
		return
	}

	p.cam = obj
}

// return camera for portal
func (p *Portal) GetCam() pub_object.Camera {
	if p == nil {
		return nil
	}

	return p.cam
}

// return normal for portal plane
func (p *Portal) GetNormal() rl.Vector3 {
	if p == nil || p.obj == nil {
		return rl.Vector3{}
	}

	norm := rl.NewVector3(0, 0, 1)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(p.obj.GetPitch()), float64(p.obj.GetYaw()), float64(p.obj.GetRoll()))
	matRot := rl.QuaternionToMatrix(rl.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	norm = rl.Vector3Transform(norm, matRot)
	return norm
}

// calculate transform for exit portal to entry portal
func (p *Portal) GetTransform() rl.Matrix {
	if p == nil {
		return rl.MatrixIdentity()
	}

	//calculate portal camera position based on calling render camera
	wldToLcl := rl.MatrixInvert(p.GetModelMatrix())
	lclToWld := rl.MatrixIdentity()
	if p.exit != nil {
		//need to flip to place camera looking out not in
		flip := rl.MatrixRotateY(rl.Pi)
		lclToWld = rl.MatrixMultiply(flip, p.exit.GetModelMatrix())
	}
	return rl.MatrixMultiply(lclToWld, wldToLcl)
}

// prerender hook
func (p *Portal) Prerender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if p == nil {
		return cmds
	}

	//guards to ensure we only render when we should be
	if p.target != nil && p.scene != nil && !p.rendering && p.visible {
		//prevent rerendering a portal a second time
		p.rendering = true
		p.visible = false
		if p.exit != nil {
			p.exit.visible = false
		}

		mat := p.GetTransform()
		pos := cam.GetPos()
		tar := cam.GetTar()
		p.cam.SetPos(rl.Vector3Transform(pos, mat))
		p.cam.SetTar(rl.Vector3Transform(tar, mat))

		//render from portals perspective
		rlx.BeginTextureMode(*p.target)
		rlx.ClearBackground(color.RGBA{255, 255, 255, 255})
		sh := app.CurApp.GetShader()
		err := sh.SetDefine("PORTAL_SCN", true)
		if err != nil {
			log.Printf("%+v\n", err)
			p.visible = false
		}
		err = sh.SetUniform("portalPos", p.exit.obj.GetPos())
		if err != nil {
			log.Printf("%+v\n", err)
			p.visible = false
		}
		err = sh.SetUniform("portalNorm", p.exit.GetNormal())
		if err != nil {
			log.Printf("%+v\n", err)
			p.visible = false
		}

		cmds = p.cam.Prerender()
		cmds = append(cmds, p.scene.Prerender(p.cam)...)
		for _, cmd := range cmds {
			cmd()
		}

		cmds = p.cam.Render()

		err = sh.SetUniform("portalMat", mat)
		if err != nil {
			log.Printf("%+v\n", err)
			p.visible = false
		}
		for _, obj := range p.exit.touching {
			cmds = append(cmds, obj.Render(p.cam)...)
		}
		err = sh.SetUniform("portalMat", rl.MatrixIdentity())
		if err != nil {
			log.Printf("%+v\n", err)
			p.visible = false
		}

		cmds = append(cmds, p.scene.Render(p.cam)...)
		for _, cmd := range cmds {
			cmd()
		}

		cmds = p.cam.Postrender()
		cmds = append(cmds, p.scene.Postrender(p.cam)...)
		for _, cmd := range cmds {
			cmd()
		}

		err = sh.SetDefine("PORTAL_SCN", false)
		if err != nil {
			log.Printf("%+v\n", err)
			p.visible = false
		}
		rlx.EndTextureMode()

		p.rendering = false
		p.visible = true
		if p.exit != nil {
			p.exit.visible = true
		}
	}

	return cmds
}

// render hook
func (p *Portal) Render(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if p == nil {
		return cmds
	}

	//avoid rendering if not visible
	if p.visible {
		//render portal surface
		if p.target != nil && p.obj != nil {
			sh := app.CurApp.GetShader()
			err := sh.SetDefine("PORTAL_OBJ", true)
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}
			cmds = p.obj.Render(cam)
			err = sh.SetDefine("PORTAL_OBJ", false)
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}
		}

		col := p.exit.obj.GetCollider()
		if col != nil && p.obj != nil {
			//render objects exiting portal
			sh := app.CurApp.GetShader()
			err := sh.SetDefine("PORTAL_SCN", true)
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}
			err = sh.SetUniform("portalPos", p.obj.GetPos())
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}
			err = sh.SetUniform("portalNorm", p.GetNormal())
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}

			//prevent rerendering a portal a second time
			p.rendering = true
			if p.exit != nil {
				p.exit.visible = false
			}

			mat := rl.MatrixInvert(p.GetTransform())
			err = sh.SetUniform("portalMat", mat)
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}
			for _, obj := range p.touching {
				obj.Render(cam)
			}
			err = sh.SetUniform("portalMat", rl.MatrixIdentity())
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}

			p.rendering = false
			if p.exit != nil {
				p.exit.visible = true
			}

			err = sh.SetDefine("PORTAL_SCN", false)
			if err != nil {
				log.Printf("%+v\n", err)
				p.visible = false
			}
		}
	}

	return cmds
}

// postrender hook
func (p *Portal) Postrender(cam pub_object.Camera) []func() {
	cmds := []func(){}
	if p == nil {
		return cmds
	}

	return cmds
}

func (p *Portal) OnResize(w int32, h int32) {
	if p == nil {
		return
	}
}

// handle update cycle
func (p *Portal) Update(dt float32) {
	if p == nil {
		return
	}

	if p.obj != nil {
		p.obj.Update(dt)
	}

	var exit_scene *world.Cell
	if p.exit != nil {
		exit_scene = p.exit.scene
		p.exit.scene = nil
	}
	if p.scene != nil {
		p.scene.Update(dt)
	}
	if p.exit != nil {
		p.exit.scene = exit_scene
	}
	col := p.exit.obj.GetCollider()
	if col != nil {
		p.touching = col.GetTouching()
	}
	for _, obj := range p.touching {
		if rl.Vector3DotProduct(rl.Vector3Subtract(obj.GetPos(), p.GetPos()), p.GetNormal()) > 0 {
			if slices.Contains(p.GetScene().GetChilds(), obj) {
				p.GetScene().RemChild(obj)
				p.exit.GetScene().AddChild(obj)
				mat := rl.MatrixInvert(p.GetTransform())
				obj.SetPos(obj.GetPos().Transform(mat))
			}
		}
	}
}

// retrieve the collider for collision detection
func (p *Portal) GetCollider() pub_object.Collider {
	if p == nil || p.obj == nil {
		return nil
	}

	return p.obj.GetCollider()
}

// handle add event
func (p *Portal) OnAdd(obj pub_object.Object) {
	if p == nil {
		return
	}
	p.parent = obj
}

// handle remove event
func (p *Portal) OnRemove() {
	if p == nil {
		return
	}
	p.parent = nil
}

// handle resize event
func (p *Portal) OnResize(w int32, h int32) {
	if p == nil {
		return
	}

	if p.target != nil {
		rlx.UnloadRenderTexture(*p.target)
		p.target = nil
	}
	text := rlx.LoadRenderTexture(app.CurApp.GetWidth(), app.CurApp.GetHeight())
	rlx.SetTextureFilter(text.Texture, rl.FilterBilinear)
	rlx.SetTextureWrap(text.Texture, rl.WrapRepeat)
	p.target = &text
	if p.cleaner != nil {
		p.cleaner.Stop()
	}
	p.cleaner = &runtime.AddCleanup(p, func(text rl.RenderTexture2D) {
		rlx.UnloadRenderTexture(text)
	}, text)
}

// add child to object
func (p *Portal) AddChild(obj pub_object.Object) {
	if p == nil {
		return
	}

	if obj != nil {
		p.obj.AddChild(obj.(pub_object.Object))
	}
}

// removes child from object
func (p *Portal) RemChild(obj pub_object.Object) {
	if p == nil {
		return
	}

	if obj != nil {
		p.obj.RemChild(obj.(pub_object.Object))
	}
}

// gets all childs recursively
func (p *Portal) GetChilds() []pub_object.Object {
	if p == nil {
		return []pub_object.Object{}
	}

	if p.obj != nil {
		return p.obj.GetChilds()
	} else {
		return []pub_object.Object{}
	}
}

func (p *Portal) GetParent() pub_object.Object {
	if p == nil {
		return nil
	}
	return p.parent
}
