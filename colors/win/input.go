package win

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

// Action is a configurable abstraction of a key press
type Action int

const (
	PLAYER_FORWARD Action = iota
	PLAYER_BACKWARD Action = iota
	PLAYER_LEFT Action = iota
	PLAYER_RIGHT Action = iota
	PROGRAM_QUIT Action = iota
)

type InputManager struct {
	actionToKeyMap map[Action]glfw.Key
	keysPressed [glfw.KeyLast]bool

	firstCursorAction bool
	cursor mgl64.Vec2
	cursorChange mgl64.Vec2
	cursorLast mgl64.Vec2
	bufferedCursorChange mgl64.Vec2
}

func NewInputManager() *InputManager {
	actionToKeyMap := map[Action]glfw.Key{
		PLAYER_FORWARD: glfw.KeyW,
		PLAYER_BACKWARD: glfw.KeyS,
		PLAYER_LEFT: glfw.KeyA,
		PLAYER_RIGHT: glfw.KeyD,
		PROGRAM_QUIT: glfw.KeyEscape,
	}

	return &InputManager{
		actionToKeyMap: actionToKeyMap,
		firstCursorAction: false,
	}
}

// IsActive returns whether the given Action is currently active
func (im *InputManager) IsActive(a Action) bool {
	return im.keysPressed[im.actionToKeyMap[a]]
}

// Cursor returns the value of the cursor at the last time that CheckpointCursorChange() was called.
func (im *InputManager) Cursor() mgl64.Vec2 {
	return im.cursor
}

// CursorChange returns the amount of change in the underlying cursor
// since the last time CheckpointCursorChange was called
func (im *InputManager) CursorChange() mgl64.Vec2 {
	return im.cursorChange
}

// CheckpointCursorChange updates the publicly available Cursor() and CursorChange()
// methods to return the current Cursor and change since last time this method was called.
func (im *InputManager) CheckpointCursorChange() {
	im.cursorChange[0] = im.bufferedCursorChange[0]
	im.cursorChange[1] = im.bufferedCursorChange[1]
	im.cursor[0] = im.cursorLast[0]
	im.cursor[1] = im.cursorLast[1]

	im.bufferedCursorChange[0] = 0
	im.bufferedCursorChange[1] = 0
}

func (im *InputManager) keyCallback(window *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {

	// timing for key events occurs differently from what the program loop requires
	// so just track what key actions occur and then access them in the program loop
	switch action {
	case glfw.Press:
		im.keysPressed[key] = true
	case glfw.Release:
		im.keysPressed[key] = false
	}
}

func (im *InputManager) mouseCallback(window *glfw.Window, xpos, ypos float64) {

	if im.firstCursorAction {
		im.cursorLast[0] = xpos
		im.cursorLast[1] = ypos
		im.firstCursorAction = false
	}

	im.bufferedCursorChange[0] += xpos - im.cursorLast[0]
	im.bufferedCursorChange[1] += ypos - im.cursorLast[1]

	im.cursorLast[0] = xpos
	im.cursorLast[1] = ypos
}
