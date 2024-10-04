package assets

import (
	_ "embed"
)

var (
	//go:embed Slide_sprite.png
	SlidingPng []byte

	//go:embed Walk_sprite.png
	WalkingPng []byte
)
