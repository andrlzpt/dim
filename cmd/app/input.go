package main

import "github.com/andrlzpt/dim/internal/window"

// arenaControls combines keyboard input with the first available gamepad.
func arenaControls(app *AppState) window.Controls {
	controls := app.Keyboard
	if len(app.Controllers) == 0 {
		return controls
	}

	gamepad := app.Controllers[0].Controls

	// Keep movement axes together so keyboard movement takes priority.
	if controls.MoveX == 0 && controls.MoveY == 0 {
		controls.MoveX = gamepad.MoveX
		controls.MoveY = gamepad.MoveY
	}

	if controls.SteerX == 0 {
		controls.SteerX = gamepad.SteerX
	}
	if controls.SteerY == 0 {
		controls.SteerY = gamepad.SteerY
	}

	controls.PrimaryAction = controls.PrimaryAction || gamepad.PrimaryAction
	controls.SecondaryAction = controls.SecondaryAction || gamepad.SecondaryAction

	return controls
}
