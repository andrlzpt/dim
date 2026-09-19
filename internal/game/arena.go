package game

import "math"

// Camera holds a position in world units and a heading in radians.
// X increases east, Y increases south, and angle zero faces east.
type Camera struct {
	X     float64
	Y     float64
	Angle float64
	Tilt  float64 // Fraction of screen height; positive looks up.

}

// Arena holds the room layout and the player's viewpoint.
type Arena struct {
	Camera Camera

	seed         int64
	chunks       map[chunkPosition]chunk
	loadedCenter chunkPosition
}

func NewArena(seed int64) Arena {
	arena := Arena{
		Camera: Camera{
			X: float64(chunkSize/2) + 0.5,
			Y: float64(chunkSize/2) + 0.5,
		},
		seed:   seed,
		chunks: make(map[chunkPosition]chunk),
	}
	arena.loadNearbyChunks()

	return arena
}

// isWall reads loaded terrain without generating more chunks.
// Unloaded terrain is solid as a safety boundary.
func (arena *Arena) isWall(x, y int) bool {
	chunkX, localX := chunkCoordinate(x)
	chunkY, localY := chunkCoordinate(y)

	current, loaded := arena.chunks[chunkPosition{chunkX, chunkY}]
	if !loaded {
		return true
	}
	return current.Tiles[localY][localX] == '#'
}

func (arena *Arena) loadNearbyChunks() {
	chunkX, _ := chunkCoordinate(int(math.Floor(arena.Camera.X)))
	chunkY, _ := chunkCoordinate(int(math.Floor(arena.Camera.Y)))
	center := chunkPosition{chunkX, chunkY}

	if center == arena.loadedCenter && len(arena.chunks) == 9 {
		return
	}

	// Unload first, keeping at most nine chunks in the cache.
	for position := range arena.chunks {
		if position.X < chunkX-1 || position.X > chunkX+1 ||
			position.Y < chunkY-1 || position.Y > chunkY+1 {
			delete(arena.chunks, position)
		}
	}

	for y := chunkY - 1; y <= chunkY+1; y++ {
		for x := chunkX - 1; x <= chunkX+1; x++ {
			position := chunkPosition{x, y}
			if _, loaded := arena.chunks[position]; !loaded {
				arena.chunks[position] = generateChunk(arena.seed, position)
			}
		}
	}

	arena.loadedCenter = center
}

// chunkCoordinate splits a world tile coordinate into chunk and local indices.
// For example, -1 belongs to chunk -1, local tile 23.
func chunkCoordinate(value int) (int, int) {
	index := value / chunkSize
	local := value % chunkSize
	if local < 0 {
		index--
		local += chunkSize
	}
	return index, local
}
