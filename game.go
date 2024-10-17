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

type Game struct {
	Width, Height int
	World         *World
	Debug         bool
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.Width, g.Height
}

func (g *Game) Update() error {
	var quit error

	g.World.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.Debug = !g.Debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
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
	g.World.Draw(screen)

	if g.Debug {
		g.drawDebugText(screen, g.World.GenerateDebugString())
	}
}

func NewGame() (*Game, error) {
	g := &Game{
		Width:  ScreenWidth,
		Height: ScreenHeight,
	}

	g.World = NewWorld(g, "Level1")

	return g, nil
}
