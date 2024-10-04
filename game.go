package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	g.World.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
}

func NewGame() (*Game, error) {
	g := &Game{
		Width:  640,
		Height: 480,
	}

	g.World = NewWorld(g)

	return g, nil
}
