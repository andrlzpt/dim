package main

import "math"

// Camera holds a position in world units and a heading in radians.
// X increases east, Y increases south, and angle zero faces east.
type Camera struct {
	X     float64
	Y     float64
	Angle float64
}

// ArenaState holds the room layout and the player's viewpoint.
type ArenaState struct {
	Tiles  []string
	Camera Camera
}

func newArena() ArenaState {
	return ArenaState{
		Tiles: []string{
			"##########",
			"#........#",
			"#........#",
			"#........#",
			"#........#",
			"#........#",
			"#........#",
			"##########",
		},
		Camera: Camera{
			X:     2.5,
			Y:     4.0,
			Angle: 0,
		},
	}
}

// isWall treats positions outside the room as solid.
func (arena *ArenaState) isWall(x, y int) bool {
	if y < 0 || y >= len(arena.Tiles) {
		return true
	}
	if x < 0 || x >= len(arena.Tiles[y]) {
		return true
	}
	return arena.Tiles[y][x] == '#'
}

// castWall steps through grid cells until it reaches a wall.
// The returned distance is perpendicular camera depth when the
// ray comes from the camera-plane construction in drawArena.
// side reports whether the ray crossed a horizontal grid boundary.
func castWall(arena *ArenaState, rayX, rayY float64) (float64, bool) {
	camera := arena.Camera
	cellX := int(math.Floor(camera.X))
	cellY := int(math.Floor(camera.Y))

	stepX, stepY := 1, 1
	deltaX, deltaY := math.Inf(1), math.Inf(1)
	nextX, nextY := math.Inf(1), math.Inf(1)

	if rayX < 0 {
		stepX = -1
		deltaX = -1 / rayX
		nextX = (camera.X - float64(cellX)) * deltaX
	} else if rayX > 0 {
		deltaX = 1 / rayX
		nextX = (float64(cellX+1) - camera.X) * deltaX
	}

	if rayY < 0 {
		stepY = -1
		deltaY = -1 / rayY
		nextY = (camera.Y - float64(cellY)) * deltaY
	} else if rayY > 0 {
		deltaY = 1 / rayY
		nextY = (float64(cellY+1) - camera.Y) * deltaY
	}

	// Also makes calls with an empty arena or zero ray safe.
	if arena.isWall(cellX, cellY) || (rayX == 0 && rayY == 0) {
		return 0, false
	}

	for {
		var distance float64
		var side bool

		if nextX < nextY {
			distance = nextX
			nextX += deltaX
			cellX += stepX
		} else {
			distance = nextY
			nextY += deltaY
			cellY += stepY
			side = true
		}

		if arena.isWall(cellX, cellY) {
			return distance, side
		}
	}
}
