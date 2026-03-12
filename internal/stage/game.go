package stage

import (
	"fmt"

	"karalis/internal/object/character"
	"karalis/internal/object/prim"
	"karalis/internal/object/world"
	"karalis/internal/scene"
	pub_object "karalis/pkg/object"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Game struct {
	scene  *scene.Scene
	player *character.Player
	rt     raylib.RenderTexture2D
	resize bool
}

// initialize game object
func (g *Game) Init() error {
	if g == nil {
		return fmt.Errorf("Invalid stage")
	}
	g.rt = raylib.LoadRenderTexture(int32(raylib.GetRenderWidth()), int32(raylib.GetRenderHeight()))

	g.scene = &scene.Scene{}
	err := g.scene.Init()
	if err != nil {
		return err
	}
	g.scene.OnAdd(nil)

	g.player, err = character.NewPlayer()
	if err != nil {
		return err
	}
	g.scene.AddChild(g.player)

	grid, err := prim.NewGrid()
	if err != nil {
		return err
	}
	g.scene.AddChild(grid)

	w, err := world.NewTerrainWorld()
	if err != nil {
		return err
	}
	g.scene.AddChild(w)

	return nil
}

// handle resize event
func (g *Game) OnResize(w int32, h int32) {
	if g == nil {
		return
	}

	g.resize = true
	if g.player != nil {
		g.player.GetCam().OnResize(w, h)
	}
	if g.scene != nil {
		g.scene.OnResize(w, h)
	}
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
	if g.resize {
		raylib.UnloadRenderTexture(g.rt)
		g.rt = raylib.LoadRenderTexture(int32(raylib.GetRenderWidth()), int32(raylib.GetRenderHeight()))
		g.resize = false
	}
	cam.UpdateCam()

	cmds := []func(){}

	cmds = append(cam.Prerender(), g.scene.Prerender(cam)...)
	for _, cmd := range cmds {
		cmd()
	}

	raylib.BeginTextureMode(g.rt)
	raylib.ClearBackground(raylib.Black)
	cmds = cam.Render()
	cmds = append(g.scene.Render(cam), cmds...)
	for _, cmd := range cmds {
		cmd()
	}
	raylib.EndTextureMode()

	cmds = append(cam.Postrender(), g.scene.Postrender(cam)...)
	raylib.BeginTextureMode(g.rt)
	for _, cmd := range cmds {
		cmd()
	}
	raylib.EndTextureMode()
	raylib.DrawTexturePro(g.rt.Texture, raylib.Rectangle{0, 0, float32(g.rt.Texture.Width), -float32(g.rt.Texture.Height)}, raylib.Rectangle{0, 0, float32(g.rt.Texture.Width), float32(g.rt.Texture.Height)}, raylib.Vector2{0, 0}, 0.0, raylib.RayWhite)

	return
}

// handle update cycle
func (g *Game) Update(dt float32) {
	if g == nil || g.scene == nil {
		return
	}

	g.scene.Update(dt)
}

// handle player input
func (g *Game) OnInput(dt float32) {
	if g == nil || g.scene == nil {
		return
	}

	g.player.OnInput(dt)
}

// handle add event
func (g *Game) OnAdd() {
	if g == nil || g.scene == nil {
		return
	}
}

// handle remove event
func (g *Game) OnRemove() {
	if g == nil || g.scene == nil {
		return
	}
}

// get the player object
func (g *Game) GetPlayer() pub_object.Object {
	if g == nil || g.player == nil {
		return nil
	}

	return g.player
}

// get currently active scene (where player is at)
func (g *Game) GetCurrentScene() pub_object.Object {
	if g == nil {
		return nil
	}

	return g.scene
}
