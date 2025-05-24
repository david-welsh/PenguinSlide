package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
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
)

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
	Buttons  []*Button
	Selected int
	bgColor  color.RGBA
}

func NewMenu(bgColor *color.RGBA, menuItems ...*MenuItem) *Menu {
	bgCol := DefaultBgColor
	if bgColor != nil {
		bgCol = *bgColor
	}
	height := (ScreenHeight - 100) / len(menuItems)
	menu := &Menu{
		Selected: 0,
		bgColor:  bgCol,
	}
	buttons := make([]*Button, len(menuItems))
	for i, item := range menuItems {
		newButton := &Button{
			X:      ScreenWidth / 4,
			Y:      50 + height*i,
			Action: item.Action,
			OnMouseOver: func(b bool) {
				if b {
					menu.Selected = i
				}
				menu.updateHighlighting()
			},
			Text:   item.Title,
			Width:  ScreenWidth / 2,
			Height: height,
		}
		buttons[i] = newButton
	}
	menu.Buttons = buttons

	return menu
}

func (m *Menu) updateHighlighting() {
	for i, button := range m.Buttons {
		button.Highlighted = i == m.Selected
	}
}

func (m *Menu) Update() error {
	for _, button := range m.Buttons {
		button.Update()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.Selected = (m.Selected + 1) % len(m.Buttons)
		m.updateHighlighting()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.Selected -= 1
		if m.Selected == -1 {
			m.Selected = len(m.Buttons) - 1
		}
		m.updateHighlighting()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		m.Buttons[m.Selected].Action()
	}

	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) error {
	vector.DrawFilledRect(screen, -5, -5, ScreenWidth+10, ScreenHeight+10, m.bgColor, false)

	for _, button := range m.Buttons {
		err := button.Draw(screen)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Menu) GenerateDebugString() string {
	return fmt.Sprintf("Selected item: %d", m.Selected)
}
