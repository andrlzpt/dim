package main

import (
	"github.com/andrlzpt/dim/internal/ascii"
	"github.com/andrlzpt/dim/internal/game"
	"github.com/andrlzpt/dim/internal/window"
)

func currentMenu(app *AppState) (string, []string) {
	if app.Screen == ScreenMain {
		return "New Game", []string{"Start"}
	}
	if app.Paused {
		return "Paused", []string{"Resume", "Main menu"}
	}
	return "", nil
}

func menuWidth(options []string) int {
	width := 0
	for _, label := range options {
		width = max(width, len([]rune(label))+2)
	}
	return width
}

// Shared layout for drawing and mouse hit areas.
func menuPosition(
	text *ascii.Buffer,
	options []string,
	item int,
) (int, int) {
	return (text.Width - menuWidth(options)) / 2,
		text.Height/2 + item*2
}

func drawAppMenu(text *ascii.Buffer, app *AppState) {
	title, options := currentMenu(app)
	if len(options) == 0 {
		return
	}
	drawCenteredText(text, text.Height/2-3, title)
	for i, label := range options {
		prefix := "  "
		if i == app.MenuItem {
			prefix = "> "
		}
		x, y := menuPosition(text, options, i)
		ascii.DrawText(text, x, y, prefix+label)
	}
	drawCenteredText(text, text.Height-2, "arrows/stick: move | enter/A/click: select")
	drawCenteredText(text, text.Height-1, "esc/start: pause | B: back")
}

func drawCenteredText(buffer *ascii.Buffer, row int, text string) {
	ascii.DrawText(buffer, (buffer.Width-len([]rune(text)))/2, row, text)
}

func handleMenus(state *window.State, app *AppState, text *ascii.Buffer) bool {
	if state.PausePressed || state.BackPressed {
		if app.Screen == ScreenArena {
			app.Paused = !app.Paused
			app.MenuItem = 0
		}
		return true
	}
	_, options := currentMenu(app)
	if len(options) == 0 {
		return false
	}

	app.MenuItem = max(0, min(len(options)-1, app.MenuItem))

	if state.MenuUpPressed {
		app.MenuItem = (app.MenuItem + len(options) - 1) % len(options)
	}
	if state.MenuDownPressed {
		app.MenuItem = (app.MenuItem + 1) % len(options)
	}

	hovered := -1
	if state.IsMouseInsideWindow {
		for i := range options {
			x, y := menuPosition(text, options, i)
			width := menuWidth(options)

			left := float64(x * textCellWidth)
			top := float64(y * textCellHeight)
			right := left + float64(width*textCellWidth)
			bottom := top + float64(textCellHeight)

			if state.MouseX >= left && state.MouseX < right &&
				state.MouseY >= top && state.MouseY < bottom {
				hovered = i
				break
			}
		}
	}

	if hovered >= 0 && (state.MouseMoved || state.MousePressed) {
		app.MenuItem = hovered
	}

	if state.ConfirmPressed || (state.MousePressed && hovered >= 0) {
		activateMenuItem(app)
	}

	return true
}

func activateMenuItem(app *AppState) {
	if app.Screen == ScreenMain {
		app.Arena = game.NewArena(42)
		app.Screen = ScreenArena
		app.Paused = false
	} else if app.Paused {
		if app.MenuItem == 1 {
			app.Screen = ScreenMain
		}
		app.Paused = false
	}
	app.MenuItem = 0
}
