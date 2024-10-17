package main

import (
	"PenguinSlide/assets"
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp/v2"
	"image"
	"log"
	"math"
)

type Player struct {
	Space       *cp.Space
	Body        *cp.Body
	Shape       *cp.Shape
	CurrentBox  cp.Vector
	OnGround    bool
	FacingRight bool
	Sliding     bool
	Game        *Game
}

var (
	walkingImage *ebiten.Image
	slidingImage *ebiten.Image
	walkBox      cp.Vector
	slideBox     cp.Vector
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

	walkBox = cp.Vector{X: 50, Y: 70}
	slideBox = cp.Vector{X: 85, Y: 35}
}

func (p *Player) Update() {
	maxWalkSpeed := 100.0
	maxSlideSpeed := 300.0

	maxSpeed := maxWalkSpeed
	if p.Sliding {
		maxSpeed = maxSlideSpeed
	}
	if math.Abs(p.Body.Velocity().X) > maxSpeed {
		vel := p.Body.Velocity()
		p.Body.SetVelocity(math.Copysign(maxSpeed, vel.X), vel.Y)
	}

	wasSliding := p.Sliding
	if ebiten.IsKeyPressed(ebiten.KeySpace) || math.Abs(p.Body.Velocity().X) > maxWalkSpeed {
		p.Sliding = true
		if !wasSliding {
			p.Shape.Space().RemoveShape(p.Shape)
			p.Shape = p.Space.AddShape(cp.NewBox(p.Body, slideBox.X, 1, slideBox.Y))
			p.Shape.SetFriction(0.3)
			p.Shape.SetElasticity(0)
			newPos := p.Body.Position()
			newPos.Y += p.CurrentBox.Y / 2
			p.Body.SetPosition(newPos)
			p.CurrentBox = slideBox
		}
	} else {
		p.Sliding = false
		if wasSliding {
			p.Shape.Space().RemoveShape(p.Shape)
			p.Shape = p.Space.AddShape(cp.NewBox(p.Body, walkBox.X/2, walkBox.Y, walkBox.X/2))
			p.Shape.SetFriction(0.7)
			p.Shape.SetElasticity(0)
			newPos := p.Body.Position()
			newPos.Y -= p.CurrentBox.Y
			p.Body.SetPosition(newPos)
			p.CurrentBox = walkBox
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxisValue(0, 0) > 0.1 {
		p.Body.ApplyForceAtLocalPoint(cp.Vector{X: 750, Y: 0}, cp.Vector{})
		p.FacingRight = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxisValue(0, 0) < -0.1 {
		p.Body.ApplyForceAtLocalPoint(cp.Vector{X: -750, Y: 0}, cp.Vector{})
		p.FacingRight = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsGamepadButtonJustPressed(0, 0) {
		p.Body.ApplyImpulseAtLocalPoint(cp.Vector{X: 0, Y: -400}, cp.Vector{})
	}
}

func (p *Player) Draw(screen *ebiten.Image, mat ebiten.GeoM) {
	img := walkingImage
	if p.Sliding {
		img = slidingImage
	}

	op := &ebiten.DrawImageOptions{}
	sign := 1
	if !p.FacingRight {
		sign = -1
	}
	op.GeoM.Scale(0.093*float64(sign), 0.093)
	op.GeoM.Translate(-float64(sign)*p.CurrentBox.X, -p.CurrentBox.Y)
	op.GeoM.Rotate(p.Body.Angle())
	op.GeoM.Translate(p.Body.Position().X, p.Body.Position().Y)
	op.GeoM.Concat(mat)
	screen.DrawImage(img, op)
}

func (p *Player) GenerateDebugText() string {
	speed := fmt.Sprintf("Speed: %.1f, %.1f", p.Body.Velocity().X, p.Body.Velocity().Y)
	position := fmt.Sprintf("X: %.1f, %.1f", p.Body.Position().X, p.Body.Position().Y)
	return fmt.Sprintf("%s\n%s", speed, position)
}

func NewPlayer(space *cp.Space, game *Game, pos cp.Vector) *Player {
	body := space.AddBody(cp.NewBody(2.0, 500))
	body.SetPosition(pos)

	walkingShape := space.AddShape(cp.NewBox(body, walkBox.X/2, walkBox.Y, walkBox.X/2))
	walkingShape.SetFriction(0.7)

	return &Player{
		Space:       space,
		Body:        body,
		Shape:       walkingShape,
		Game:        game,
		CurrentBox:  walkBox,
		FacingRight: true,
	}
}
