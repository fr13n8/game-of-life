package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"golang.org/x/image/math/f64"

	"github.com/fr13n8/game-of-life/camera"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type World struct {
	Cells      [][]Cell
	Width      int
	Height     int
	Percentage float64
	camera     camera.Camera
	LifeTime   int
	Paused     bool
}

func NewWorld(width int, height int) *World {
	world := &World{
		Width:  width,
		Height: height,
		camera: camera.Camera{ViewPort: f64.Vec2{ScreenWidth, ScreenHeight}},
		Paused: true,
	}
	return world
}

func (w *World) Run() {
	w.Init()
}

func (w *World) GetCursorCoordinates() (int, int) {
	worldX, worldY := w.camera.ScreenToWorld(ebiten.CursorPosition())
	x, y := int(worldX), int(worldY)
	return x, y
}

func (w *World) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := w.GetCursorCoordinates()
		if y <= w.Height && x <= w.Width && y >= 0 && x >= 0 {
			if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
				w.Get(x, y).SetStatus(false)
				return nil
			}
			w.Get(x, y).SetStatus(true)
		}
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		x, y := w.GetCursorCoordinates()
		glider := Glider[Rotate]
		w.ClearShadows()

		if y <= w.Height && x <= w.Width && y >= 0 && x >= 0 {
			for _, coords := range glider {
				w.Get(x+coords[0], y+coords[1]).SetShadow(true)
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			if Rotate == 3 {
				Rotate = 0
			} else {
				Rotate++
			}
		}
	}

	if inpututil.IsKeyJustReleased(ebiten.Key1) {
		x, y := w.GetCursorCoordinates()
		glider := Glider[Rotate]

		if y <= w.Height && x <= w.Width && y >= 0 && x >= 0 {
			for _, coords := range glider {
				w.Get(x+coords[0], y+coords[1]).SetShadow(false)
				w.Get(x+coords[0], y+coords[1]).SetStatus(true)
			}
		}
		w.ClearShadows()
	}

	if ebiten.IsKeyPressed(ebiten.KeyG) {
		w.Clear()
		w.Random()
	}
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		w.Clear()
		w.Paused = true
		w.LifeTime = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		w.Next()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		w.Paused = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyNumpadEnter) {
		w.Paused = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		w.camera.Position[0] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		w.camera.Position[0] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		w.camera.Position[1] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		w.camera.Position[1] += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if w.camera.ZoomFactor > -2400 {
			w.camera.ZoomFactor -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if w.camera.ZoomFactor < 2400 {
			w.camera.ZoomFactor += 1
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		w.camera.Rotation += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		w.camera.Reset()
	}

	return nil
}

func (w *World) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	background := ebiten.NewImage(ScreenWidth, ScreenHeight)

	w.Print(background)
	if !w.Paused {
		w.Next()
	}
	w.camera.Render(background, screen)

	worldX, worldY := w.camera.ScreenToWorld(ebiten.CursorPosition())
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("TPS: %0.2f\n"+
			"Pause/Play (ESC)(Enter)\n"+
			"Move (WASD/Arrows)\n"+
			"Zoom (QE)\n"+
			"Rotate (R)\n"+
			"Reset camera (Space)\n"+
			"Random fill (G)\n"+
			"Reset world (C)\n"+
			"Draw/Eraser (LCM)/(LShift+LCM)\n"+
			"Rotate model (modelKey+RCM)\n"+
			"Next step (N)\n", ebiten.CurrentTPS()),
		2, 1,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f\nWorld life time: %d | Models(numpad): 1-Glider",
			w.camera.String(),
			worldX, worldY, w.LifeTime),
		2, ScreenHeight-48,
	)
}

func Cells(width, height int) [][]Cell {
	cells := make([][]Cell, height)
	for i := 0; i < height; i++ {
		cells[i] = make([]Cell, width)
	}

	return cells
}

func (w *World) Init() {
	w.Cells = Cells(w.Width, w.Height)
	w.InitCells()
}

func (w *World) InitCells() {
	around := [8][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			var aroundCells []*Cell
			w.Cells[i][j].UnLink()
			for _, a := range around {
				y := i + a[0]
				x := j + a[1]

				if x < 0 || y < 0 || x >= w.Width || y >= w.Height {
					continue
				}

				aroundCells = append(aroundCells, &w.Cells[y][x])
			}
			w.Cells[i][j].Link(aroundCells)
		}
	}
}

func (w *World) Print(screen *ebiten.Image) {
	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {

			if w.Cells[i][j].Status() {
				ebitenutil.DrawRect(screen, float64(j), float64(i), 1, 1, LiveColor)
			}
			if w.Cells[i][j].Shadow() {
				ebitenutil.DrawRect(screen, float64(j), float64(i), 1, 1, ShadowColor)
			}
			if i == 0 {
				ebitenutil.DrawRect(screen, float64(j), float64(i), 1, 1, BorderColor)
			} else if i == w.Height-1 {
				ebitenutil.DrawRect(screen, float64(j), float64(i), 1, 1, BorderColor)
			} else if j == 0 {
				ebitenutil.DrawRect(screen, float64(j), float64(i), 1, 1, BorderColor)
			} else if j == w.Width-1 {
				ebitenutil.DrawRect(screen, float64(j), float64(i), 1, 1, BorderColor)
			}
		}
	}
}

func (w *World) Next() {
	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			w.Cells[i][j].CalcNextState()
		}
	}

	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			w.Cells[i][j].Flush()
		}
	}
	w.LifeTime++
}

func (w *World) Random() {
	grid := make([][]bool, w.Height)
	for i := 0; i < w.Height; i++ {
		grid[i] = make([]bool, w.Width)
	}

	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			grid[i][j] = rand.Int()%2 == 0
		}
	}
	w.Set(0, 0, grid)
}

func (w *World) Set(x, y int, cells [][]bool) {
	for i, row := range cells {
		for j, state := range row {
			if y+i < w.Height && x+j < w.Width {
				w.Cells[y+i][x+j].SetStatus(state)
			}
		}
	}
}

func (w *World) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (w *World) Get(x, y int) *Cell {
	return &w.Cells[y][x]
}

func (w *World) ClearShadows() {
	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			w.Cells[i][j].ClearShadow()
		}
	}
}

func (w *World) Clear() {
	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			w.Cells[i][j].Clear()
		}
	}
}
