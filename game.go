package main

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type WorldDescriptor struct {
	key  string
	name string
}

type Game struct {
	Width, Height int
	Scene         Scene
	Debug         bool
	ShouldQuit    bool
	Worlds        []WorldDescriptor
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.Width, g.Height
}

func (g *Game) Update() error {
	var quit error

	err := g.Scene.Update()
	if err != nil {
		return err
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.Debug = !g.Debug
	}

	if g.ShouldQuit {
		quit = errors.New("quit")
	}

	return quit
}

func (g *Game) generateDebugString() string {
	screenInfo := fmt.Sprintf("%d FPS, Screen: %dx%d", int(ebiten.ActualFPS()), g.Width, g.Height)
	return fmt.Sprintf("%s", screenInfo)
}

func (g *Game) drawDebugText(screen *ebiten.Image, worldDebug string) {
	debugText := fmt.Sprintf("%s\n%s", g.generateDebugString(), worldDebug)
	ebitenutil.DebugPrint(screen, debugText)
}

func (g *Game) Draw(screen *ebiten.Image) {
	err := g.Scene.Draw(screen)
	if err != nil {
		return
	}

	if g.Debug {
		g.drawDebugText(screen, g.Scene.GenerateDebugString())
	}
}

func (g *Game) LoadWorld(world string) {
	g.Scene = NewWorld(g, world)
}

func (g *Game) LoadMenu() {
	menuItems := make([]*MenuItem, len(g.Worlds))
	for i, t := range g.Worlds {
		menuItems[i] = NewMenuItem(t.name, func() {
			g.LoadWorld(t.key)
		})
	}
	menuItems = append(menuItems, NewMenuItem("Quit", func() {
		g.ShouldQuit = true
	}))

	g.Scene = NewMenu(nil, menuItems...)
}

func NewGame() (*Game, error) {
	MenuInit()

	worlds := []WorldDescriptor{
		{
			key:  "Level1",
			name: "Tutorial",
		},
		{
			key:  "Level2",
			name: "Main World",
		},
	}

	g := &Game{
		Width:  ScreenWidth,
		Height: ScreenHeight,
		Worlds: worlds,
	}

	g.LoadMenu()

	return g, nil
}
