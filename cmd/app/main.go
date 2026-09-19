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
	worldCellWidth       = 7
	worldCellHeight      = 13
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
	Screen            Screen
	MenuItem          int
	Paused            bool
	Keyboard          window.Controls
	Controllers       []window.Controller
	Arena             game.Arena
	Width             int
	Height            int
	FramebufferWidth  int
	FramebufferHeight int
}

func main() {
	runtime.LockOSThread()
	app := AppState{Screen: ScreenMain}
	var blocks, world, text *ascii.Buffer

	window.Run(window.Config{
		Width:      windowWidth,
		Height:     windowHeight,
		Title:      "New Game",
		FitDisplay: true,
	}, func(state *window.State) {
		// A minimized window may have no drawable area.
		if state.Width <= 0 || state.Height <= 0 ||
			state.FramebufferWidth <= 0 || state.FramebufferHeight <= 0 {
			state.CaptureMouse = false
			if app.Screen == ScreenArena {
				app.Paused = true
				app.MenuItem = 0
			}
			return
		}

		if blocks == nil || app.Width != state.Width || app.Height != state.Height {
			app.Width = state.Width
			app.Height = state.Height

			blocks = ascii.NewBuffer(
				max(1, app.Width/blockCellWidth),
				max(1, app.Height/blockCellHeight),
			)
			world = ascii.NewBuffer(
				max(1, app.Width/worldCellWidth),
				max(1, app.Height/worldCellHeight),
			)
			text = ascii.NewBuffer(
				max(1, app.Width/textCellWidth),
				max(1, app.Height/textCellHeight),
			)
		}

		app.FramebufferWidth = state.FramebufferWidth
		app.FramebufferHeight = state.FramebufferHeight
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
