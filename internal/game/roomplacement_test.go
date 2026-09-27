package game

import "testing"

func TestGeneratedRoomIdentity(t *testing.T) {
	position := chunkPosition{X: -2, Y: 3}
	current := generateChunk(42, position)

	expected := RoomID{Seed: 42, X: -2, Y: 3}
	if current.RoomID != expected {
		t.Fatalf("expected room %v, got %v", expected, current.RoomID)
	}
}

func TestRoomLandmarkPillars(t *testing.T) {
	seen := make(map[LandmarkKind]bool)

	for seed := range int64(100) {
		current := generateChunk(seed, chunkPosition{0, 0})
		seen[current.Landmark] = true

		walls := 0
		bounds := current.Bounds
		for y := bounds.Y; y < bounds.Y+bounds.Height; y++ {
			for x := bounds.X; x < bounds.X+bounds.Width; x++ {
				if current.Tiles[y][x] == '#' {
					walls++
				}
			}
		}

		expected := 0
		switch current.Landmark {
		case LandmarkNone:
			expected = 0
		case LandmarkPillarPair:
			expected = 2
		case LandmarkFourPillars:
			expected = 4
		default:
			t.Fatalf("unknown landmark: %v", current.Landmark)
		}

		if walls != expected {
			t.Fatalf("seed %d: expected %d pillar cells, got %d", seed, expected, walls)
		}
	}

	if len(seen) != 3 {
		t.Fatal("expected all three room styles across the tested seeds")
	}
}

func TestRoomPlacement(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		arena := NewArena(seed)

		for _, current := range arena.chunks {
			position, found := arena.FindRoomPlacement(current.RoomID, radius)
			if !found {
				t.Fatalf("seed %d: no placement in room %v", seed, current.RoomID)
			}
			if position.Room != current.RoomID {
				t.Fatal("placement belongs to the wrong room")
			}

			x, y := position.WorldPosition()
			if !arena.canStand(x, y, radius) {
				t.Fatal("placement overlaps terrain")
			}

			again, found := arena.FindRoomPlacement(current.RoomID, radius)
			if !found || again != position {
				t.Fatal("unchanged terrain must produce the same placement")
			}
		}
	}
}

func TestRoomPlacementRejectsUnavailableRoom(t *testing.T) {
	arena := NewArena(42)

	for _, id := range []RoomID{
		{Seed: 43, X: 0, Y: 0},
		{Seed: 42, X: 100, Y: 100},
	} {
		if _, found := arena.FindRoomPlacement(id, radius); found {
			t.Fatalf("unexpected placement for room %v", id)
		}
	}
}
