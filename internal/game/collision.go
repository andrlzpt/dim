package game

import "math"

// canStand checks a circular player against nearby wall cells.
func (arena *Arena) canStand(x, y, radius float64) bool {
	minX := int(math.Floor(x - radius))
	maxX := int(math.Floor(x + radius))
	minY := int(math.Floor(y - radius))
	maxY := int(math.Floor(y + radius))

	for cellY := minY; cellY <= maxY; cellY++ {
		for cellX := minX; cellX <= maxX; cellX++ {
			if !arena.isWall(cellX, cellY) {
				continue
			}

			nearestX := max(float64(cellX), min(x, float64(cellX+1)))
			nearestY := max(float64(cellY), min(y, float64(cellY+1)))
			dx := x - nearestX
			dy := y - nearestY

			if dx*dx+dy*dy < radius*radius {
				return false
			}
		}
	}
	return true
}
