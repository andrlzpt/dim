package game

import (
	"fmt"
	"math"
)

type CheckpointStatus int

const (
	CheckpointUndiscovered CheckpointStatus = iota
	CheckpointDiscovered
	CheckpointReached
)

// Checkpoint stores narrative state independently of loaded terrain
type Checkpoint struct {
	ID          string
	Anchor      RoomPosition
	Destination RoomPosition
	Status      CheckpointStatus
}

// CreateCheckpoint captures the player's current position as its anchor
// Room offsets are east/west and south/north grid displacements
// Reusing an ID returns the existing checkpoint without moving or resetting it
func (arena *Arena) CreateCheckpoint(
	id string,
	offsetX, offsetY int,
) (Checkpoint, error) {
	if id == "" {
		return Checkpoint{}, fmt.Errorf("checkpoint ID must not be empty")
	}

	if existing, found := arena.checkpoints[id]; found {
		return existing, nil
	}

	anchor := arena.CurrentRoomPosition()
	targetRoom := anchor.Room.Offset(offsetX, offsetY)

	destination, found := arena.checkpointPlacement(targetRoom)
	if !found {
		return Checkpoint{}, fmt.Errorf(
			"no safe placement for checkpoint %q in room %v",
			id,
			targetRoom,
		)
	}

	checkpoint := Checkpoint{
		ID:          id,
		Anchor:      anchor,
		Destination: destination,
		Status:      CheckpointUndiscovered,
	}

	if arena.checkpoints == nil {
		arena.checkpoints = make(map[string]Checkpoint)
	}
	arena.checkpoints[id] = checkpoint

	return checkpoint, nil
}

// Checkpoint returns a copy of the stored record
func (arena *Arena) Checkpoint(id string) (Checkpoint, bool) {
	checkpoint, found := arena.checkpoints[id]
	return checkpoint, found
}

// checkpointPlacemente can inspect unloaded terrain without changing the cache
func (arena *Arena) checkpointPlacement(id RoomID) (RoomPosition, bool) {
	position := chunkPosition{X: id.X, Y: id.Y}

	if _, loaded := arena.chunks[position]; loaded {
		return arena.FindRoomPlacement(id, radius)
	}

	// A temporary arena lets us reuse placement and collision checks
	// Placement keeps the footprint inside this room so neighbos aren't needed
	preview := Arena{
		seed: arena.seed,
		chunks: map[chunkPosition]chunk{
			position: generateChunk(arena.seed, position),
		},
	}

	return preview.FindRoomPlacement(id, radius)
}

// updateCheckpoints advances status without moving stored destinations
func (arena *Arena) updateCheckpoints() {
	const reachDistance = 0.75

	playerRoom := arena.CurrentRoomPosition().Room

	for id, checkpoint := range arena.checkpoints {
		if checkpoint.Status == CheckpointReached {
			continue
		}
		if checkpoint.Destination.Room != playerRoom {
			continue
		}

		// Discovery currently means entering the destination's room region
		checkpoint.Status = CheckpointDiscovered

		x, y := checkpoint.Destination.WorldPosition()
		distance := math.Hypot(
			arena.Camera.X-x,
			arena.Camera.Y-y,
		)
		if distance <= reachDistance {
			checkpoint.Status = CheckpointReached
		}

		arena.checkpoints[id] = checkpoint
	}
}
