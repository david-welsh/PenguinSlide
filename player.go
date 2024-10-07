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
	Sliding     bool
}

var (
	walkingImage *ebiten.Image
	slidingImage *ebiten.Image
	walkBox      resolv.Vector
	slideBox     resolv.Vector
)

func loadImage(img []byte) *ebiten.Image {
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

func init() {
	walkingImage = loadImage(assets.WalkingPng)
	slidingImage = loadImage(assets.SlidingPng)

	walkBox = resolv.Vector{X: 100, Y: 140}
	slideBox = resolv.Vector{X: 170, Y: 70}
}

func (p *Player) adjustCollisionBox(newWidth, newHeight float64) {
	oldHeight := p.Object.Size.Y
	oldWidth := p.Object.Size.X

	p.Object.Size.X = newWidth
	p.Object.Size.Y = newHeight

	p.Object.SetShape(resolv.NewRectangle(0, 0, newWidth, newHeight))
	p.Object.Position.Y += oldHeight - newHeight
	if !p.FacingRight {
		p.Object.Position.X -= newWidth - oldWidth
	}
}

func (p *Player) Update() {
	friction := 0.01
	accel := 0.2 + friction
	maxWalkSpeed := 1.0
	maxSlideSpeed := 6.0
	gravity := 0.75

	p.Speed.Y += gravity

	if ebiten.IsKeyPressed(ebiten.KeySpace) || (p.Speed.X > maxWalkSpeed || p.Speed.X < -maxWalkSpeed) {
		p.Sliding = true
		p.adjustCollisionBox(slideBox.X, slideBox.Y)
	} else {
		p.Sliding = false
		p.adjustCollisionBox(walkBox.X, walkBox.Y)
	}

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

	maxSpeed := maxWalkSpeed
	if p.Sliding {
		maxSpeed = maxSlideSpeed
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
	if p.Sliding {
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
