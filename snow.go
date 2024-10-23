package main

import (
	"container/list"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jakecoffman/cp/v2"
	"image/color"
	"math/rand"
)

type SnowHolder struct {
	Snows *list.List
}

func (s *SnowHolder) Init(mat ebiten.GeoM) {
	for i := 0; i < 5000; i++ {
		s.Update(mat)
	}
}

func (s *SnowHolder) Update(mat ebiten.GeoM) {
	if s.Snows == nil {
		s.Snows = list.New()
	}

	if s.Snows.Len() < 5000 && rand.Intn(4) < 3 {
		s.Snows.PushBack(NewSnow(mat))
	}

	for e := s.Snows.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Snow)
		t.Update()
		if t.Terminated() {
			s.Snows.Remove(e)
		}
	}
}

func (s *SnowHolder) Draw(screen *ebiten.Image, mat ebiten.GeoM) {
	for e := s.Snows.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Snow)
		t.Draw(screen, mat)
	}
}

type Snow struct {
	life      int
	startLife int
	pos       cp.Vector
	lastX     float64
	speed     float64
	size      int
}

func NewSnow(mat ebiten.GeoM) *Snow {
	life := rand.Intn(600) + 2000
	x := -3000 + rand.Float64()*(ScreenWidth*15)
	y := -1000.0

	speed := 0.4 + rand.Float64()*(0.8-0.4)

	size := rand.Intn(2) + 2

	x2, y2 := mat.Apply(x, y)

	return &Snow{
		pos:       cp.Vector{X: x2, Y: y2},
		startLife: life,
		life:      life,
		speed:     speed,
		size:      size,
	}
}

func (s *Snow) Update() {
	if s.life == 0 {
		return
	}
	s.life--

	dir := -0.7
	x := s.pos.X + (dir + rand.Float64())
	y := s.pos.Y + s.speed
	s.pos = cp.Vector{X: x, Y: y}
}

func (s *Snow) Terminated() bool {
	return s.life == 0
}

func (s *Snow) Draw(screen *ebiten.Image, mat ebiten.GeoM) {
	x, y := mat.Apply(s.pos.X, s.pos.Y)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(s.size), float32(s.size), color.Gray{Y: 255}, false)
}
