package game

import "math"

func (arena *Arena) FindRoomPlacement(
	id RoomID,
	clearance float64,
) (RoomPosition, bool) {
	if id.Seed != arena.seed {
		return RoomPosition{}, false
	}

	current, loaded := arena.chunks[chunkPosition{id.X, id.Y}]
	if !loaded {
		return RoomPosition{}, false
	}

	bounds := current.Bounds
	maxClearance := float64(min(bounds.Width, bounds.Height)) / 2
	if math.IsNaN(clearance) || clearance <= 0 || clearance > maxClearance {
		return RoomPosition{}, false
	}

	centerX := float64(bounds.X) + float64(bounds.Width)/2
	centerY := float64(bounds.Y) + float64(bounds.Height)/2

	var best RoomPosition
	found := false
	bestDistance := math.Inf(1)

	for y := bounds.Y; y < bounds.Y+bounds.Height; y++ {
		for x := bounds.X; x < bounds.X+bounds.Width; x++ {
			position := RoomPosition{
				Room: id,
				X:    float64(x) + 0.5,
				Y:    float64(y) + 0.5,
			}

			if position.X-clearance < float64(bounds.X) ||
				position.Y-clearance < float64(bounds.Y) ||
				position.X+clearance > float64(bounds.X+bounds.Width) ||
				position.Y+clearance > float64(bounds.Y+bounds.Height) {
				continue
			}

			worldX, worldY := position.WorldPosition()
			if !arena.canStand(worldX, worldY, clearance) {
				continue
			}

			dx := position.X - centerX
			dy := position.Y - centerY

			distance := dx*dx + dy*dy

			if distance < bestDistance {
				best = position
				bestDistance = distance
				found = true
			}
		}
	}

	return best, found
}
