package main

import (
	"PenguinSlide/assets/fonts"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	X            int
	Y            int
	Action       func()
	OnMouseOver  func(bool)
	Text         string
	Width        int
	Height       int
	Highlighted  bool
	mouseWasOver bool
}

func (b *Button) Update() {
	x, y := ebiten.CursorPosition()
	if b.inButtonBounds(x, y) {
		if !b.mouseWasOver {
			b.OnMouseOver(true)
			b.mouseWasOver = true
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			b.Action()
		}
	} else {
		if b.mouseWasOver {
			b.mouseWasOver = false
			b.OnMouseOver(false)
		}
	}
}

func (b *Button) Draw(screen *ebiten.Image) error {
	normalFontFace := &text.GoTextFace{
		Source: fonts.FontFaceSource(),
		Size:   float64(b.Height / 3),
	}
	selectedFontFace := &text.GoTextFace{
		Source: fonts.FontFaceSource(),
		Size:   float64(b.Height/2) * 0.75,
	}

	fontFace := normalFontFace
	if b.Highlighted {
		fontFace = selectedFontFace
	}
	op := &text.DrawOptions{}
	x := float64(b.X) + (0.5 * float64(b.Width))
	op.GeoM.Translate(x, float64(b.Y))
	op.ColorScale.ScaleWithColor(TextColor)
	op.PrimaryAlign = text.AlignCenter

	st := b.Text
	if b.Highlighted {
		st = fmt.Sprintf("> %s <", b.Text)
	}
	text.Draw(screen, st, fontFace, op)

	return nil
}

func (b *Button) inButtonBounds(x, y int) bool {
	if x < b.X || x > b.X+b.Width {
		return false
	} else if y < b.Y || y > b.Y+b.Height {
		return false
	}
	return true
}
