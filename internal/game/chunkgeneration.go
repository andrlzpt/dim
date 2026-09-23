package game

import (
	"encoding/binary"
	"hash/fnv"
	"math/rand"
)

const chunkSize = 24

// chunkPosition identifies a chunk, including west/north of the origin.
type chunkPosition struct {
	X, Y int
}

type chunk struct {
	RoomID RoomID
	Tiles  []string
}

func generateChunk(seed int64, position chunkPosition) chunk {
	random := rand.New(rand.NewSource(
		chunkSeed(seed, position.X, position.Y, 0),
	))

	grid := make([][]byte, chunkSize)
	for y := range grid {
		grid[y] = make([]byte, chunkSize)
		for x := range grid[y] {
			grid[y][x] = '#'
		}
	}

	// Every room contains the chunk center.
	width := 8 + random.Intn(5)
	height := 8 + random.Intn(5)
	x := (chunkSize-width)/2 + random.Intn(3) - 1
	y := (chunkSize-height)/2 + random.Intn(3) - 1
	carveArea(grid, x, y, width, height)

	center := chunkSize / 2

	// Vertical boundaries use the coordinate of the chunk to their east.
	west := boundaryOpening(seed, position.X, position.Y, 1)
	east := boundaryOpening(seed, position.X+1, position.Y, 1)

	// Horizontal boundaries use the coordinate of the chunk to their south.
	north := boundaryOpening(seed, position.X, position.Y, 2)
	south := boundaryOpening(seed, position.X, position.Y+1, 2)

	// Join each two-cell opening to the central room.
	carveArea(grid, 0, west, center+2, 2)
	carveVerticalConnection(grid, center, west, center)

	carveArea(grid, center, east, chunkSize-center, 2)
	carveVerticalConnection(grid, center, east, center)

	carveArea(grid, north, 0, 2, center+2)
	carveHorizontalConnection(grid, north, center, center)

	carveArea(grid, south, center, 2, chunkSize-center)
	carveHorizontalConnection(grid, south, center, center)

	tiles := make([]string, chunkSize)
	for y := range grid {
		tiles[y] = string(grid[y])
	}

	return chunk{
		RoomID: RoomID{
			Seed: seed,
			X:    position.X,
		},
		Tiles: tiles,
	}
}

func carveVerticalConnection(grid [][]byte, x, fromY, toY int) {
	carveArea(
		grid,
		x,
		min(fromY, toY),
		2,
		max(fromY, toY)-min(fromY, toY)+2,
	)
}

func carveHorizontalConnection(grid [][]byte, fromX, toX, y int) {
	carveArea(
		grid,
		min(fromX, toX),
		y,
		max(fromX, toX)-min(fromX, toX)+2,
		2,
	)
}

// boundaryOpening leaves solid space around both ends of an edge.
func boundaryOpening(seed int64, x, y int, kind uint64) int {
	random := rand.New(rand.NewSource(chunkSeed(seed, x, y, kind)))
	return 2 + random.Intn(chunkSize-5)
}

// chunkSeed derives repeatable randomness from coordinates and purpose.
// kind separates room layout, vertical edges, and horizontal edges.
func chunkSeed(seed int64, x, y int, kind uint64) int64 {
	var data [32]byte
	binary.LittleEndian.PutUint64(data[0:8], uint64(seed))
	binary.LittleEndian.PutUint64(data[8:16], uint64(int64(x)))
	binary.LittleEndian.PutUint64(data[16:24], uint64(int64(y)))
	binary.LittleEndian.PutUint64(data[24:32], kind)

	hash := fnv.New64a()
	_, _ = hash.Write(data[:])
	return int64(hash.Sum64())
}
