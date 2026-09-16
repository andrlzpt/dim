package main

import (
	"testing"

	"example.com/game/internal/ascii"
	"example.com/game/internal/window"
)

func TestNavigation(t *testing.T) {
	app := AppState{Screen: ScreenMain}
	text := ascii.NewBuffer(windowWidth/textCellWidth, windowHeight/textCellHeight)
	handleMenus(&window.State{ConfirmPressed: true}, &app, text)
	if app.Screen != ScreenArena || app.Paused {
		t.Fatal("Start must open the active arena")
	}
	drawAppMenu(text, &app)
	for _, cell := range text.Cells {
		if cell.Rune != 0 {
			t.Fatal("active arena must have no menu text")
		}
	}
	handleMenus(&window.State{PausePressed: true}, &app, text)
	if !app.Paused {
		t.Fatal("Escape must open pause menu")
	}
	handleMenus(&window.State{BackPressed: true}, &app, text)
	if app.Paused {
		t.Fatal("Back must resume arena")
	}
	handleMenus(&window.State{PausePressed: true}, &app, text)
	handleMenus(&window.State{MenuDownPressed: true}, &app, text)
	handleMenus(&window.State{ConfirmPressed: true}, &app, text)
	if app.Screen != ScreenMain || app.Paused || app.MenuItem != 0 {
		t.Fatal("Main menu must restore title state")
	}
}

func TestMouseStart(t *testing.T) {
	app := AppState{Screen: ScreenMain}
	text := ascii.NewBuffer(windowWidth/textCellWidth, windowHeight/textCellHeight)
	_, options := currentMenu(&app)
	x, y := menuPosition(text, options, 0)
	handleMenus(&window.State{
		IsMouseInsideWindow: true,
		MousePressed:        true,
		MouseX:              float64(x*textCellWidth + 1),
		MouseY:              float64(y*textCellHeight + 1),
	}, &app, text)
	if app.Screen != ScreenArena {
		t.Fatal("clicking Start must open arena")
	}
}
