package main

import (
	"fmt"
	"math/rand"

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
	State       CellState
	IsMine      bool
	MinesAround uint8
}

type Grid struct {
	Cells [][]Cell
}

func (g *Grid) ForEachNeighbor(x, y int, fn func(nx, ny int, n *Cell)) {
	if len(g.Cells) == 0 {
		return
	}
	H, W := len(g.Cells), len(g.Cells[0])
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < W && ny >= 0 && ny < H {
				fn(nx, ny, &g.Cells[ny][nx])
			}
		}
	}
}
func shouldPlaceMine(minesLeft, cellsLeft int) bool {
	if minesLeft <= 0 {
		return false
	}
	return rand.Intn(cellsLeft) < minesLeft
}

func NewGrid(W, H, Mines int) Grid {
	if Mines < 0 {
		Mines = 0
	}
	if Mines > W*H {
		Mines = W * H
	}
	minesLeft := Mines
	cellsLeft := W * H

	cells := make([][]Cell, H)
	for y := range cells {
		// Zero value is already Hidden, no mine, count 0.
		cells[y] = make([]Cell, W)
	}
	grid := Grid{Cells: cells}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if shouldPlaceMine(minesLeft, cellsLeft) {
				cells[y][x].IsMine = true
				grid.ForEachNeighbor(x, y, func(_, _ int, n *Cell) {
					n.MinesAround++
				})
				minesLeft--
			}
			cellsLeft--
		}
	}
	return grid
}

func DrawGrid(s tcell.Screen, g *Grid) {
	for y, row := range g.Cells {
		for x, c := range row {
			DrawCell(s, c, x, y)
		}
	}
}

func DrawCell(s tcell.Screen, c Cell, x, y int) {
	const cellW, ox, oy = 2, 0, 0

	sx := ox + x*cellW
	sy := oy + y
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

// Reveal floods from a Hidden zero cell. Caller checks IsMine/MinesAround.
func (g *Grid) Reveal(x, y int) [][2]int {
	c := &g.Cells[y][x]
	c.State = Revealed
	changed := [][2]int{{x, y}}
	g.ForEachNeighbor(x, y, func(nx, ny int, n *Cell) {
		if n.State != Hidden || n.IsMine {
			return
		}
		if n.MinesAround > 0 {
			n.State = Revealed
			changed = append(changed, [2]int{nx, ny})
			return
		}
		changed = append(changed, g.Reveal(nx, ny)...)
	})
	return changed
}
