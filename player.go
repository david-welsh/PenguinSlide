package main

import (
	"PenguinSlide/assets"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"image"
	"image/color"
	"log"
	"math"
)

type Player struct {
	Object      *resolv.Object
	Speed       resolv.Vector
	OnGround    *resolv.Object
	FacingRight bool
}

var (
	walkingImage *ebiten.Image
	slidingImage *ebiten.Image
)

func init() {
	wImg, _, err := image.Decode(bytes.NewReader(assets.WalkingPng))
	if err != nil {
		log.Fatal(err)
	}
	origWalkingImage := ebiten.NewImageFromImage(wImg)

	s := origWalkingImage.Bounds().Size()
	walkingImage = ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	walkingImage.DrawImage(origWalkingImage, op)

	sImg, _, sErr := image.Decode(bytes.NewReader(assets.SlidingPng))
	if sErr != nil {
		log.Fatal(sErr)
	}
	origSlidingImage := ebiten.NewImageFromImage(sImg)

	sz := origSlidingImage.Bounds().Size()
	slidingImage = ebiten.NewImage(sz.X, sz.Y)

	op2 := &ebiten.DrawImageOptions{}
	slidingImage.DrawImage(origSlidingImage, op2)
}

func (p *Player) Update() {
	friction := 0.05
	accel := 0.5 + friction
	maxSpeed := 4.0
	gravity := 0.75

	p.Speed.Y += gravity

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxisValue(0, 0) > 0.1 {
		p.Speed.X += accel
		p.FacingRight = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxisValue(0, 0) < -0.1 {
		p.Speed.X -= accel
		p.FacingRight = false
	}

	if p.Speed.X > friction {
		p.Speed.X -= friction
	} else if p.Speed.X < -friction {
		p.Speed.X += friction
	} else {
		p.Speed.X = 0
	}

	if p.Speed.X > maxSpeed {
		p.Speed.X = maxSpeed
	} else if p.Speed.X < -maxSpeed {
		p.Speed.X = -maxSpeed
	}

	dx := p.Speed.X

	p.Object.Position.X += dx

	if p.Speed.X > 0 || p.Speed.Y < 0 {
		p.Object.Shape.Rotate(math.Pi / 2)
	}

	p.OnGround = nil

	dy := p.Speed.Y

	dy = math.Max(math.Min(dy, 16), -16)

	checkDistance := dy
	if dy >= 0 {
		checkDistance++
	}

	if check := p.Object.Check(0, checkDistance, "solid", "platform", "ramp"); check != nil {
		slide, slideOK := check.SlideAgainstCell(check.Cells[0], "solid")

		if dy < 0 && check.Cells[0].ContainsTags("solid") && slideOK && math.Abs(slide.X) <= 8 {
			p.Object.Position.X += slide.X
		} else {
			if solids := check.ObjectsByTags("solid"); len(solids) > 0 && (p.OnGround == nil || p.OnGround.Position.Y >= solids[0].Position.Y) {
				dy = check.ContactWithObject(solids[0]).Y
				p.Speed.Y = 0

				if solids[0].Position.Y > p.Object.Position.Y {
					p.OnGround = solids[0]
				}

			}
		}
	}

	p.Object.Position.Y += dy

	p.Object.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {
	img := walkingImage
	if p.Speed.X > 0 || p.Speed.X < 0 {
		img = slidingImage
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.1, 0.1)
	if !p.FacingRight {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(p.Object.Size.X, 0)
	}
	op.GeoM.Translate(float64(p.Object.Position.X), float64(p.Object.Position.Y))
	screen.DrawImage(img, op)

	cl := color.RGBA{R: 255, A: 120}
	vector.DrawFilledRect(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), float32(p.Object.Size.X), float32(p.Object.Size.Y), cl, false)
}

func NewPlayer(space *resolv.Space) *Player {
	p := &Player{
		Object: resolv.NewObject(32, 0, 100, 140),
	}

	p.Object.SetShape(resolv.NewRectangle(0, 0, p.Object.Size.X, p.Object.Size.Y))

	space.Add(p.Object)

	return p
}
