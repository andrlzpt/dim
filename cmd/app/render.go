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

// drawArena projects the room into the world character buffer.
func drawArena(world *ascii.Buffer, app *AppState) {
	const (
		horizontalFOV = math.Pi / 3
		wallHeight    = 3.0
	)

	camera := app.Arena.Camera
	forwardX := math.Cos(camera.Angle)
	forwardY := math.Sin(camera.Angle)
	planeScale := math.Tan(horizontalFOV / 2)

	planeX := -forwardY * planeScale
	planeY := forwardX * planeScale

	projection := float64(world.Width*worldCellWidth) /
		(2 * planeScale * float64(worldCellHeight))
	horizon := float64(world.Height) / 2

	for x := 0; x < world.Width; x++ {
		screenX := 2*(float64(x)+0.5)/float64(world.Width) - 1
		rayX := forwardX + planeX*screenX
		rayY := forwardY + planeY*screenX
		rayLength := math.Hypot(rayX, rayY)

		// Draw quiet floor and ceiling surfaces behind the walls.
		for y := 0; y < world.Height; y++ {
			offset := float64(y) + 0.5 - horizon
			if math.Abs(offset) < 0.001 {
				continue
			}

			depth := (wallHeight / 2) * projection / math.Abs(offset)
			surfaceX := camera.X + rayX*depth
			surfaceY := camera.Y + rayY*depth

			// Use actual horizontal distance for fading.
			distance := depth * rayLength
			brightness := surfaceBrightness(
				surfaceX,
				surfaceY,
				distance,
				offset < 0,
			)
			world.SetShaded(x, y, '.', brightness)
		}

		depth, side := app.Arena.CastWall(rayX, rayY)
		depth = max(depth, 0.001)

		height := wallHeight * projection / depth
		top := horizon - height/2
		bottom := horizon + height/2

		start := max(0, int(math.Ceil(top-0.5)))
		end := min(world.Height, int(math.Ceil(bottom-0.5)))

		brightness := worldBrightness(depth * rayLength)
		if side {
			brightness *= 0.75
		}

		for y := start; y < end; y++ {
			world.SetShaded(x, y, '+', brightness)
		}
	}
}

// worldBrightness fades soft gray smoothly toward black.
func worldBrightness(distance float64) float32 {
	const (
		nearBrightness = 0.65
		fadeDistance   = 12.0
	)

	t := max(0.0, min(distance/fadeDistance, 1.0))
	fade := 1 - t*t*(3-2*t)
	return float32(nearBrightness * fade)
}

// surfaceBrightness keeps surfaces quieter than walls.
func surfaceBrightness(x, y, distance float64, ceiling bool) float32 {
	brightness := worldBrightness(distance)

	if ceiling {
		return brightness * 0.25
	}

	// Broad, gentle variations provide movement cues without hard tile edges.
	pattern := 0.9 + 0.1*math.Cos(x*math.Pi/2)*math.Cos(y*math.Pi/2)
	return brightness * 0.5 * float32(pattern)
}
