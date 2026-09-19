package game

import (
	"reflect"
	"testing"
)

func TestChunkGenerationRepeatable(t *testing.T) {
	position := chunkPosition{X: -2, Y: 3}
	first := generateChunk(42, position)

	// Generating other chunks must not affect this one.
	generateChunk(42, chunkPosition{X: 5, Y: -1})

	second := generateChunk(42, position)
	if !reflect.DeepEqual(first, second) {
		t.Fatal("the same seed and coordinates must reproduce the chunk")
	}
}

func TestChunkBoundariesMatch(t *testing.T) {
	for seed := int64(0); seed < 5; seed++ {
		for y := -2; y <= 2; y++ {
			for x := -2; x <= 2; x++ {
				current := generateChunk(seed, chunkPosition{x, y})
				east := generateChunk(seed, chunkPosition{x + 1, y})
				south := generateChunk(seed, chunkPosition{x, y + 1})

				var eastOpen, southOpen int
				for i := 0; i < chunkSize; i++ {
					if current.Tiles[i][chunkSize-1] != east.Tiles[i][0] {
						t.Fatalf("seed %d: east mismatch at chunk %d,%d", seed, x, y)
					}
					if current.Tiles[chunkSize-1][i] != south.Tiles[0][i] {
						t.Fatalf("seed %d: south mismatch at chunk %d,%d", seed, x, y)
					}

					if current.Tiles[i][chunkSize-1] == '.' {
						eastOpen++
					}
					if current.Tiles[chunkSize-1][i] == '.' {
						southOpen++
					}
				}

				if eastOpen != 2 || southOpen != 2 {
					t.Fatalf("seed %d: expected two-cell openings at %d,%d", seed, x, y)
				}
			}
		}
	}
}

func TestChunkFloorConnected(t *testing.T) {
	for seed := range int64(20) {
		for _, position := range []chunkPosition{
			{0, 0}, {-1, 0}, {0, -1}, {-2, -3}, {2, 3},
		} {
			current := generateChunk(seed, position)
			start := chunkPosition{chunkSize / 2, chunkSize / 2}

			if current.Tiles[start.Y][start.X] != '.' {
				t.Fatalf("seed %d at %v: center must be open", seed, position)
			}

			queue := []chunkPosition{start}
			seen := map[chunkPosition]bool{start: true}
			directions := []chunkPosition{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

			for head := 0; head < len(queue); head++ {
				point := queue[head]
				for _, direction := range directions {
					next := chunkPosition{
						X: point.X + direction.X,
						Y: point.Y + direction.Y,
					}
					if next.X < 0 || next.X >= chunkSize ||
						next.Y < 0 || next.Y >= chunkSize {
						continue
					}
					if seen[next] || current.Tiles[next.Y][next.X] != '.' {
						continue
					}
					seen[next] = true
					queue = append(queue, next)
				}
			}

			for y, row := range current.Tiles {
				for x := range row {
					if row[x] == '.' && !seen[chunkPosition{x, y}] {
						t.Fatalf(
							"seed %d at %v: disconnected floor at %d,%d",
							seed, position, x, y,
						)
					}
				}
			}
		}
	}
}
