package game

import "math"

// RoomID identifies a room within a seeded world.
// X and Y are room-grid coordinates, not world-unit positions.
// The current generator assigns one room to each chunk.
type RoomID struct {
	Seed int64
	X    int
	Y    int
}

// Offset returns another room relative to this one.
// Positive X is east; positive Y is south.
// The offset describes grid displacement, not walking distance.
func (id RoomID) Offset(x, y int) RoomID {
	return RoomID{
		Seed: id.Seed,
		X:    id.X + x,
		Y:    id.Y + y,
	}
}

// RoomPosition stores a location relative to a room's region origin.
// X and Y are world units measured from the chunk's top-left corner.
// It contains no references to loaded terrain.
type RoomPosition struct {
	Room RoomID
	X    float64
	Y    float64
}

// WorldPosition converts this reference into world coordinates.
// It does not load terrain or check whether the position is walkable.
func (position RoomPosition) WorldPosition() (float64, float64) {
	return float64(position.Room.X)*float64(chunkSize) + position.X,
		float64(position.Room.Y)*float64(chunkSize) + position.Y
}

// RoomPositionAt identifies a world position without loading its chunk.
func (arena *Arena) RoomPositionAt(x, y float64) RoomPosition {
	roomX, _ := chunkCoordinate(int(math.Floor(x)))
	roomY, _ := chunkCoordinate(int(math.Floor(y)))

	return RoomPosition{
		Room: RoomID{
			Seed: arena.seed,
			X:    roomX,
			Y:    roomY,
		},
		X: x - float64(roomX)*float64(chunkSize),
		Y: y - float64(roomY)*float64(chunkSize),
	}
}

// CurrentRoomPosition returns a snapshot of the player's location.
func (arena *Arena) CurrentRoomPosition() RoomPosition {
	return arena.RoomPositionAt(arena.Camera.X, arena.Camera.Y)
}
