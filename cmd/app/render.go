package main

import (
	"math"

	"github.com/andrlzpt/dim/internal/ascii"
	"github.com/andrlzpt/dim/internal/renderer"
)

func drawFrame(blocks, world, text *ascii.Buffer, app *AppState) {
	blocks.Clear()
	world.Clear()
	text.Clear()
	if app.Screen == ScreenArena {
		drawArena(world, app)
	}
	drawAppMenu(text, app)

	renderer.BeginFrame(windowWidth, windowHeight)
	renderer.RenderBlocks(blocks, renderer.Config{
		WindowWidth: windowWidth, WindowHeight: windowHeight,
		CellWidth: blockCellWidth, CellHeight: blockCellHeight,
	})
	renderer.RenderText(world, renderer.Config{
		WindowWidth: windowWidth, WindowHeight: windowHeight,
		CellWidth: worldCellWidth, CellHeight: worldCellHeight,
	})
	renderer.RenderText(text, renderer.Config{
		WindowWidth: windowWidth, WindowHeight: windowHeight,
		CellWidth: textCellWidth, CellHeight: textCellHeight,
	})
}

// drawArena is the extension point for the new game's visuals.
// Leaving both buffers clear produces an empty black arena.
func drawArena(world *ascii.Buffer, app *AppState) {
	const horizontalFOV = math.Pi / 3 // 60 degrees.

	camera := app.Arena.Camera
	forwardX := math.Cos(camera.Angle)
	forwardY := math.Sin(camera.Angle)
	planeScale := math.Tan(horizontalFOV / 2)

	// The camera plane points to the right of the viewing direction.
	planeX := -forwardY * planeScale
	planeY := forwardX * planeScale

	// Convert the projection distance from pixels to character rows.
	projection := float64(world.Width*worldCellWidth) /
		(2 * planeScale * float64(worldCellHeight))
	horizon := float64(world.Height) / 2

	for x := 0; x < world.Width; x++ {
		screenX := 2*(float64(x)+0.5)/float64(world.Width) - 1
		rayX := forwardX + planeX*screenX
		rayY := forwardY + planeY*screenX

		distance, side := castWall(&app.Arena, rayX, rayY)
		distance = max(distance, 0.001)

		const wallHeight = 3.0
		height := wallHeight * projection / distance
		top := horizon - height/2
		bottom := horizon + height/2

		// Include rows whose centers fall inside the wall slice.
		start := max(0, int(math.Ceil(top-0.5)))
		end := min(world.Height, int(math.Ceil(bottom-0.5)))
		glyph := wallGlyph(distance, side)

		for y := start; y < end; y++ {
			world.Set(x, y, glyph)
		}
	}
}

func wallGlyph(distance float64, side bool) rune {
	const shades = ".:-=+*#%@"

	brightness := 1 / (1 + distance*0.18)
	if side {
		brightness *= 0.7
	}

	index := int(brightness * float64(len(shades)-1))
	return rune(shades[index])
}
