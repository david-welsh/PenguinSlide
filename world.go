package main

import (
	"fmt"
	"github.com/demouth/ebitencp"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp/v2"
	"image/color"
)

var (
	SkyColor = color.RGBA{
		R: 180,
		G: 180,
		B: 200,
		A: 255,
	}
)

type World struct {
	Game   *Game
	Space  *cp.Space
	Player *Player
	Drawer *ebitencp.Drawer
	Camera Camera
}

func NewWorld(game *Game) *World {
	w := &World{Game: game}
	w.Init()
	return w
}

func (world *World) Init() {
	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)

	world.Camera = Camera{
		Zoom:   1,
		Offset: cp.Vector{X: -ScreenWidth / 2, Y: -ScreenHeight / 2},
	}

	world.Drawer = ebitencp.NewDrawer(0, 0)
	world.Drawer.FlipYAxis = true
	world.Drawer.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	world.Space = cp.NewSpace()

	world.Space.SleepTimeThreshold = 0.5
	world.Space.SetCollisionSlop(0.5)

	fs1 := cp.NewSegment(world.Space.StaticBody, cp.Vector{X: 0, Y: gh - 110}, cp.Vector{X: gw / 2, Y: gh - 110}, 0)
	fs1.SetFriction(0.7)
	world.Space.AddShape(fs1)

	fs2 := cp.NewSegment(world.Space.StaticBody, cp.Vector{X: gw / 2, Y: gh - 110}, cp.Vector{X: gw, Y: gh - 80}, 0)
	fs2.SetFriction(0.2)
	world.Space.AddShape(fs2)

	fs3 := cp.NewSegment(world.Space.StaticBody, cp.Vector{X: gw, Y: gh - 80}, cp.Vector{X: gw * 2, Y: gh - 80}, 0)
	fs3.SetFriction(0.2)
	world.Space.AddShape(fs3)

	world.Space.Iterations = 10
	world.Space.SetGravity(cp.Vector{Y: 400})
	world.Space.SetCollisionSlop(0.5)

	world.Player = NewPlayer(world.Space, world.Game)
}

func (world *World) GenerateDebugString() string {
	worldDebug := ""
	cameraDebug := world.Camera.GenerateDebugText()
	playerDebug := world.Player.GenerateDebugText()
	return fmt.Sprintf("%s\n%s\n%s", worldDebug, cameraDebug, playerDebug)
}

func (world *World) Update() {
	world.Player.Update()

	world.Space.Step(1.0 / float64(ebiten.TPS()))

	world.Camera.Offset.X = -world.Player.Body.Position().X

	world.Drawer.GeoM.Reset()
	world.Drawer.GeoM.Translate(world.Camera.Offset.X, world.Camera.Offset.Y)
	world.Drawer.GeoM.Scale(world.Camera.Zoom, world.Camera.Zoom)
	world.Drawer.GeoM.Rotate(world.Camera.Rotate)
	world.Drawer.GeoM.Translate(float64(world.Game.Width)/2, float64(world.Game.Height)/2)
}

func (world *World) Draw(screen *ebiten.Image) {
	screen.Fill(SkyColor)

	cp.DrawSpace(
		world.Space,
		world.Drawer.WithScreen(screen),
	)

	world.Player.Draw(screen, *world.Drawer.GeoM)
}
