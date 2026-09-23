package main

import (
	"log"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Difficulty struct{ W, H, Mines int }

var (
	Easy   = Difficulty{9, 9, 10}
	Medium = Difficulty{16, 16, 40}
	Hard   = Difficulty{30, 16, 99}
)

func main() {
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.
	// xmax, ymax := s.Size()

	// Here's an example of how to inject a keystroke where it will
	// be picked up by a future read of the event queue.  Note that
	// care should be used to avoid blocking writes to the queue if
	// this is done from the same thread that is responsible for reading
	// the queue, or else a single-party deadlock might occur.
	// s.EventQ() <- tcell.NewEventKey(tcell.KeyRune, rune('a'), 0)

	// Event loop
	grid := NewGrid(Easy.W, Easy.H, Easy.Mines)
	const cellW, ox, oy = 2, 0, 0

	DrawGrid(s, &grid)
	s.Show()

	for {
		// Poll event (this can be in a select statement as well)
		ev := <-s.EventQ()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			s.Clear()
			DrawGrid(s, &grid)
			s.Show()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
				s.Clear()
				DrawGrid(s, &grid)
				s.Show()
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			b := ev.Buttons()

			bx := (x - ox) / cellW
			by := y - oy
			if by < 0 || by >= len(grid.Cells) || bx < 0 || (len(grid.Cells) > 0 && bx >= len(grid.Cells[0])) {
				continue
			}

			if b&tcell.Button1 != 0 {
				c := &grid.Cells[by][bx]
				if c.State == Hidden {
					c.State = Revealed
					DrawCell(s, *c)
					s.Show()
				}
			} else if b&tcell.Button2 != 0 {
				c := &grid.Cells[by][bx]
				if c.State == Hidden {
					c.State = Flagged
					DrawCell(s, *c)
					s.Show()
				} else if c.State == Flagged {
					c.State = Hidden
					DrawCell(s, *c)
					s.Show()
				}
			}
		}
	}
}
