package main

import (
	"fmt"
	"github.com/jakecoffman/cp/v2"
)

type Camera struct {
	Offset cp.Vector
	Zoom   float64
	Rotate float64
}

func (camera *Camera) GenerateDebugText() string {
	return fmt.Sprintf(
		"CAMERA: Offset: (%.1f, %.1f), Zoom: %.1f, Rotate: %.1f",
		camera.Offset.X,
		camera.Offset.Y,
		camera.Zoom,
		camera.Rotate,
	)
}
