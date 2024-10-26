package main

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image) error
	GenerateDebugString() string
}
