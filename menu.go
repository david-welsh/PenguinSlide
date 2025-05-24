package main

import (
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
	buttons := make([]*Button, len(menuItems))
	for i, item := range menuItems {
		buttons[i] = &Button{
			X:      ScreenWidth / 4,
			Y:      50 + height*i,
			Action: item.Action,
			OnMouseOver: func(b bool) {
				for j, button := range buttons {
					if b {
						if j == i {
							button.Highlighted = true
						} else {
							button.Highlighted = false
						}
					} else {
						if j == i {
							button.Highlighted = false
						}
					}
				}
			},
			Text:   item.Title,
			Width:  ScreenWidth / 2,
			Height: height,
		}
	}

	return &Menu{
		Buttons:  buttons,
		Selected: 0,
		bgColor:  bgCol,
	}
}

func (m *Menu) Update() error {
	for _, button := range m.Buttons {
		button.Update()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		selected := 0
		for i, button := range m.Buttons {
			if button.Highlighted {
				selected = i
				break
			}
		}
		nextSelect := (selected + 1) % len(m.Buttons)
		m.Buttons[selected].Highlighted = false
		m.Buttons[nextSelect].Highlighted = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		selected := 0
		for i, button := range m.Buttons {
			if button.Highlighted {
				selected = i
				break
			}
		}
		nextSelect := selected
		nextSelect -= 1
		if nextSelect == -1 {
			nextSelect = len(m.Buttons) - 1
		}
		m.Buttons[selected].Highlighted = false
		m.Buttons[nextSelect].Highlighted = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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
	return ""
}
