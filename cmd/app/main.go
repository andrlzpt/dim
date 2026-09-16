package main

import (
	"runtime"

	"example.com/game/internal/ascii"
	"example.com/game/internal/window"
)

const (
	windowWidth     = 840
	windowHeight    = 624
	blockCellWidth  = 4
	blockCellHeight = 4
	textCellWidth   = 14
	textCellHeight  = 26
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
}

func main() {
	runtime.LockOSThread()
	app := AppState{Screen: ScreenMain}
	blocks := ascii.NewBuffer(windowWidth/blockCellWidth, windowHeight/blockCellHeight)
	text := ascii.NewBuffer(windowWidth/textCellWidth, windowHeight/textCellHeight)

	window.Run(window.Config{
		Width: windowWidth, Height: windowHeight, Title: "New Game",
	}, func(state *window.State) {
		app.Keyboard = state.Keyboard
		app.Controllers = state.Controllers
		handleMenus(state, &app, text)
		drawFrame(blocks, text, &app)
	})
}
