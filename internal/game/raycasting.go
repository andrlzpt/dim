package game

import "math"

// CastWall steps through grid cells from the camera until it hits a wall.
// It returns t for the intersection Camera + t*ray.
// With camera-plane rays, t is perpendicular camera depth.
// The boolean reports crossing a horizontal grid boundary.
func (arena *Arena) CastWall(rayX, rayY float64) (float64, bool) {
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
