package game

import (
	"math"
	"testing"
)

func TestRoomPositionCoordinates(t *testing.T) {
	arena := NewArena(42)

	cases := []struct {
		worldX, worldY float64
		roomX, roomY   int
		localX, localY float64
	}{
		{12.5, 12.5, 0, 0, 12.5, 12.5},
		{24, 24, 1, 1, 0, 0},
		{-0.5, -0.5, -1, -1, 23.5, 23.5},
		{-24, -24, -1, -1, 0, 0},
		{-24.5, 48.5, -2, 2, 23.5, 0.5},
	}

	for _, test := range cases {
		position := arena.RoomPositionAt(test.worldX, test.worldY)
		expectedID := RoomID{
			Seed: 42,
			X:    test.roomX,
			Y:    test.roomY,
		}

		if position.Room != expectedID {
			t.Fatalf(
				"world %v,%v: expected room %v, got %v",
				test.worldX, test.worldY, expectedID, position.Room,
			)
		}

		if position.X != test.localX || position.Y != test.localY {
			t.Fatalf(
				"world %v,%v: expected local %v,%v, got %v,%v",
				test.worldX, test.worldY,
				test.localX, test.localY,
				position.X, position.Y,
			)
		}

		x, y := position.WorldPosition()
		if math.Abs(x-test.worldX) > 1e-9 ||
			math.Abs(y-test.worldY) > 1e-9 {
			t.Fatalf("world position did not round-trip: got %v,%v", x, y)
		}
	}
}

func TestRoomReferenceSurvivesUnloading(t *testing.T) {
	arena := NewArena(42)
	saved := arena.CurrentRoomPosition()
	savedX, savedY := arena.Camera.X, arena.Camera.Y
	origin := chunkPosition{0, 0}

	if arena.chunks[origin].RoomID != saved.Room {
		t.Fatal("generated room identity must match the player's reference")
	}

	arena.Camera.X += float64(chunkSize * 5)
	arena.loadNearbyChunks()

	if _, loaded := arena.chunks[origin]; loaded {
		t.Fatal("original room should now be unloaded")
	}

	x, y := saved.WorldPosition()
	if x != savedX || y != savedY {
		t.Fatal("saved position changed after leaving the room")
	}

	arena.Camera.X = x
	arena.Camera.Y = y
	arena.loadNearbyChunks()

	if arena.CurrentRoomPosition() != saved {
		t.Fatal("returning must restore the same room-relative position")
	}
	if arena.chunks[origin].RoomID != saved.Room {
		t.Fatal("regenerated room must retain its identity")
	}
}

func TestRoomReferenceDoesNotLoadTerrain(t *testing.T) {
	arena := NewArena(42)
	loadedBefore := len(arena.chunks)

	position := arena.RoomPositionAt(
		float64(chunkSize*100)+0.5,
		float64(-chunkSize*100)+0.5,
	)

	if position.Room != (RoomID{Seed: 42, X: 100, Y: -100}) {
		t.Fatalf("unexpected distant room: %v", position.Room)
	}
	if len(arena.chunks) != loadedBefore {
		t.Fatal("looking up a room reference must not load terrain")
	}
	if _, loaded := arena.chunks[chunkPosition{100, -100}]; loaded {
		t.Fatal("distant room should remain unloaded")
	}
}

func TestRoomReferenceOffset(t *testing.T) {
	anchor := RoomID{Seed: 42, X: -2, Y: 3}
	target := anchor.Offset(3, -1)

	if target != (RoomID{Seed: 42, X: 1, Y: 2}) {
		t.Fatalf("unexpected target room: %v", target)
	}
	if anchor != (RoomID{Seed: 42, X: -2, Y: 3}) {
		t.Fatal("calculating a target must not change the anchor")
	}
}

func TestRoomIdentityIncludesWorldSeed(t *testing.T) {
	first := NewArena(42)
	second := NewArena(43)

	firstID := first.CurrentRoomPosition().Room
	secondID := second.CurrentRoomPosition().Room

	if firstID.X != secondID.X || firstID.Y != secondID.Y {
		t.Fatal("test requires matching room-grid coordinates")
	}
	if firstID == secondID {
		t.Fatal("rooms in different seeded worlds must have different identities")
	}
}
