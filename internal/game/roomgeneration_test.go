package game

import (
	"reflect"
	"testing"
)

func TestRoomGenerationRepeatable(t *testing.T) {
	first := NewArena(42)
	second := NewArena(42)

	if !reflect.DeepEqual(first, second) {
		t.Fatal("the same seed must produce the same arena")
	}
	if reflect.DeepEqual(first.Tiles, NewArena(43).Tiles) {
		t.Fatal("these different seeds should produce different layouts")
	}
}

func TestRoomGenerationTraversable(t *testing.T) {
	type point struct {
		X, Y int
	}

	for seed := range int64(100) {
		arena := NewArena(seed)

		if !arena.canStand(arena.Camera.X, arena.Camera.Y, radius) {
			t.Fatalf("seed %d: unsafe spawn", seed)
		}

		start := point{int(arena.Camera.X), int(arena.Camera.Y)}
		queue := []point{start}
		seen := map[point]bool{start: true}
		directions := []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

		for head := 0; head < len(queue); head++ {
			current := queue[head]
			for _, direction := range directions {
				next := point{
					current.X + direction.X,
					current.Y + direction.Y,
				}
				if seen[next] || arena.isWall(next.X, next.Y) {
					continue
				}
				seen[next] = true
				queue = append(queue, next)
			}
		}

		for y, row := range arena.Tiles {
			for x := range row {
				boundary := y == 0 || y == len(arena.Tiles)-1 ||
					x == 0 || x == len(row)-1

				if boundary && !arena.isWall(x, y) {
					t.Fatalf("seed %d: open boundary at %d,%d", seed, x, y)
				}
				if !arena.isWall(x, y) && !seen[point{x, y}] {
					t.Fatalf("seed %d: unreachable floor at %d,%d", seed, x, y)
				}
			}
		}
	}
}
