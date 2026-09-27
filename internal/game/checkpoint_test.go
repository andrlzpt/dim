package game

import "testing"

func TestCheckpointInUnloadedRoom(t *testing.T) {
	arena := NewArena(42)
	anchor := arena.CurrentRoomPosition()
	loadedBefore := len(arena.chunks)

	checkpoint, err := arena.CreateCheckpoint("test", 3, -1)
	if err != nil {
		t.Fatal(err)
	}

	if checkpoint.Anchor != anchor {
		t.Fatal("checkpoint must capture its creation position")
	}
	if checkpoint.Destination.Room != anchor.Room.Offset(3, -1) {
		t.Fatal("checkpoint must use the requested relative room")
	}
	if checkpoint.Status != CheckpointUndiscovered {
		t.Fatal("new checkpoint must start undiscovered")
	}
	if len(arena.chunks) != loadedBefore {
		t.Fatal("checkpoint creation must not expand the loaded cache")
	}

	target := checkpoint.Destination.Room
	if _, loaded := arena.chunks[chunkPosition{target.X, target.Y}]; loaded {
		t.Fatal("distant destination must remain unloaded")
	}

	// Load the destination and verify the stored location against real terrain.
	x, y := checkpoint.Destination.WorldPosition()
	arena.Camera.X = x
	arena.Camera.Y = y
	arena.loadNearbyChunks()

	if !arena.canStand(x, y, radius) {
		t.Fatal("checkpoint destination must be safe")
	}

	resolved, found := arena.FindRoomPlacement(target, radius)
	if !found || resolved != checkpoint.Destination {
		t.Fatal("loaded terrain must agree with preview placement")
	}
}

func TestCheckpointKeepsItsAnchor(t *testing.T) {
	arena := NewArena(42)

	first, err := arena.CreateCheckpoint("test", 3, -1)
	if err != nil {
		t.Fatal(err)
	}

	arena.Camera.X += float64(chunkSize * 10)
	arena.Camera.Y -= float64(chunkSize * 10)
	arena.loadNearbyChunks()

	// Repeating the narrative event must not move the existing checkpoint.
	again, err := arena.CreateCheckpoint("test", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if again != first {
		t.Fatal("existing checkpoint must retain its original record")
	}

	stored, found := arena.Checkpoint("test")
	if !found || stored != first {
		t.Fatal("checkpoint must survive unloading its surrounding terrain")
	}
}

func TestCheckpointStatus(t *testing.T) {
	arena := NewArena(42)
	checkpoint, err := arena.CreateCheckpoint("test", 3, -1)
	if err != nil {
		t.Fatal(err)
	}

	x, y := checkpoint.Destination.WorldPosition()

	// Enter the destination region, but remain outside the reach distance.
	arena.Camera.X = x + 2
	arena.Camera.Y = y
	arena.loadNearbyChunks()
	arena.updateCheckpoints()

	stored, _ := arena.Checkpoint("test")
	if stored.Status != CheckpointDiscovered {
		t.Fatal("entering the destination room must discover the checkpoint")
	}

	// Use the normal update path to verify arrival detection is connected.
	arena.Camera.X = x
	arena.Camera.Y = y
	arena.Update(0, 0, 0, 1.0/60.0)

	stored, _ = arena.Checkpoint("test")
	if stored.Status != CheckpointReached {
		t.Fatal("arriving at the destination must reach the checkpoint")
	}

	arena.Camera.X += float64(chunkSize * 10)
	arena.Update(0, 0, 0, 1.0/60.0)

	stored, _ = arena.Checkpoint("test")
	if stored.Status != CheckpointReached {
		t.Fatal("leaving and unloading must not reset reached status")
	}
}

func TestCheckpointRejectsEmptyID(t *testing.T) {
	arena := NewArena(42)

	if _, err := arena.CreateCheckpoint("", 3, -1); err == nil {
		t.Fatal("empty checkpoint IDs must be rejected")
	}
	if len(arena.checkpoints) != 0 {
		t.Fatal("failed creation must not store a checkpoint")
	}
}
