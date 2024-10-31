package main

import (
	"PenguinSlide/assets"
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

var (
	DefaultBgColor = color.RGBA{
		R: 180,
		G: 180,
		B: 200,
		A: 255,
	}
	TextColor = color.RGBA{
		R: 165,
		G: 110,
		B: 145,
		A: 255,
	}
	fontFaceSource *text.GoTextFaceSource
)

func MenuInit() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(assets.MenuFontTtf))
	if err != nil {
		log.Fatal(err)
	}

	fontFaceSource = s
}

type MenuItem struct {
	Action func()
	Title  string
}

func NewMenuItem(title string, action func()) *MenuItem {
	return &MenuItem{
		Title:  title,
		Action: action,
	}
}

type Menu struct {
	MenuItems []*MenuItem
	Selected  int
	bgColor   color.RGBA
}

func NewMenu(bgColor *color.RGBA, menuItems ...*MenuItem) *Menu {
	bgCol := DefaultBgColor
	if bgColor != nil {
		bgCol = *bgColor
	}
	return &Menu{
		MenuItems: menuItems,
		Selected:  0,
		bgColor:   bgCol,
	}
}

func (m *Menu) Update() error {
	x, y := ebiten.CursorPosition()
	if i := m.inButtonBounds(x, y); i != nil {
		m.Selected = *i
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.MenuItems[m.Selected].Action()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.Selected += 1
		m.Selected %= len(m.MenuItems)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.Selected -= 1
		if m.Selected == -1 {
			m.Selected = len(m.MenuItems) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		m.MenuItems[m.Selected].Action()
	}

	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) error {
	vector.DrawFilledRect(screen, -5, -5, ScreenWidth+10, ScreenHeight+10, m.bgColor, false)

	height := (ScreenHeight - 100) / len(m.MenuItems)

	normalFontFace := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   float64(height / 3),
	}
	selectedFontFace := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   float64(height/2) * 0.75,
	}

	for i, item := range m.MenuItems {
		fontFace := normalFontFace
		if i == m.Selected {
			fontFace = selectedFontFace
		}
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(ScreenWidth/2), float64(50+(height*i)))
		op.ColorScale.ScaleWithColor(TextColor)
		op.PrimaryAlign = text.AlignCenter
		st := item.Title
		if i == m.Selected {
			st = fmt.Sprintf("> %s <", item.Title)
		}
		text.Draw(screen, st, fontFace, op)
	}

	return nil
}

func (m *Menu) getButtonBounds(i int) (t, b, l, r float64) {
	height := (ScreenHeight - 100) / len(m.MenuItems)
	t = float64(50+height*i) - 10
	b = t + float64(height/2)*0.75 + 10
	l = 50
	r = ScreenWidth - 50

	return t, b, l, r
}

func (m *Menu) inButtonBounds(x, y int) *int {
	for i := range m.MenuItems {
		t, b, l, r := m.getButtonBounds(i)
		if float64(x) > l && float64(x) < r && float64(y) > t && float64(y) < b {
			return &i
		}
	}
	return nil
}

func (m *Menu) GenerateDebugString() string {
	return ""
}
