package main

import "image/color"

var (
	neighbours = [8][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	LiveColor   = color.RGBA{R: 136, G: 0, B: 255, A: 255}
	ShadowColor = color.RGBA{R: 255, G: 255, B: 255, A: 100}
	BorderColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	Rotate = 0
	Glider = [][][]int{
		{
			{1, 0}, {1, 1}, {-1, 0}, {0, 1}, {1, -1},
		},
		{
			{0, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1},
		},
		{
			{0, -1}, {-1, -1}, {1, -1}, {1, 0}, {0, 1},
		},
		{
			{0, -1}, {-1, 1}, {-1, 0}, {-1, -1}, {1, 0},
		},
	}
)

type Cell struct {
	status      bool
	nextStatus  bool
	aroundCells []*Cell
	shadow      bool
}

func (c *Cell) Switch() {
	c.status = !c.status
}

func (c *Cell) UnLink() {
	c.aroundCells = []*Cell{}
}

func (c *Cell) Link(aroundCells []*Cell) {
	c.aroundCells = aroundCells
}

func (c *Cell) Status() bool {
	return c.status
}

func (c *Cell) Shadow() bool {
	return c.shadow
}

func (c *Cell) SetStatus(status bool) {
	c.status = status
}

func (c *Cell) SetShadow(shadow bool) {
	c.shadow = shadow
}

func (c *Cell) CalcNextState() {
	var nextStatus bool
	liveCells := 0
	for _, ac := range c.aroundCells {
		if ac.status {
			liveCells++
		}
	}

	if !c.status && liveCells == 3 {
		nextStatus = true
	}
	if c.status && (liveCells < 2 || liveCells > 3) {
		nextStatus = false
	}
	if c.status && (liveCells == 2 || liveCells == 3) {
		nextStatus = true
	}
	c.nextStatus = nextStatus
}

func (c *Cell) Flush() {
	c.status = c.nextStatus
}

func (c *Cell) ClearShadow() {
	c.shadow = false
}

func (c *Cell) Clear() {
	c.shadow = false
	c.status = false
}
