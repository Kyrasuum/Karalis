package stage

import (
	"fmt"

	"karalis/internal/cell"
	"karalis/internal/character"
	"karalis/internal/object/prim"
	"karalis/internal/object/world"
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

	ter, err := world.RandTerrain(0.0, 0.0, 256, 256, 1234567)
	if err != nil {
		return err
	}
	ter.SetScale(raylib.NewVector3(10, 5, 10))
	ter.SetPos(raylib.NewVector3(-5, 0, -5))
	g.curcell.AddChild(ter)

	return nil
}

// handle resize event
func (g *Game) OnResize(w int32, h int32) {
	if g == nil {
		return
	}

	g.player.GetCam().OnResize(w, h)
	g.curcell.OnResize(w, h)
}

// render hook
func (g *Game) Render() {
	if g == nil {
		return
	}
	cam := g.player.GetCam()
	if cam == nil {
		return
	}

	cmds := []func(){}

	cmds = append(cam.Prerender(), g.curcell.Prerender(cam)...)
	for _, cmd := range cmds {
		cmd()
	}

	cmds = cam.Render()
	cmds = append(g.curcell.Render(cam), cmds...)
	for _, cmd := range cmds {
		cmd()
	}

	cmds = append(cam.Postrender(), g.curcell.Postrender(cam)...)
	for _, cmd := range cmds {
		cmd()
	}

	return
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
