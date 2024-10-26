package assets

import (
	_ "embed"
)

var (
	//go:embed Slide_sprite.png
	SlidingPng []byte

	//go:embed Walk_sprite.png
	WalkingPng []byte

	//go:embed fonts/Lilita_One/LilitaOne-Regular.ttf
	MenuFontTtf []byte
)
