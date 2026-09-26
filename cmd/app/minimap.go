package main

import (
	"math"

	"github.com/andrlzpt/dim/internal/ascii"
	"github.com/andrlzpt/dim/internal/game"
	"github.com/andrlzpt/dim/internal/renderer"
)

type MinimapCorner int

const (
	MinimapBottomLeft MinimapCorner = iota
	MinimapBottomRight
	MinimapTopLeft
	MinimapTopRight
)

type MinimapConfig struct {
	Enabled      bool
	Corner       MinimapCorner
	Radius       int
	Margin       int
	RoomSpacing  int
	CheckpointID string
}

func defaultMinimapConfig() MinimapConfig {
	return MinimapConfig{
		Enabled:      true,
		Corner:       MinimapBottomLeft,
		Radius:       96,
		Margin:       16,
		RoomSpacing:  32,
		CheckpointID: "first-checkpoint",
	}
}

func drawMinimap(app *AppState) {
	const (
		cellWidth  = 7
		cellHeight = 13
	)

	settings := app.Minimap
	if !settings.Enabled {
		return
	}

	margin := max(0, settings.Margin)
	radius := min(
		settings.Radius,
		(min(app.Width, app.Height)-2*margin)/2,
	)
	if radius < 40 {
		return
	}

	diameter := radius * 2
	left, top := minimapPosition(
		settings.Corner,
		app.Width,
		app.Height,
		diameter,
		margin,
	)

	columns := (diameter + cellWidth - 1) / cellWidth
	rows := (diameter + cellHeight - 1) / cellHeight

	if app.minimapText == nil ||
		app.minimapText.Width != columns ||
		app.minimapText.Height != rows {
		app.minimapText = ascii.NewBuffer(columns, rows)
	}
	buffer := app.minimapText
	buffer.Clear()

	center := float64(radius)
	innerRadius := center - 3

	// Coordinates here are relative to the minimap's top-left corner.
	// Check the entire glyph rectangle against the circular boundary.
	put := func(x, y float64, glyph rune, brightness float32) {
		column := int(math.Floor(x / cellWidth))
		row := int(math.Floor(y / cellHeight))

		glyphCenterX := (float64(column) + 0.5) * cellWidth
		glyphCenterY := (float64(row) + 0.5) * cellHeight

		dx := math.Abs(glyphCenterX-center) + float64(cellWidth)/2
		dy := math.Abs(glyphCenterY-center) + float64(cellHeight)/2
		if dx*dx+dy*dy > innerRadius*innerRadius {
			return
		}

		buffer.SetShaded(column, row, glyph, brightness)
	}

	playerRoom := app.Arena.CurrentRoomPosition().Room
	spacing := float64(max(settings.RoomSpacing, cellHeight*2))
	extent := int(math.Ceil(center / spacing))

	// Draw connections first. Every cardinal neighbor is connected
	// by the current generator; both rooms must have been visited.
	for dy := -extent; dy <= extent; dy++ {
		for dx := -extent; dx <= extent; dx++ {
			id := playerRoom.Offset(dx, dy)
			if _, found := app.Arena.DiscoveredRoom(id); !found {
				continue
			}

			x := center + float64(dx)*spacing
			y := center + float64(dy)*spacing

			if _, found := app.Arena.DiscoveredRoom(id.Offset(1, 0)); found {
				put(x+spacing/2, y, '-', 0.35)
			}
			if _, found := app.Arena.DiscoveredRoom(id.Offset(0, 1)); found {
				put(x, y+spacing/2, '|', 0.35)
			}
		}
	}

	// Draw remembered room landmarks over the connections.
	for dy := -extent; dy <= extent; dy++ {
		for dx := -extent; dx <= extent; dx++ {
			id := playerRoom.Offset(dx, dy)
			discovery, found := app.Arena.DiscoveredRoom(id)
			if !found {
				continue
			}

			put(
				center+float64(dx)*spacing,
				center+float64(dy)*spacing,
				minimapRoomGlyph(discovery.Landmark),
				0.65,
			)
		}
	}

	checkpoint, found := app.Arena.Checkpoint(settings.CheckpointID)
	if found {
		targetRoom := checkpoint.Destination.Room

		// Show the destination room only after it has been visited.
		if _, discovered := app.Arena.DiscoveredRoom(targetRoom); discovered {
			glyph := 'X'
			if checkpoint.Status == game.CheckpointReached {
				glyph = 'V'
			}
			put(
				center+float64(targetRoom.X-playerRoom.X)*spacing,
				center+float64(targetRoom.Y-playerRoom.Y)*spacing,
				glyph,
				0.9,
			)
		}

		// The rim marker gives a bearing, not a path through the corridors.
		if checkpoint.Status != game.CheckpointReached {
			targetX, targetY := checkpoint.Destination.WorldPosition()
			dx := targetX - app.Arena.Camera.X
			dy := targetY - app.Arena.Camera.Y
			distance := math.Hypot(dx, dy)

			if distance > 0 {
				markerRadius := center - 17
				put(
					center+dx/distance*markerRadius,
					center+dy/distance*markerRadius,
					'*',
					0.95,
				)
			}
		}
	}

	// North stays at the top; the map does not rotate with the camera.
	put(center, 14, 'N', 0.7)

	angle := app.Arena.Camera.Angle
	put(
		center+math.Cos(angle)*17,
		center+math.Sin(angle)*17,
		minimapHeading(angle),
		0.9,
	)
	put(center, center, '@', 1)

	config := renderer.Config{
		WindowWidth:  app.Width,
		WindowHeight: app.Height,
		CellWidth:    cellWidth,
		CellHeight:   cellHeight,
	}

	renderer.RenderCircle(
		float32(left+radius),
		float32(top+radius),
		float32(radius),
		config,
	)
	renderer.RenderTextAt(buffer, float32(left), float32(top), config)
}

func minimapPosition(
	corner MinimapCorner,
	width, height, diameter, margin int,
) (int, int) {
	left := margin
	top := height - margin - diameter

	switch corner {
	case MinimapBottomRight:
		left = width - margin - diameter
	case MinimapTopLeft:
		top = margin
	case MinimapTopRight:
		left = width - margin - diameter
		top = margin
	}

	return left, top
}

func minimapRoomGlyph(landmark game.LandmarkKind) rune {
	switch landmark {
	case game.LandmarkPillarPair:
		return '2'
	case game.LandmarkFourPillars:
		return '4'
	default:
		return 'o'
	}
}

func minimapHeading(angle float64) rune {
	direction := int(math.Round(angle / (math.Pi / 2)))
	direction = (direction%4 + 4) % 4
	return []rune{'>', 'v', '<', '^'}[direction]
}
