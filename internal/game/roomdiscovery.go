package game

// RoomDiscovery remembers a visited room's identifying features.
// It does not retain the room's terrain.
type RoomDiscovery struct {
	Landmark LandmarkKind
}

// DiscoveredRoom returns remembered information without loading terrain.
// The boolean distinguishes an undiscovered room from an open room.
func (arena *Arena) DiscoveredRoom(id RoomID) (RoomDiscovery, bool) {
	discovery, found := arena.discoveredRooms[id]
	return discovery, found
}

// discoverCurrentRoom records only the region occupied by the player.
func (arena *Arena) discoverCurrentRoom() {
	id := arena.CurrentRoomPosition().Room

	if _, discovered := arena.discoveredRooms[id]; discovered {
		return
	}

	current, loaded := arena.chunks[chunkPosition{X: id.X, Y: id.Y}]
	if !loaded {
		return
	}

	if arena.discoveredRooms == nil {
		arena.discoveredRooms = make(map[RoomID]RoomDiscovery)
	}

	arena.discoveredRooms[id] = RoomDiscovery{
		Landmark: current.Landmark,
	}
}
