package world

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"slices"

	"karalis/pkg/app"
	"karalis/pkg/lmath"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var (
	CellScale  = raylib.Vector3{10, 10, 10}
	CellOffset = raylib.Vector3{-5, 0, -5}
	CellRender = 1
	Seed       = int64(1234567)
)

type World struct {
	parent  pub_object.Object
	childs  []pub_object.Object
	cells   map[string]*Cell
	sky     *Skybox
	pos     raylib.Vector2
	creator func(raylib.Vector3, raylib.Vector3, int64) (*Cell, error)
}

func NewTerrainWorld() (*World, error) {
	w := &World{}
	err := w.Init()
	if err != nil {
		return nil, err
	}
	w.creator = NewTerrainCell

	return w, nil
}

func NewCityWorld() (*World, error) {
	w := &World{}
	err := w.Init()
	if err != nil {
		return nil, err
	}
	w.creator = NewCityCell

	return w, nil
}

func NewDungeonWorld() (*World, error) {
	w := &World{}
	err := w.Init()
	if err != nil {
		return nil, err
	}
	w.creator = NewDungeonCell

	return w, nil
}

func (w *World) Init() error {
	if w == nil {
		return errors.New("Invalid world")
	}
	w.parent = nil
	w.childs = []pub_object.Object{}
	w.cells = map[string]*Cell{}

	sky, err := NewSkybox(nil)
	if err != nil {
		return err
	}
	w.sky = sky
	sky.OnAdd(w)

	return nil
}

func (w *World) Prerender(cam pub_object.Camera) []func() {
	if w == nil {
		return []func(){}
	}
	w.GenCells()

	cmds := []func(){}
	if w.sky != nil {
		cmds = append(cmds, w.sky.Prerender(cam)...)
	}
	for _, cell := range w.cells {
		cmds = append(cmds, cell.Prerender(cam)...)
	}
	for _, child := range w.childs {
		cmds = append(cmds, child.Prerender(cam)...)
	}
	return cmds
}

func (w *World) Render(cam pub_object.Camera) []func() {
	if w == nil {
		return []func(){}
	}

	cmds := []func(){}
	if w.sky != nil {
		cmds = append(cmds, w.sky.Render(cam)...)
	}
	for _, cell := range w.cells {
		cmds = append(cmds, cell.Render(cam)...)
	}
	for _, child := range w.childs {
		cmds = append(cmds, child.Render(cam)...)
	}
	return cmds
}

func (w *World) Postrender(cam pub_object.Camera) []func() {
	if w == nil {
		return []func(){}
	}

	cmds := []func(){}
	if w.sky != nil {
		cmds = append(cmds, w.sky.Postrender(cam)...)
	}
	for _, cell := range w.cells {
		cmds = append(cmds, cell.Postrender(cam)...)
	}
	for _, child := range w.childs {
		cmds = append(cmds, child.Postrender(cam)...)
	}
	cmds = append(cmds, DrawUnderwater(w)...)

	return cmds
}

func (w *World) OnResize(wi int32, he int32) {
	if w == nil {
		return
	}

	//perform update on objects
	if w.sky != nil {
		w.sky.OnResize(wi, he)
	}
	for _, child := range w.childs {
		child.OnResize(wi, he)
	}
	for _, cell := range w.cells {
		cell.OnResize(wi, he)
	}
}

func (w *World) Update(dt float32) {
	if w == nil {
		return
	}

	//if player is in world we generate cells around player
	ply := app.CurApp.GetStage().GetPlayer()
	if w.GetParent() == ply.GetParent() {
		ply := app.CurApp.GetStage().GetPlayer()
		pos := ply.GetPos()
		w.pos = raylib.Vector2{pos.X, pos.Z}
	}

	//perform update on objects
	if w.sky != nil {
		w.sky.Update(dt)
	}
	for _, child := range w.childs {
		child.Update(dt)
	}
	for _, cell := range w.cells {
		cell.Update(dt)
	}
}

