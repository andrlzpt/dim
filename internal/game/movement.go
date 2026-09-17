package game

import "math"

// Update advances the player using movement input and elapsed seconds.
func (arena *Arena) Update(moveX, moveY, turn, deltaTime float64) {
	const (
		moveSpeed = 2.5
		turnSpeed = 2.0
		radius    = 0.2
	)

	// Avoid a large jump after a stall or debugger pause.
	dt := max(0.0, min(deltaTime, 0.1))
	if dt == 0 {
		return
	}

	// Diagonal movement should not be faster.
	length := math.Hypot(moveX, moveY)
	if length > 1 {
		moveX /= length
		moveY /= length
	}

	// Small steps prevent crossing walls on slow frames.
	steps := max(1, int(math.Ceil(dt/(1.0/120.0))))
	stepTime := dt / float64(steps)

	camera := &arena.Camera
	for range steps {
		camera.Angle = math.Remainder(
			camera.Angle+turn*turnSpeed*stepTime,
			2*math.Pi,
		)

		forwardX := math.Cos(camera.Angle)
		forwardY := math.Sin(camera.Angle)

		// Negative MoveY means forward; positive MoveX means right.
		dx := (-moveY*forwardX - moveX*forwardY) * moveSpeed * stepTime
		dy := (-moveY*forwardY + moveX*forwardX) * moveSpeed * stepTime

		// Separate axes allow sliding along walls.
		if arena.canStand(camera.X+dx, camera.Y, radius) {
			camera.X += dx
		}
		if arena.canStand(camera.X, camera.Y+dy, radius) {
			camera.Y += dy
		}
	}
}

// Look rotates the camera by an angle in radians.
func (arena *Arena) Look(angle float64) {
	arena.Camera.Angle = math.Remainder(
		arena.Camera.Angle+angle,
		2*math.Pi,
	)
}
