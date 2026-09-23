package main

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type CellState uint8
type CellPayload uint8

const (
	Hidden CellState = iota
	Flagged
	Revealed
)

type Cell struct {
	X, Y        int
	State       CellState
	IsMine      bool
	MinesAround uint8
}

type Grid struct {
	Cells [][]Cell
}

func NewGrid(W, H, Mines int) Grid {
	cells := make([][]Cell, H)
	for y := range cells {
		cells[y] = make([]Cell, W)
		for x := range cells[y] {
			cells[y][x] = Cell{X: x, Y: y, State: Hidden}
		}
	}
	// TODO: place Mines randomly + compute MinesAround
	return Grid{Cells: cells}
}

func DrawGrid(s tcell.Screen, g *Grid) {
	for _, row := range g.Cells {
		for _, c := range row {
			DrawCell(s, c)
		}
	}
}

func DrawCell(s tcell.Screen, c Cell) {
	const cellW, ox, oy = 2, 0, 0

	sx := ox + c.X*cellW
	sy := oy + c.Y
	sw, sh := s.Size()
	if sx < 0 || sy < 0 || sx+cellW > sw || sy >= sh {
		return
	}

	var text string
	style := tcell.StyleDefault

	switch c.State {
	case Hidden:
		text = "■ "
		style = style.Foreground(color.Gray).Background(color.Reset)
	case Flagged:
		text = "F "
		style = style.Foreground(color.Yellow).Background(color.Reset)
	case Revealed:
		if c.IsMine {
			text = "* "
			style = style.Foreground(color.Red).Background(color.Reset)
		} else if c.MinesAround == 0 {
			text = "  "
		} else {
			text = fmt.Sprintf("%d ", c.MinesAround)
			style = style.Foreground(color.White).Background(color.Reset)
		}
	}

	s.Put(sx, sy, text, style)
}
