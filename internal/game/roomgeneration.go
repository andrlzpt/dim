package game

import "math/rand"

type room struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r room) center() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

func generateRooms(seed int64) ([]string, Camera) {
	const (
		columns    = 3
		rows       = 2
		slotWidth  = 16
		slotHeight = 14
		width      = columns * slotWidth
		height     = rows * slotHeight
	)

	// Local randomness makes generation independent of other game systems.
	random := rand.New(rand.NewSource(seed))

	grid := make([][]byte, height)
	for y := range grid {
		grid[y] = make([]byte, width)
		for x := range grid[y] {
			grid[y][x] = '#'
		}
	}

	var rooms []room

	// Visit regions in a snake pattern to keep successive rooms nearby.
	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			slotX := column
			if row%2 != 0 {
				slotX = columns - 1 - column
			}

			roomWidth := 6 + random.Intn(5)
			roomHeight := 6 + random.Intn(4)

			r := room{
				X: slotX*slotWidth + 1 +
					random.Intn(slotWidth-roomWidth-1),
				Y: row*slotHeight + 1 +
					random.Intn(slotHeight-roomHeight-1),
				Width:  roomWidth,
				Height: roomHeight,
			}

			carveArea(grid, r.X, r.Y, r.Width, r.Height)

			if len(rooms) > 0 {
				previous := rooms[len(rooms)-1]
				connectRooms(grid, previous, r, random.Intn(2) == 0)
			}

			rooms = append(rooms, r)
		}
	}

	tiles := make([]string, height)
	for y := range grid {
		tiles[y] = string(grid[y])
	}

	spawnX, spawnY := rooms[0].center()
	camera := Camera{
		X: float64(spawnX) + 0.5,
		Y: float64(spawnY) + 0.5,
	}

	return tiles, camera
}

// carveArea opens a rectangle. Callers keep it inside the outer walls.
func carveArea(grid [][]byte, x, y, width, height int) {
	for row := y; row < y+height; row++ {
		for column := x; column < x+width; column++ {
			grid[row][column] = '.'
		}
	}
}

// connectRooms joins room centers with a two-cell-wide L-shaped corridor.
func connectRooms(grid [][]byte, from, to room, horizontalFirst bool) {
	fromX, fromY := from.center()
	toX, toY := to.center()

	horizontalY := fromY
	verticalX := toX
	if !horizontalFirst {
		horizontalY = toY
		verticalX = fromX
	}

	carveArea(
		grid,
		min(fromX, toX),
		horizontalY,
		max(fromX, toX)-min(fromX, toX)+2,
		2,
	)

	carveArea(
		grid,
		verticalX,
		min(fromY, toY),
		2,
		max(fromY, toY)-min(fromY, toY)+2,
	)
}
