package game

// Camera holds a position in world units and a heading in radians.
// X increases east, Y increases south, and angle zero faces east.
type Camera struct {
	X     float64
	Y     float64
	Angle float64
	Tilt  float64 // Fraction of screen height; positive looks up.

}

// Arena holds the room layout and the player's viewpoint.
type Arena struct {
	Tiles  []string
	Camera Camera
}

func NewArena() Arena {
	return Arena{
		Tiles: []string{
			"########################",
			"##############.........#",
			"#........#####.........#",
			"#........#####.........#",
			"#.................##...#",
			"#.................##...#",
			"#........#####.........#",
			"#........#####.........#",
			"##############.........#",
			"##############.........#",
			"########################",
		},
		Camera: Camera{
			X:     2.5,
			Y:     4.0,
			Angle: 0,
		},
	}
}

// isWall treats positions outside the room as solid.
func (arena *Arena) isWall(x, y int) bool {
	if y < 0 || y >= len(arena.Tiles) {
		return true
	}
	if x < 0 || x >= len(arena.Tiles[y]) {
		return true
	}
	return arena.Tiles[y][x] == '#'
}
