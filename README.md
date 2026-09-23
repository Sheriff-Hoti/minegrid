# minegrid

Terminal Minesweeper in Go with `tcell`.

## Run

```sh
go run .
```

## Controls

- Left click: reveal (empty area flood-fills)
- Right click: flag / unflag
- `Esc` / `Ctrl+C`: quit
- `Ctrl+L`: redraw

## Difficulties (`main.go`)

- Easy: 9x9, 10 mines
- Medium: 16x16, 40 mines
- Hard: 30x16, 99 mines

Currently starts on Easy.

## Layout

- `main.go`: screen setup, event loop, incremental redraw
- `grid.go`: `Grid`/`Cell` model, sequential mine placement (`shouldPlaceMine`), neighbor counts via `ForEachNeighbor`, recursive zero `Reveal`, `DrawGrid`/`DrawCell`
