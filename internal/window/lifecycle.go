package window

import "github.com/go-gl/glfw/v3.4/glfw"

type Loop func(state *State)

type State struct {
	MouseX              float64
	MouseY              float64
	IsMouseInsideWindow bool
	DeltaTime           float64
	ElapsedTime         float64
	FrameCount          uint64
	Keyboard            Controls
	Controllers         []Controller
	MenuUpPressed       bool
	MenuDownPressed     bool
	ConfirmPressed      bool
	MousePressed        bool
	MouseMoved          bool
	PausePressed        bool
	BackPressed         bool
	Focused             bool
	FocusLost           bool
	CaptureMouse        bool
	MouseDeltaX         float64
	MouseDeltaY         float64
	Width               int
	Height              int
	FramebufferWidth    int
	FramebufferHeight   int
}

type Config struct {
	Width      int
	Height     int
	Title      string
	FitDisplay bool
}

func Run(config Config, loop Loop) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	var workX, workY, workWidth, workHeight int

	if config.FitDisplay {
		if monitor := glfw.GetPrimaryMonitor(); monitor != nil {
			workX, workY, workWidth, workHeight = monitor.GetWorkarea()
			if workWidth > 0 && workHeight > 0 {
				config.Width = max(1, workWidth*9/10)
				config.Height = max(1, workHeight*9/10)
			}
		}
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(
		config.Width,
		config.Height,
		config.Title,
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	if config.FitDisplay && workWidth > 0 && workHeight > 0 {
		window.SetPos(
			workX+(workWidth-config.Width)/2,
			workY+(workHeight-config.Height)/2,
		)
	}

	window.MakeContextCurrent()
	state := State{}
	window.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		state.IsMouseInsideWindow = entered
	})

	window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if !focused {
			state.FocusLost = true
		}
	})

	var previousUp, previousDown, previousConfirm, previousMouse, previousPause, mouseCaptured, mouseReady bool
	controllerMenus := controllerMenuTracker{}
	lastTime := glfw.GetTime()

	for !window.ShouldClose() {
		currentTime := glfw.GetTime()
		state.DeltaTime = currentTime - lastTime
		state.ElapsedTime = currentTime
		state.FrameCount++
		lastTime = currentTime
		state.Width, state.Height = window.GetSize()
		state.FramebufferWidth, state.FramebufferHeight = window.GetFramebufferSize()

		state.Focused = window.GetAttrib(glfw.Focused) == glfw.True

		x, y := window.GetCursorPos()
		state.MouseDeltaX = 0
		state.MouseDeltaY = 0
		state.MouseMoved = false

		if mouseReady && state.Focused && !state.FocusLost {
			state.MouseMoved = x != state.MouseX || y != state.MouseY
			if mouseCaptured {
				state.MouseDeltaX = x - state.MouseX
				state.MouseDeltaY = y - state.MouseY
			}
		}

		state.MouseX = x
		state.MouseY = y
		mouseReady = state.Focused && !state.FocusLost

		up := window.GetKey(glfw.KeyUp) == glfw.Press
		down := window.GetKey(glfw.KeyDown) == glfw.Press
		confirm := window.GetKey(glfw.KeyEnter) == glfw.Press
		mouse := window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
		pause := window.GetKey(glfw.KeyEscape) == glfw.Press

		state.MenuUpPressed = up && !previousUp
		state.MenuDownPressed = down && !previousDown
		state.ConfirmPressed = confirm && !previousConfirm
		state.MousePressed = mouse && !previousMouse
		state.PausePressed = pause && !previousPause
		state.BackPressed = false

		previousUp = up
		previousDown = down
		previousConfirm = confirm
		previousMouse = mouse
		previousPause = pause

		state.Keyboard = readKeyboard(window)
		state.Controllers = readControllers()
		controllerMenus.addPresses(&state)

		loop(&state)

		capture := state.CaptureMouse && state.Focused
		if capture != mouseCaptured {
			mode := glfw.CursorNormal
			if capture {
				mode = glfw.CursorDisabled
			}
			window.SetInputMode(glfw.CursorMode, mode)
			mouseCaptured = capture

			// Ignore cursor repositioning caused by changing capture mode.
			mouseReady = false
			state.MouseDeltaX = 0
			state.MouseDeltaY = 0
			state.MouseX, state.MouseY = window.GetCursorPos()
		}

		state.FocusLost = false
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
