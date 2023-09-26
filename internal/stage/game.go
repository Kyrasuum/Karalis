package stage

import (
	"godev/internal/cell"
	"godev/internal/character"
	"godev/internal/object"
	"godev/internal/object/prim"
	pub_object "godev/pkg/object"
)

var ()

type Game struct {
	curcell *cell.Cell
	player  *character.Player
	portal1 object.Portal
	portal2 object.Portal
}

// initialize game object
func (g *Game) Init() {
	g.curcell = &cell.Cell{}
	g.curcell.Init()

	g.player = &character.Player{}
	g.player.Init()
	pcam := g.player.GetCam()
	g.curcell.AddCam(pcam)
	g.curcell.AddChild(g.player)

	grid2 := prim.Grid{}
	grid2.Init()
	g.curcell.AddChild(&grid2)

	g.portal1 = object.Portal{}
	g.portal1.Init(nil, nil, nil, nil)
	g.curcell.AddChild(&g.portal1)

	grid1 := prim.Grid{}
	grid1.Init()
	g.portal1.AddChild(&grid1)

	box1 := prim.Cube{}
	box1.Init()
	g.portal1.AddChild(&box1)

	plane1 := prim.Square{}
	plane1.Init()
	g.portal1.SetPortal(&plane1)
}

// handle resize event
func (g *Game) OnResize(w int32, h int32) {
	g.curcell.OnResize(w, h)
	g.portal1.OnResize(w, h)
	g.portal2.OnResize(w, h)
}

// prerender hook
func (g *Game) Prerender() []func() {
	cmds := g.curcell.Prerender()
	return cmds
}

// render hook
func (g *Game) Render() []func() {
	cmds := g.curcell.Render()
	return cmds
}

// postrender hook
func (g *Game) Postrender() []func() {
	cmds := g.curcell.Postrender()
	return cmds
}

// handle update cycle
func (g *Game) Update(dt float32) {
	g.curcell.Update(dt)
}

// handle player input
func (g *Game) OnInput(dt float32) {
	g.player.OnInput(dt)
}

// handle add event
func (g *Game) OnAdd() {
	g.curcell.OnAdd()
}

// handle remove event
func (g *Game) OnRemove() {
	g.curcell.OnRemove()
}

// get the player object
func (g *Game) GetPlayer() pub_object.Object {
	return g.player
}

// get currently active cell (where player is at)
func (g *Game) GetCurrentCell() *cell.Cell {
	return g.curcell
}
