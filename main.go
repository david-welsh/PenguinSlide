package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(game.Width, game.Height)
	ebiten.SetWindowTitle("Penguin Slide")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
