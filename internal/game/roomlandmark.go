package game

import "math/rand"

type LandmarkKind int

const (
	LandmarkNone LandmarkKind = iota
	LandmarkPillarPair
	LandmarkFourPillars
)

func generateRoomLandmark(
	seed int64,
	position chunkPosition,
	grid [][]byte,
	bounds room,
) LandmarkKind {

	// Landmarks randomness is independent of room geometry
	random := rand.New(rand.NewSource(
		chunkSeed(seed, position.X, position.Y, 3),
	))

	kind := LandmarkKind(random.Intn(3))

	// Leave walkable border around the pillars and an open center
	left := bounds.X + 1
	right := bounds.X + bounds.Width - 2
	top := bounds.Y + 1
	bottom := bounds.Y + bounds.Height - 2

	switch kind {
	case LandmarkPillarPair:
		grid[top][left] = '#'
		grid[bottom][right] = '#'

	case LandmarkFourPillars:
		grid[top][left] = '#'
		grid[top][right] = '#'
		grid[bottom][left] = '#'
		grid[bottom][right] = '#'
	}

	return kind
}
