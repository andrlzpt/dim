package main

import (
	"runtime"

	"github.com/andrlzpt/dim/internal/ascii"
	"github.com/andrlzpt/dim/internal/game"
	"github.com/andrlzpt/dim/internal/window"
)

const (
	windowWidth          = 1260
	windowHeight         = 936
	blockCellWidth       = 4
	blockCellHeight      = 4
	textCellWidth        = 14
	textCellHeight       = 26
	worldCellWidth       = 6
	worldCellHeight      = 11
	worldGlyphWidth      = 7
	worldGlyphHeight     = 13
	mouseSensitivity     = 0.0025
	mouseTiltSensitivity = 0.0005
	controllerTiltSpeed  = 0.25
)

type Screen string

const (
	ScreenMain  Screen = "main"
	ScreenArena Screen = "arena"
)

// AppState holds navigation and current device state independently of gameplay.
type AppState struct {
	Screen      Screen
	MenuItem    int
	Paused      bool
	Keyboard    window.Controls
	Controllers []window.Controller
	Arena       game.Arena
}

func main() {
	runtime.LockOSThread()
	app := AppState{Screen: ScreenMain}
	blocks := ascii.NewBuffer(windowWidth/blockCellWidth, windowHeight/blockCellHeight)
	text := ascii.NewBuffer(windowWidth/textCellWidth, windowHeight/textCellHeight)
	world := ascii.NewBuffer(
		windowWidth/worldCellWidth,
		windowHeight/worldCellHeight,
	)
	window.Run(window.Config{
		Width: windowWidth, Height: windowHeight, Title: "New Game",
	}, func(state *window.State) {
		app.Keyboard = state.Keyboard
		app.Controllers = state.Controllers

		// Losing focus pauses gameplay. Returning requires an explicit resume.
		if (state.FocusLost || !state.Focused) && app.Screen == ScreenArena {
			app.Paused = true
			app.MenuItem = 0
		}

		menuHandled := true
		if state.Focused && !state.FocusLost {
			menuHandled = handleMenus(state, &app, text)
		}

		active := state.Focused &&
			!state.FocusLost &&
			app.Screen == ScreenArena &&
			!app.Paused

		state.CaptureMouse = active

		if active && !menuHandled {
			controls := arenaControls(&app)
			dt := max(0.0, min(state.DeltaTime, 0.1))

			// Mouse and stick Y are negative when pushed upward.
			tilt := -state.MouseDeltaY*mouseTiltSensitivity -
				controls.SteerY*controllerTiltSpeed*dt

			app.Arena.Look(
				state.MouseDeltaX*mouseSensitivity,
				tilt,
			)

			app.Arena.Update(
				controls.MoveX,
				controls.MoveY,
				controls.SteerX,
				state.DeltaTime,
			)
		}

		drawFrame(blocks, world, text, &app)
	})
}
