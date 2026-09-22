package main

import (
	"math"

	"github.com/andrlzpt/dim/internal/ascii"
	"github.com/andrlzpt/dim/internal/renderer"
)

const (
	nearBrightness    = 0.65
	worldViewDistance = 12.0
	horizontalFOV     = math.Pi / 3
	wallHeight        = 3.0
)

type surfaceRow struct {
	Depth   float64
	Ceiling bool
}

func drawFrame(world, text *ascii.Buffer, app *AppState) {
	world.Clear()
	text.Clear()
	if app.Screen == ScreenArena {
		drawArena(world, app)
	}
	drawAppMenu(text, app)

	renderer.BeginFrame(app.FramebufferWidth, app.FramebufferHeight)
	renderer.RenderText(world, renderer.Config{
		WindowWidth: app.Width, WindowHeight: app.Height,
		CellWidth: worldCellWidth, CellHeight: worldCellHeight,
		GlyphWidth: worldGlyphWidth, GlyphHeight: worldGlyphHeight,
	})
	renderer.RenderText(text, renderer.Config{
		WindowWidth: app.Width, WindowHeight: app.Height,
		CellWidth: textCellWidth, CellHeight: textCellHeight,
	})
}

// drawArena projects the room into the world character buffer.
func drawArena(world *ascii.Buffer, app *AppState) {
	camera := app.Arena.Camera
	forwardX := math.Cos(camera.Angle)
	forwardY := math.Sin(camera.Angle)
	planeScale := math.Tan(horizontalFOV / 2)

	planeX := -forwardY * planeScale
	planeY := forwardX * planeScale

	projection := float64(world.Width*worldCellWidth) /
		(2 * planeScale * float64(worldCellHeight))

	// Looking up moves the horizon down, revealing more ceiling.
	horizon := float64(world.Height) * (0.5 + camera.Tilt)

	// Reuse storage unless the character grid changes height.
	if len(app.surfaceRows) != world.Height {
		app.surfaceRows = make([]surfaceRow, world.Height)
	}

	for y := range app.surfaceRows {
		offset := float64(y) + 0.5 - horizon
		row := surfaceRow{
			Ceiling: offset < 0,
		}

		if math.Abs(offset) >= 0.001 {
			row.Depth = (wallHeight / 2) * projection / math.Abs(offset)
		}

		app.surfaceRows[y] = row
	}

	for x := 0; x < world.Width; x++ {
		screenX := 2*(float64(x)+0.5)/float64(world.Width) - 1
		rayX := forwardX + planeX*screenX
		rayY := forwardY + planeY*screenX
		rayLength := math.Hypot(rayX, rayY)

		// Find and draw the wall before calculating other surfaces.
		wallStart, wallEnd := 0, 0
		depth, side, hit := app.Arena.CastWall(
			rayX,
			rayY,
			worldViewDistance,
		)

		if hit {
			depth = max(depth, 0.001)

			height := wallHeight * projection / depth
			top := horizon - height/2
			bottom := horizon + height/2

			wallStart = max(0, min(world.Height, int(math.Ceil(top-0.5))))
			wallEnd = max(0, min(world.Height, int(math.Ceil(bottom-0.5))))

			brightness := worldBrightness(depth * rayLength)
			if side {
				brightness *= 0.75
			}

			for y := wallStart; y < wallEnd; y++ {
				world.SetShaded(x, y, '+', brightness)
			}
		}

		for y := 0; y < world.Height; y++ {
			// Jump past rows already covered by the wall.
			if y >= wallStart && y < wallEnd {
				y = wallEnd - 1
				continue
			}

			row := app.surfaceRows[y]
			if row.Depth == 0 {
				continue
			}

			// Ray length varies by column, so distance still does too.
			distance := row.Depth * rayLength
			if distance >= worldViewDistance {
				continue
			}

			surfaceX := camera.X + rayX*row.Depth
			surfaceY := camera.Y + rayY*row.Depth

			brightness := surfaceBrightness(
				surfaceX,
				surfaceY,
				distance,
				row.Ceiling,
			)
			world.SetShaded(x, y, '.', brightness)
		}
	}
}

// worldBrightness fades soft gray smoothly toward black.
func worldBrightness(distance float64) float32 {

	t := max(0.0, min(distance/worldViewDistance, 1.0))
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