func (w *World) GenCells() {
	if w == nil {
		return
	}

	ply := app.CurApp.GetStage().GetPlayer()
	plypos := raylib.Vector2{float32(lmath.Round(ply.GetPos().X / CellScale.X)), float32(lmath.Round(ply.GetPos().Z / CellScale.Z))}

	//drop off cells
	remove := []string{}
	for _, cell := range w.cells {
		cpos := raylib.Vector2{float32(lmath.Round((cell.GetPos().X + CellScale.X/2) / CellScale.X)), float32(lmath.Round((cell.GetPos().Z + CellScale.Z/2) / CellScale.Z))}
		diff := raylib.Vector2{lmath.Abs(plypos.X - cpos.X), lmath.Abs(plypos.Y - cpos.Y)}
		if diff.X > float32(CellRender) || diff.Y > float32(CellRender) {
			spos := fmt.Sprintf("%d %d", cpos.X, cpos.Y)
			remove = append(remove, spos)
		}
	}
	for _, spos := range remove {
		delete(w.cells, spos)
	}

	//add cells as needed
	for x := range CellRender*2 + 1 {
		for y := range CellRender*2 + 1 {
			cpos := raylib.Vector2{float32(x-CellRender) + plypos.X, float32(y-CellRender) + plypos.Y}
			spos := fmt.Sprintf("%d %d", cpos.X, cpos.Y)
			if _, ok := w.cells[spos]; !ok {
				offset := raylib.Vector3{cpos.X*CellScale.X + CellOffset.X, CellOffset.Y, cpos.Y*CellScale.Z + CellOffset.Z}
				cell, err := w.creator(offset, CellScale, Seed)
				if err != nil {
					log.Printf("%+v\n", err)
				}
				w.cells[spos] = cell
				cell.OnAdd(w)
			}
		}
	}
}

func (w *World) OnAdd(obj pub_object.Object) {
	if w == nil {
		return
	}
	w.parent = obj
}

func (w *World) OnRemove() {
	if w == nil {
		return
	}
	w.parent = nil
}

func (w *World) AddChild(obj pub_object.Object) {
	if w == nil {
		return
	}

	w.childs = append(w.childs, obj)
	obj.OnAdd(w)
}

func (w *World) RemChild(obj pub_object.Object) {
	if w == nil {
		return
	}
	//find index of child
	index := -1
	for i, child := range w.childs {
		if obj == child {
			index = i
			break
		}
	}

	//remove child
	if index >= 0 {
		w.childs[index] = w.childs[len(w.childs)-1]
		w.childs = w.childs[:len(w.childs)-1]
		obj.OnRemove()
	}
}

func (w *World) GetChilds() []pub_object.Object {
	if w == nil {
		return []pub_object.Object{}
	}

	childs := append(w.childs, w.sky)
	grandchilds := []pub_object.Object{}
	for _, cell := range w.cells {
		childs = append(w.childs, cell)
	}
	for _, child := range childs {
		grandchilds = append(grandchilds, child.GetChilds()...)
	}

	return slices.Concat(grandchilds, childs)
}

func (w *World) GetCollider() pub_object.Collider {
	if w == nil {
		return nil
	}

	return nil
}

func (w *World) GetModelMatrix() raylib.Matrix {
	if w == nil {
		return raylib.Matrix{}
	}
	return raylib.Matrix{}
}

func (w *World) GetModel() *raylib.Model {
	if w == nil {
		return nil
	}

	return nil
}

func (w *World) SetColor(col color.Color) {
	if w == nil {
		return
	}
}

func (w *World) GetColor() color.Color {
	if w == nil {
		return nil
	}

	return nil
}

func (w *World) GetScale() raylib.Vector3 {
	if w == nil {
		return raylib.Vector3{}
	}

	return raylib.Vector3{}
}

func (w *World) SetScale(sc raylib.Vector3) {
	if w == nil {
		return
	}
}

func (w *World) SetPos(pos raylib.Vector3) {
	if w == nil {
		return
	}
	w.pos = raylib.Vector2{pos.X, pos.Z}
}

func (w *World) GetPos() raylib.Vector3 {
	if w == nil {
		return raylib.Vector3{}
	}

	return raylib.Vector3{w.pos.X, 0, w.pos.Y}
}

func (w *World) GetPitch() float32 {
	if w == nil {
		return 0
	}

	return 0
}

func (w *World) SetPitch(pitch float32) {
	if w == nil {
		return
	}
}

func (w *World) GetYaw() float32 {
	if w == nil {
		return 0
	}

	return 0
}

func (w *World) SetYaw(yaw float32) {
	if w == nil {
		return
	}
}

func (w *World) GetRoll() float32 {
	if w == nil {
		return 0
	}

	return 0
}

func (w *World) SetRoll(roll float32) {
	if w == nil {
		return
	}
}

func (w *World) GetVertices() []raylib.Vector3 {
	if w == nil {
		return []raylib.Vector3{}
	}
	return []raylib.Vector3{}
}

func (w *World) GetUVs() []raylib.Vector2 {
	if w == nil {
		return []raylib.Vector2{}
	}
	return []raylib.Vector2{}
}

func (w *World) SetUVs(uvs []raylib.Vector2) {
	if w == nil {
		return
	}
}

func (w *World) GetMaterials() *raylib.Material {
	if w == nil {
		return nil
	}

	return nil
}

func (w *World) SetTexture(tex raylib.Texture2D) {
	if w == nil {
		return
	}
}

func (w *World) GetTexture() *raylib.Texture2D {
	if w == nil {
		return nil
	}

	return nil
}

func (w *World) GetParent() pub_object.Object {
	if w == nil {
		return nil
	}
	return w.parent
}
