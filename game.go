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
	Screen        *ebiten.Image
	World         *World
	Debug         bool
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
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

func generateDebugString() string {
	fps := fmt.Sprintf("%d FPS", int(ebiten.ActualFPS()))
	screenSize := fmt.Sprintf("Screen: %dx%d", ScreenWidth, ScreenHeight)
	return fmt.Sprintf("%s\n%s\n", fps, screenSize)
}

func drawDebugText(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, generateDebugString())
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
	if g.Debug {
		drawDebugText(screen)
	}
}

func NewGame() (*Game, error) {
	g := &Game{
		Width:  640,
		Height: 480,
	}

	g.World = NewWorld(g)

	return g, nil
}
