package stage

import (
	"fmt"

	"karalis/internal/cell"
	"karalis/internal/character"
	"karalis/internal/object"
	"karalis/internal/object/prim"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Game struct {
	curcell *cell.Cell
	player  *character.Player
}

// initialize game object
func (g *Game) Init() error {
	if g == nil {
		return fmt.Errorf("Invalid stage")
	}

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

	grid, err := prim.NewGrid()
	if err != nil {
		return err
	}
	g.curcell.AddChild(grid)

	// ter, err := object.NewTerrain("tex/map/heightmap.png", "tex/map/heightmap.png")
	ter, err := object.RandTerrain()
	if err != nil {
		return err
	}
	ter.SetScale(raylib.NewVector3(10, 5, 10))
	ter.SetPos(raylib.NewVector3(-5, -1.5, -5))
	g.curcell.AddChild(ter)

	return nil
}

// handle resize event
func (g *Game) OnResize(w int32, h int32) {
	if g == nil {
		return
	}

	g.player.GetCam().OnResize(w, h)
}

// prerender hook
func (g *Game) Prerender() []func() {
	if g == nil {
		return []func(){}
	}

	cam := g.player.GetCam()
	cmds := cam.Prerender()
	cmds = append(cmds, g.curcell.Prerender(cam)...)
	return cmds
}

// render hook
func (g *Game) Render() []func() {
	if g == nil {
		return []func(){}
	}

	cam := g.player.GetCam()
	cmds := cam.Render()
	cmds = append(cmds, g.curcell.Render(cam)...)
	return cmds
}

// postrender hook
func (g *Game) Postrender() []func() {
	if g == nil {
		return []func(){}
	}

	cam := g.player.GetCam()
	cmds := cam.Postrender()
	cmds = append(cmds, g.curcell.Postrender(cam)...)
	return cmds
}

// handle update cycle
func (g *Game) Update(dt float32) {
	if g == nil {
		return
	}
}

// handle player input
func (g *Game) OnInput(dt float32) {
	if g == nil {
		return
	}

	g.player.OnInput(dt)
}

// handle add event
func (g *Game) OnAdd() {
	if g == nil {
		return
	}

	g.curcell.OnAdd()
}

// handle remove event
func (g *Game) OnRemove() {
	if g == nil {
		return
	}

	g.curcell.OnRemove()
}

// get the player object
func (g *Game) GetPlayer() pub_object.Object {
	if g == nil {
		return nil
	}

	return g.player
}

// get currently active cell (where player is at)
func (g *Game) GetCurrentCell() pub_object.Cell {
	if g == nil {
		return nil
	}

	return g.curcell
}
