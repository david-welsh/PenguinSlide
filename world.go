package main

import (
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
)

type World struct {
	Game       *Game
	Level      string
	Space      *cp.Space
	Player     *Player
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

func ParseLevel(level string) (segments []LevelSegment, playerPos cp.Vector) {
	data := readCsvFile(fmt.Sprintf("levels/%s.csv", level))
	for _, seg := range data {
		if seg[0] == "line" {
			segments = append(segments, LevelSegment{
				A: cp.Vector{X: parseFloatIgnore(seg[1]), Y: parseFloatIgnore(seg[2])},
				B: cp.Vector{X: parseFloatIgnore(seg[3]), Y: parseFloatIgnore(seg[4])},
			})
		} else if seg[0] == "player" {
			playerPos.X = parseFloatIgnore(seg[1])
			playerPos.Y = parseFloatIgnore(seg[2])
		}
	}

	return segments, playerPos
}

func (world *World) Reset() {
	world.Init()
}

func (world *World) Init() {
	world.Menu = NewMenu(
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

	levelSegments, playerPos := ParseLevel(world.Level)

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

	world.SnowHolder.Init(*world.Drawer.GeoM)

	world.Player = NewPlayer(world.Space, world.Game, playerPos)
}

func (world *World) GenerateDebugString() string {
	worldDebug := ""
	cameraDebug := world.Camera.GenerateDebugText()
	playerDebug := world.Player.GenerateDebugText()
	return fmt.Sprintf("%s\n%s\n%s", worldDebug, cameraDebug, playerDebug)
}

func (world *World) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
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

func (world *World) Draw(screen *ebiten.Image) error {
	screen.Fill(SkyColor)

	if world.Game.Debug {
		cp.DrawSpace(
			world.Space,
			world.Drawer.WithScreen(screen),
		)
	}

	world.Player.Draw(screen, *world.Drawer.GeoM)

	world.SnowHolder.Draw(screen, *world.Drawer.GeoM)

	if world.Paused {
		err := world.Menu.Draw(screen)
		if err != nil {
			return err
		}
	}

	return nil
}
