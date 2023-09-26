package stage

import (
	"godev/internal/camera"
	"godev/internal/cell"
	"godev/internal/character"
	"godev/internal/object"
	"godev/internal/object/prim"
)

var ()

type Game struct {
	curcell cell.Cell
	player  character.Player
	portal  object.Portal
}

func (g *Game) Init() {
	g.curcell = cell.Cell{}
	g.curcell.Init()

	g.player = character.Player{}
	g.player.Init()
	pcam := g.player.GetCam()
	g.curcell.AddCam(pcam)
	g.curcell.AddChild(&g.player)

	grid2 := prim.Grid{}
	grid2.Init()
	g.curcell.AddChild(&grid2)

	g.portal = object.Portal{}
	g.portal.Init()
	g.curcell.AddChild(&g.portal)

	cam := camera.Cam{}
	cam.Init()
	g.portal.SetCam(&cam)

	grid1 := prim.Grid{}
	grid1.Init()
	g.portal.AddChild(&grid1)

	box1 := prim.Cube{}
	box1.Init()
	g.portal.AddChild(&box1)

	plane1 := prim.Square{}
	plane1.Init()
	g.portal.SetPortal(&plane1)
}

func (g *Game) OnResize(w int32, h int32) {
	g.curcell.OnResize(w, h)
	g.portal.OnResize(w, h)
}

func (g *Game) Prerender() []func() {
	cmds := g.curcell.Prerender()
	return cmds
}

func (g *Game) Render() []func() {
	cmds := g.curcell.Render()
	return cmds
}

func (g *Game) Postrender() []func() {
	cmds := g.curcell.Postrender()
	return cmds
}

func (g *Game) Update(dt float32) {
	g.curcell.Update(dt)
}

func (g *Game) OnInput(dt float32) {
	g.player.OnInput(dt)
}

func (g *Game) OnAdd() {
	g.curcell.OnAdd()
}

func (g *Game) OnRemove() {
	g.curcell.OnRemove()
}
