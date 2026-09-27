package game

import "testing"

func TestStartingRoomDiscovered(t *testing.T) {
	arena := NewArena(42)
	id := arena.CurrentRoomPosition().Room

	discovery, found := arena.DiscoveredRoom(id)
	if !found {
		t.Fatal("starting room must be discovered")
	}

	current := arena.chunks[chunkPosition{id.X, id.Y}]
	if discovery.Landmark != current.Landmark {
		t.Fatal("discovery must remember the room's landmark")
	}

	// Neighbors are loaded for rendering, but haven't been visited.
	if _, found := arena.DiscoveredRoom(id.Offset(1, 0)); found {
		t.Fatal("loading a neighboring room must not discover it")
	}
	if len(arena.discoveredRooms) != 1 {
		t.Fatal("only the starting room should be discovered")
	}
}

func TestDiscoverySurvivesUnloading(t *testing.T) {
	arena := NewArena(42)
	origin := arena.CurrentRoomPosition().Room
	original, _ := arena.DiscoveredRoom(origin)

	// Move to a distant room center, including negative coordinates.
	arena.Camera.X = float64(-4*chunkSize) + 12.5
	arena.Camera.Y = float64(3*chunkSize) + 12.5
	arena.loadNearbyChunks()
	arena.Update(0, 0, 0, 1.0/60.0)

	destination := arena.CurrentRoomPosition().Room
	if _, found := arena.DiscoveredRoom(destination); !found {
		t.Fatal("the normal update must discover the current room")
	}

	if _, loaded := arena.chunks[chunkPosition{origin.X, origin.Y}]; loaded {
		t.Fatal("original terrain should be unloaded")
	}

	remembered, found := arena.DiscoveredRoom(origin)
	if !found || remembered != original {
		t.Fatal("discovery must survive terrain unloading")
	}

	// Repeated updates must not create duplicate records.
	arena.Update(0, 0, 0, 1.0/60.0)
	if len(arena.discoveredRooms) != 2 {
		t.Fatal("expected exactly two visited rooms")
	}
}

func TestCheckpointPreviewDoesNotDiscoverRoom(t *testing.T) {
	arena := NewArena(42)

	checkpoint, err := arena.CreateCheckpoint("test", 3, -1)
	if err != nil {
		t.Fatal(err)
	}

	if _, found := arena.DiscoveredRoom(checkpoint.Destination.Room); found {
		t.Fatal("previewing checkpoint terrain must not discover its room")
	}
	if len(arena.discoveredRooms) != 1 {
		t.Fatal("checkpoint creation must preserve the explored map")
	}
}
