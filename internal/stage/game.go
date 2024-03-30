package stage

import (
	"karalis/internal/cell"
	"karalis/internal/character"
	"karalis/internal/object"
	"karalis/internal/object/prim"
	"karalis/internal/portal"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Game struct {
	curcell *cell.Cell
	player  *character.Player
	portal1 *portal.Portal
}

// initialize game object
func (g *Game) Init() error {
	g.curcell = &cell.Cell{}
	err := g.curcell.Init()
	if err != nil {
		return err
	}

	g.player, err = character.NewPlayer()
	if err != nil {
		return err
	}
	g.curcell.AddChild(g.player)

	grid2, err := prim.NewGrid()
	if err != nil {
		return err
	}
	g.curcell.AddChild(grid2)

	g.portal1, err = portal.NewPortal(nil, nil, nil, nil)
	if err != nil {
		return err
	}
	g.curcell.AddChild(g.portal1)

	grid1, err := prim.NewGrid()
	if err != nil {
		return err
	}
	g.portal1.AddChild(grid1)

	box1, err := prim.NewCube()
	if err != nil {
		return err
	}
	box1.SetPos(raylib.NewVector3(0, 0, -2))
	g.portal1.AddChild(box1)

	return nil
}

// handle resize event
func (g *Game) OnResize(w int32, h int32) {
	g.player.GetCam().OnResize(w, h)
	g.portal1.OnResize(w, h)
}

// prerender hook
func (g *Game) Prerender() []func() {
	cam := g.player.GetCam()
	cmds := cam.Prerender()
	cmds = append(cmds, g.curcell.Prerender(cam)...)
	return cmds
}

// render hook
func (g *Game) Render() []func() {
	cam := g.player.GetCam()
	cmds := cam.Render()
	cmds = append(cmds, g.curcell.Render(cam)...)
	return cmds
}

// postrender hook
func (g *Game) Postrender() []func() {
	cam := g.player.GetCam()
	cmds := cam.Postrender()
	cmds = append(cmds, g.curcell.Postrender(cam)...)
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
