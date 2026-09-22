package game

import (
	"reflect"
	"testing"
)

func TestChunkCoordinates(t *testing.T) {
	cases := []struct {
		world, chunk, local int
	}{
		{-25, -2, 23},
		{-24, -1, 0},
		{-1, -1, 23},
		{0, 0, 0},
		{23, 0, 23},
		{24, 1, 0},
	}

	for _, test := range cases {
		index, local := chunkCoordinate(test.world)
		if index != test.chunk || local != test.local {
			t.Fatalf(
				"world %d: got chunk %d, local %d",
				test.world, index, local,
			)
		}
	}
}

func TestChunkStreamingReload(t *testing.T) {
	arena := NewArena(42)
	origin := chunkPosition{0, 0}
	original := arena.chunks[origin]

	if !arena.canStand(arena.Camera.X, arena.Camera.Y, radius) {
		t.Fatal("initial spawn must be safe")
	}

	for _, position := range []chunkPosition{
		{4, 3}, {-4, -3}, {0, 0},
	} {
		arena.Camera.X = float64(position.X*chunkSize) + 12.5
		arena.Camera.Y = float64(position.Y*chunkSize) + 12.5
		arena.loadNearbyChunks()

		if len(arena.chunks) != 9 {
			t.Fatalf("expected nine loaded chunks, got %d", len(arena.chunks))
		}
		if !arena.canStand(arena.Camera.X, arena.Camera.Y, radius) {
			t.Fatalf("chunk %v: center must be safe", position)
		}

		if position != origin {
			if _, loaded := arena.chunks[origin]; loaded {
				t.Fatal("distant origin chunk was not unloaded")
			}
		}
	}

	if !reflect.DeepEqual(original, arena.chunks[origin]) {
		t.Fatal("returning must reconstruct the same terrain")
	}
}

func TestMovementAcrossChunkBoundary(t *testing.T) {
	for _, chunkX := range []int{-1, 0} {
		arena := NewArena(42)
		boundaryX := (chunkX + 1) * chunkSize
		openingY := boundaryOpening(42, chunkX+1, 0, 1)

		arena.Camera.X = float64(boundaryX) - 0.1
		arena.Camera.Y = float64(openingY) + 0.5
		arena.loadNearbyChunks()

		// Face east and walk across the shared opening.
		arena.Update(0, -1, 0, 0.1)

		if arena.Camera.X <= float64(boundaryX) {
			t.Fatalf("could not cross from chunk %d", chunkX)
		}
		if arena.loadedCenter != (chunkPosition{chunkX + 1, 0}) {
			t.Fatal("movement did not update the loaded neighborhood")
		}
	}
}

func TestRayDistanceLimit(t *testing.T) {
	arena := NewArena(42)

	// The spawn is well inside a room: a short ray cannot hit a wall.
	depth, _, hit := arena.CastWall(2, 0, 1)
	if hit || depth != 0.5 {
		t.Fatalf("expected no hit at depth 0.5, got depth %v, hit %v", depth, hit)
	}
}
