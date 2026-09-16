package main

import (
	"example.com/game/internal/ascii"
	"example.com/game/internal/renderer"
)

func drawFrame(blocks, text *ascii.Buffer, app *AppState) {
	blocks.Clear()
	text.Clear()
	if app.Screen == ScreenArena {
		drawArena(blocks, text, app)
	}
	drawAppMenu(text, app)

	renderer.BeginFrame(windowWidth, windowHeight)
	renderer.RenderBlocks(blocks, renderer.Config{
		WindowWidth: windowWidth, WindowHeight: windowHeight,
		CellWidth: blockCellWidth, CellHeight: blockCellHeight,
	})
	renderer.RenderText(text, renderer.Config{
		WindowWidth: windowWidth, WindowHeight: windowHeight,
		CellWidth: textCellWidth, CellHeight: textCellHeight,
	})
}

// drawArena is the extension point for the new game's visuals.
// Leaving both buffers clear produces an empty black arena.
func drawArena(blocks, text *ascii.Buffer, app *AppState) {}
