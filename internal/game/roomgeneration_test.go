package game

import (
	"reflect"
	"testing"
)

func TestRoomGenerationRepeatable(t *testing.T) {
	first, firstCamera := generateRooms(42)
	second, secondCamera := generateRooms(42)

	if !reflect.DeepEqual(first, second) || firstCamera != secondCamera {
		t.Fatal("the same seed must produce the same layout and spawn")
	}

	different, _ := generateRooms(43)
	if reflect.DeepEqual(first, different) {
		t.Fatal("these different seeds should produce different layouts")
	}
}

func TestRoomGenerationTraversable(t *testing.T) {
	type point struct {
		X, Y int
	}

	for seed := range int64(100) {
		tiles, camera := generateRooms(seed)

		isWall := func(x, y int) bool {
			return y < 0 || y >= len(tiles) ||
				x < 0 || x >= len(tiles[y]) ||
				tiles[y][x] == '#'
		}

		start := point{int(camera.X), int(camera.Y)}
		if isWall(start.X, start.Y) {
			t.Fatalf("seed %d: spawn is inside a wall", seed)
		}
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
				if seen[next] || isWall(next.X, next.Y) {
					continue
				}
				seen[next] = true
				queue = append(queue, next)
			}
		}

		for y, row := range tiles {
			for x := range row {
				boundary := y == 0 || y == len(tiles)-1 ||
					x == 0 || x == len(row)-1

				if boundary && !isWall(x, y) {
					t.Fatalf("seed %d: open boundary at %d,%d", seed, x, y)
				}
				if !isWall(x, y) && !seen[point{x, y}] {
					t.Fatalf("seed %d: unreachable floor at %d,%d", seed, x, y)
				}
			}
		}
	}
}
