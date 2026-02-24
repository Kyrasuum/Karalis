package portal

import (
	"fmt"
	"image/color"
	"slices"

	"karalis/internal/camera"
	"karalis/internal/cell"
	"karalis/internal/object/prim"
	"karalis/pkg/app"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
	lmath "karalis/pkg/lmath"
)

var ()

type Portal struct {
	scene *cell.Cell

	exit *Portal

	target *raylib.RenderTexture2D
	obj    pub_object.Object
	cam    *camera.Cam

	touching []pub_object.Object

	rendering bool
	visible   bool
}

// Constuctor
func NewPortal(scene pub_object.Cell, exit *Portal, cam *camera.Cam, obj pub_object.Object) (p *Portal, err error) {
	p = &Portal{}
	err = p.Init(scene, exit, cam, obj)

	return p, err
}

// initialize portal object
func (p *Portal) Init(scene pub_object.Cell, exit *Portal, cam *camera.Cam, obj pub_object.Object) error {
	if p == nil {
		return nil
	}

	if scene != nil {
		p.scene = scene
	} else {
		p.scene = &cell.Cell{}
		err := p.scene.Init()
		if err != nil {
			return err
		}
	}

	if exit != nil {
		exit.Pair(p)
	}

	if cam != nil {
		p.cam = cam
	} else {
		p.cam = &camera.Cam{}
		err := p.cam.Init()
		if err != nil {
			return err
		}
	}

	text := raylib.LoadRenderTexture(app.CurApp.GetWidth(), app.CurApp.GetHeight())
	raylib.SetTextureFilter(text.Texture, raylib.FilterBilinear)
	raylib.SetTextureWrap(text.Texture, raylib.WrapRepeat)
	p.target = &text

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
func (p *Portal) GetScene() *cell.Cell {
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
func (p *Portal) GetVertices() []raylib.Vector3 {
	if p == nil {
		return []raylib.Vector3{}
	}

	if p.obj != nil {
		return p.obj.GetVertices()
	} else {
		return []raylib.Vector3{}
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

func (p *Portal) SetScale(sc raylib.Vector3) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetScale(sc)
}

func (p *Portal) GetScale() raylib.Vector3 {
	if p == nil || p.obj == nil {
		return raylib.Vector3{}
	}

	return p.obj.GetScale()
}

func (p *Portal) SetPos(pos raylib.Vector3) {
	if p == nil || p.obj == nil {
		return
	}

	p.obj.SetPos(pos)
}

func (p *Portal) GetPos() raylib.Vector3 {
	if p == nil || p.obj == nil {
		return raylib.Vector3{}
	}

	return p.obj.GetPos()
}

// retrieve the portal texture uvs for the display object
func (p *Portal) GetUVs() []raylib.Vector2 {
	if p == nil {
		return []raylib.Vector2{}
	}

	if p.obj != nil {
		return p.obj.GetUVs()
	} else {
		return []raylib.Vector2{}
	}
}

// set the texture uvs for the portal display object
func (p *Portal) SetUVs(uvs []raylib.Vector2) {
	if p == nil {
		return
	}

	if p.obj != nil {
		p.obj.SetUVs(uvs)
	}
}

// get the portal render objects model matrix
func (p *Portal) GetModelMatrix() raylib.Matrix {
	if p == nil {
		return raylib.Matrix{}
	}

	if p.obj != nil {
		return p.obj.GetModelMatrix()
	} else {
		return raylib.MatrixTranslate(0, 0, 0)
	}
}

func (p *Portal) GetModel() *raylib.Model {
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
func (p *Portal) GetMaterials() *raylib.Material {
	if p == nil {
		return nil
	}

	if p.obj != nil {
		return p.obj.GetMaterials()
	}
	return &raylib.Material{}
}

// set portal render texture
func (p *Portal) SetTexture(tex raylib.Texture2D) {
	if p == nil {
		return
	}

	p.target.Texture = tex
}

// get portal render texture
func (p *Portal) GetTexture() raylib.Texture2D {
	if p == nil {
		return raylib.Texture2D{}
	}

	return p.target.Texture
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
	p.obj.OnAdd()
	p.obj.SetTexture(p.GetTexture())
	raylib.SetTextureFilter(p.target.Texture, raylib.FilterBilinear)
	raylib.SetTextureWrap(p.target.Texture, raylib.WrapRepeat)
}

// get portal render object
func (p *Portal) GetPortal() pub_object.Portal {
	if p == nil {
		return nil
	}

	return p.obj
}

// set camera for portal
func (p *Portal) SetCam(obj *camera.Cam) {
	if p == nil {
		return
	}

	p.cam = obj
}

// return camera for portal
func (p *Portal) GetCam() *camera.Cam {
	if p == nil {
		return nil
	}

	return p.cam
}

// return normal for portal plane
func (p *Portal) GetNormal() raylib.Vector3 {
	if p == nil || p.obj == nil {
		return raylib.Vector3{}
	}

	norm := raylib.NewVector3(0, 0, 1)
	Quat := lmath.Quat{}
	Quat = *Quat.FromEuler(float64(p.obj.GetPitch()), float64(p.obj.GetYaw()), float64(p.obj.GetRoll()))
	matRot := raylib.QuaternionToMatrix(raylib.NewQuaternion(float32(Quat.X), float32(Quat.Y), float32(Quat.Z), float32(Quat.W)))
	norm = raylib.Vector3Transform(norm, matRot)
	return norm
}

// calculate transform for exit portal to entry portal
func (p *Portal) GetTransform() raylib.Matrix {
	if p == nil {
		return raylib.MatrixIdentity()
	}

	//calculate portal camera position based on calling render camera
	wldToLcl := raylib.MatrixInvert(p.GetModelMatrix())
	lclToWld := raylib.MatrixIdentity()
	if p.exit != nil {
		//need to flip to place camera looking out not in
		flip := raylib.MatrixRotateY(raylib.Pi)
		lclToWld = raylib.MatrixMultiply(flip, p.exit.GetModelMatrix())
	}
	return raylib.MatrixMultiply(lclToWld, wldToLcl)
}

// prerender hook
func (p *Portal) Prerender(cam *camera.Cam) []func() {
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
		p.cam.SetPos(raylib.Vector3Transform(pos, mat))
		p.cam.SetTar(raylib.Vector3Transform(tar, mat))

		//render from portals perspective
		raylib.BeginTextureMode(*p.target)
		raylib.ClearBackground(color.RGBA{255, 255, 255, 255})
		sh := app.CurApp.GetShader()
		err := sh.SetDefine("PORTAL_SCN", true)
		if err != nil {
			fmt.Printf("%+v\n", err)
			p.visible = false
		}
		err = sh.SetUniform("portalPos", p.exit.obj.GetPos())
		if err != nil {
			fmt.Printf("%+v\n", err)
			p.visible = false
		}
		err = sh.SetUniform("portalNorm", p.exit.GetNormal())
		if err != nil {
			fmt.Printf("%+v\n", err)
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
			fmt.Printf("%+v\n", err)
			p.visible = false
		}
		for _, obj := range p.exit.touching {
			cmds = append(cmds, obj.Render(p.cam)...)
		}
		err = sh.SetUniform("portalMat", raylib.MatrixIdentity())
		if err != nil {
			fmt.Printf("%+v\n", err)
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
			fmt.Printf("%+v\n", err)
			p.visible = false
		}
		raylib.EndTextureMode()

		p.rendering = false
		p.visible = true
		if p.exit != nil {
			p.exit.visible = true
		}
	}

	return cmds
}

// render hook
func (p *Portal) Render(cam *camera.Cam) []func() {
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
				fmt.Printf("%+v\n", err)
				p.visible = false
			}
			cmds = p.obj.Render(cam)
			err = sh.SetDefine("PORTAL_OBJ", false)
			if err != nil {
				fmt.Printf("%+v\n", err)
				p.visible = false
			}
		}

		col := p.exit.obj.GetCollider()
		if col != nil && p.obj != nil {
			//render objects exiting portal
			sh := app.CurApp.GetShader()
			err := sh.SetDefine("PORTAL_SCN", true)
			if err != nil {
				fmt.Printf("%+v\n", err)
				p.visible = false
			}
			err = sh.SetUniform("portalPos", p.obj.GetPos())
			if err != nil {
				fmt.Printf("%+v\n", err)
				p.visible = false
			}
			err = sh.SetUniform("portalNorm", p.GetNormal())
			if err != nil {
				fmt.Printf("%+v\n", err)
				p.visible = false
			}

			//prevent rerendering a portal a second time
			p.rendering = true
			if p.exit != nil {
				p.exit.visible = false
			}

			mat := raylib.MatrixInvert(p.GetTransform())
			err = sh.SetUniform("portalMat", mat)
			if err != nil {
				fmt.Printf("%+v\n", err)
				p.visible = false
			}
			for _, obj := range p.touching {
				obj.Render(cam)
			}
			err = sh.SetUniform("portalMat", raylib.MatrixIdentity())
			if err != nil {
				fmt.Printf("%+v\n", err)
				p.visible = false
			}

			p.rendering = false
			if p.exit != nil {
				p.exit.visible = true
			}

			err = sh.SetDefine("PORTAL_SCN", false)
			if err != nil {
				fmt.Printf("%+v\n", err)
				p.visible = false
			}
		}
	}

	return cmds
}

