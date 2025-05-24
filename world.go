package main

import (
	"PenguinSlide/assets"
	"encoding/csv"
	"fmt"
	"github.com/demouth/ebitencp"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp/v2"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
)

var (
	SkyColor = color.RGBA{
		R: 180,
		G: 180,
		B: 200,
		A: 255,
	}
	groundImage *ebiten.Image
)

type World struct {
	Game       *Game
	Level      string
	Space      *cp.Space
	Player     *Player
	Finish     cp.Vector
	Drawer     *ebitencp.Drawer
	Camera     Camera
	SnowHolder SnowHolder
	Paused     bool
	Menu       *Menu
}

func NewWorld(game *Game, level string) *World {
	w := &World{Game: game, Level: level}
	w.Init()
	return w
}

type LevelSegment struct {
	A cp.Vector
	B cp.Vector
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Failed to close level file %s", filePath)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func parseFloatIgnore(f string) float64 {
	f64, _ := strconv.ParseFloat(f, 32)
	return f64
}

func ParseLevel(level string) (segments []LevelSegment, playerPos cp.Vector, finishPos cp.Vector, maxX float64) {
	data := readCsvFile(fmt.Sprintf("levels/%s.csv", level))
	for _, seg := range data {
		if seg[0] == "line" {
			levelSegment := LevelSegment{
				A: cp.Vector{X: parseFloatIgnore(seg[1]), Y: parseFloatIgnore(seg[2])},
				B: cp.Vector{X: parseFloatIgnore(seg[3]), Y: parseFloatIgnore(seg[4])},
			}
			segments = append(segments, levelSegment)
			segmentMaxX := math.Max(levelSegment.A.X, levelSegment.B.X)
			if segmentMaxX > maxX {
				maxX = segmentMaxX
			}
		} else if seg[0] == "player" {
			playerPos.X = parseFloatIgnore(seg[1])
			playerPos.Y = parseFloatIgnore(seg[2])
		} else if seg[0] == "finish" {
			finishPos.X = parseFloatIgnore(seg[1])
			finishPos.Y = parseFloatIgnore(seg[2])
		}
	}

	return segments, playerPos, finishPos, maxX
}

func (world *World) Reset() {
	world.Init()
}

func (world *World) Init() {
	world.Menu = NewMenu(
		&color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 100,
		},
		NewMenuItem("Return", func() {
			world.Paused = false
		}),
		NewMenuItem("Main Menu", func() {
			world.Game.LoadMenu()
		}),
		NewMenuItem("Quit Game", func() {
			world.Game.ShouldQuit = true
		}),
	)

	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)

	world.Camera = Camera{
		Zoom:   1.0,
		Offset: cp.Vector{X: -gw / 2, Y: -gh / 2},
	}

	world.Drawer = ebitencp.NewDrawer(0, 0)
	world.Drawer.FlipYAxis = true
	world.Drawer.GeoM.Translate(gw/2, gh/2)

	world.Space = cp.NewSpace()

	world.Space.SleepTimeThreshold = 0.5
	world.Space.SetCollisionSlop(0.5)

	levelSegments, playerPos, finishPos, maxX := ParseLevel(world.Level)

	world.Finish = finishPos

	for _, segment := range levelSegments {
		shape := cp.NewSegment(world.Space.StaticBody, segment.A, segment.B, 0)
		shape.SetFriction(0.1)
		shape.SetElasticity(0)
		shape.SetCollisionType(2)
		world.Space.AddShape(shape)
	}

	world.Space.Iterations = 30
	world.Space.SetGravity(cp.Vector{Y: 400})
	world.Space.SetCollisionSlop(0.5)

	world.SnowHolder.Init(*world.Drawer.GeoM, maxX)

	world.Player = NewPlayer(world.Space, world.Game, playerPos)
}

func (world *World) GenerateDebugString() string {
	worldDebug := ""
	cameraDebug := world.Camera.GenerateDebugText()
	playerDebug := world.Player.GenerateDebugText()
	return fmt.Sprintf("%s\n%s\n%s", worldDebug, cameraDebug, playerDebug)
}

func (world *World) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		world.Paused = !world.Paused
	}

	if world.Paused {
		err := world.Menu.Update()
		if err != nil {
			return err
		}
		return nil
	}

	world.Player.Update()

	world.SnowHolder.Update(*world.Drawer.GeoM)

	world.Space.Step(1.0 / float64(ebiten.TPS()))

	world.Camera.Offset.X = -world.Player.Body.Position().X
	world.Camera.Offset.Y = -world.Player.Body.Position().Y
	world.Camera.Zoom = 1.0 - ((math.Abs(world.Player.Body.Velocity().X) / 300) * 0.2)

	world.Drawer.GeoM.Reset()
	world.Drawer.GeoM.Translate(world.Camera.Offset.X, world.Camera.Offset.Y)
	world.Drawer.GeoM.Scale(world.Camera.Zoom, world.Camera.Zoom)
	world.Drawer.GeoM.Rotate(world.Camera.Rotate)
	world.Drawer.GeoM.Translate(float64(world.Game.Width)/2, float64(world.Game.Height)/2)

	return nil
}

func drawWorldSegment(screen *ebiten.Image, ta, tb cp.Vector, mat ebiten.GeoM) {
	if groundImage == nil {
		groundImage = LoadImage(assets.GroundPng)
	}
	for _, drawPoint := range drawPoints(ta, tb) {
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.2, 0.2)
		op.GeoM.Translate(drawPoint.X, drawPoint.Y)
		op.GeoM.Concat(mat)
		screen.DrawImage(groundImage, &op)
	}
}

func drawPoints(ta, tb cp.Vector) []cp.Vector {
	v := cp.Vector{
		X: tb.X - ta.X,
		Y: tb.Y - ta.Y,
	}
	l := vectorLength(v)
	u := v.Mult(1 / l)
	d := 5.0
	count := int(l/d) + 1

	points := make([]cp.Vector, count)
	for i := 0; i < count; i++ {
		points[i] = ta.Add(u.Mult(float64(i) * d))
	}
	return points
}

func vectorLength(v cp.Vector) float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func (world *World) Draw(screen *ebiten.Image) error {
	screen.Fill(SkyColor)

	world.Space.EachShape(func(s *cp.Shape) {
		switch s.Class.(type) {
		case *cp.Segment:
			segment := s.Class.(*cp.Segment)
			ta := segment.TransformA()
			tb := segment.TransformB()
			drawWorldSegment(screen, ta, tb, *world.Drawer.GeoM)
		}
	})

	world.Player.Draw(screen, *world.Drawer.GeoM)

	world.SnowHolder.Draw(screen, *world.Drawer.GeoM)

	if world.Game.Debug {
		cp.DrawSpace(
			world.Space,
			world.Drawer.WithScreen(screen),
		)
	}

	if world.Paused {
		err := world.Menu.Draw(screen)
		if err != nil {
			return err
		}
	}

	return nil
}
