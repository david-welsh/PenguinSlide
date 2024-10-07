package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"image/color"
)

var (
	SnowColor = color.RGBA{
		R: 230,
		G: 230,
		B: 240,
		A: 255,
	}
	SkyColor = color.RGBA{
		R: 180,
		G: 180,
		B: 200,
		A: 255,
	}
)

type World struct {
	Game   *Game
	Space  *resolv.Space
	Player *Player
}

func NewWorld(game *Game) *World {
	w := &World{Game: game}
	w.Init()
	return w
}

func (world *World) Init() {
	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)

	world.Space = resolv.NewSpace(int(gw), int(gh), 16, 16)

	floor := resolv.NewObject(0, gh-100, gw, 100, "solid")
	world.Space.Add(floor)

	world.Player = NewPlayer(world.Space)
}

func (world *World) GenerateDebugString() string {
	worldDebug := ""
	playerDebug := world.Player.GenerateDebugText()
	return fmt.Sprintf("%s\n%s", worldDebug, playerDebug)
}

func (world *World) Update() {
	world.Player.Update()
}

func (world *World) Draw(screen *ebiten.Image) {
	screen.Fill(SkyColor)

	for _, o := range world.Space.Objects() {
		if o.HasTags("solid") {
			drawColor := SnowColor
			vector.DrawFilledRect(screen, float32(o.Position.X), float32(o.Position.Y), float32(o.Size.X), float32(o.Size.Y), drawColor, false)
		}
	}

	world.Player.Draw(screen)
}