// postrender hook
func (p *Portal) Postrender(cam *camera.Cam) []func() {
	cmds := []func(){}
	if p == nil {
		return cmds
	}

	return cmds
}

// handle update cycle
func (p *Portal) Update(dt float32) {
	if p == nil {
		return
	}

	if p.obj != nil {
		p.obj.Update(dt)
	}

	var exit_scene *cell.Cell
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
		if raylib.Vector3DotProduct(raylib.Vector3Subtract(obj.GetPos(), p.GetPos()), p.GetNormal()) > 0 {
			if slices.Contains(p.GetScene().GetChilds(), obj) {
				p.GetScene().RemChild(obj)
				p.exit.GetScene().AddChild(obj)
				mat := raylib.MatrixInvert(p.GetTransform())
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
func (p *Portal) OnAdd() {
	if p == nil {
		return
	}

	if p.obj != nil {
		p.obj.OnAdd()
	}
}

// handle remove event
func (p *Portal) OnRemove() {
	if p == nil {
		return
	}

	if p.target != nil {
		raylib.UnloadRenderTexture(*p.target)
		p.target = nil
	}
	if p.obj != nil {
		p.obj.OnRemove()
	}
}

// handle resize event
func (p *Portal) OnResize(w int32, h int32) {
	if p == nil {
		return
	}

	p.cam.OnResize(w, h)
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
