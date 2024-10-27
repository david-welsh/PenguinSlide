package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func LoadImage(img []byte) *ebiten.Image {
	wImg, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Fatal(err)
	}
	origImage := ebiten.NewImageFromImage(wImg)

	s := origImage.Bounds().Size()
	newImage := ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	newImage.DrawImage(origImage, op)
	return newImage
}
